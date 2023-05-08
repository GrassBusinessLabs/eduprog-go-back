package eduprog

import (
	"errors"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/app"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/requests"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

const (
	MANDATORY = "MANDATORY"
	BLOC      = "BLOC" // ВБ
	//	LIST      = "LIST" // ВП
)

type EduprogcompController struct {
	eduprogcompService app.EduprogcompService
	eduprogService     app.EduprogService
	eduprogController  EduprogController
}

func NewEduprogcompController(es app.EduprogcompService, eps app.EduprogService, edc EduprogController) EduprogcompController {
	return EduprogcompController{
		eduprogcompService: es,
		eduprogService:     eps,
		eduprogController:  edc,
	}
}

func (c EduprogcompController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		eduprogcomp, err := requests.Bind(r, requests.CreateEduprogcompRequest{}, domain.Eduprogcomp{})
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		comps, _ := c.eduprogcompService.SortComponentsByMnS(eduprogcomp.EduprogId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		//Code generation logic
		if eduprogcomp.Type == "ОК" { //if educomp type is "OK"
			var maxCode uint64 = 0
			eduprogcomp.Category = MANDATORY
			for i := range comps.Mandatory {
				if comps.Mandatory[i].Name == eduprogcomp.Name {
					log.Printf("EduprogcompController: %s", err)
					controllers.BadRequest(w, errors.New("eduprog component with this name already exists"))
					return
				}
				temp, _ := strconv.ParseUint(comps.Mandatory[i].Code, 10, 64)
				if i == 0 || temp > maxCode {
					maxCode = temp
				}
			}
			eduprogcomp.Code = strconv.FormatUint(maxCode+1, 10)
			eduprogcomp.BlockName = ""
			eduprogcomp.BlockNum = ""
		} else if eduprogcomp.Type == "ВБ" { //if educomp type is "VB"
			var maxCode uint64 = 0
			eduprogcomp.Category = BLOC
			for i := range comps.Selective {
				for _, comp := range comps.Selective[i].CompsInBlock {
					if comp.Name == eduprogcomp.Name {
						log.Printf("EduprogcompController: %s", err)
						controllers.BadRequest(w, errors.New("eduprog component with this name already exists"))
						return
					}
					if comp.BlockNum == eduprogcomp.BlockNum {
						temp, _ := strconv.ParseUint(comp.Code, 10, 64)
						if i == 0 || temp > maxCode {
							maxCode = temp
						}
					}
				}

			}
			eduprogcomp.Code = strconv.FormatUint(maxCode+1, 10)
		}

		//Free credits check

		eduprog, err := c.eduprogService.FindById(eduprogcomp.EduprogId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		creditsDto, err := c.eduprogController.GetCreditsInfo(comps, eduprog.EducationLevel)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		if eduprogcomp.Type == "ОК" {
			if eduprogcomp.Credits > creditsDto.MandatoryFreeCredits {
				log.Printf("EduprogcompController: %s", err)
				controllers.BadRequest(w, errors.New("too much credits"))
				return
			}
		} else if eduprogcomp.Type == "ВБ" {
			if eduprogcomp.Credits > creditsDto.SelectiveFreeCredits {
				log.Printf("EduprogcompController: %s", err)
				controllers.BadRequest(w, errors.New("too much credits or wrong number (must be > 0)"))
				return
			}
		} else {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, errors.New(`wrong type, it can only be "ОК" or "ВБ"`))
			return
		}

		eduprogcomp, err = c.eduprogcompService.Save(eduprogcomp)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var eduprogcompDto resources.EduprogcompDto
		controllers.Created(w, eduprogcompDto.DomainToDto(eduprogcomp))
	}
}

func (c EduprogcompController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epcId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomp, err := requests.Bind(r, requests.UpdateEduprogcompRequest{}, domain.Eduprogcomp{})
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcompById, err := c.eduprogcompService.FindById(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}
		eduprogcomp.Id = eduprogcompById.Id

		comps, _ := c.eduprogcompService.SortComponentsByMnS(eduprogcomp.EduprogId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		maxCode := 1

		if eduprogcomp.Type == "ОК" { //if educomp type is "OK"
			eduprogcomp.Category = MANDATORY
			for i := range comps.Mandatory {
				if comps.Mandatory[i].Name == eduprogcomp.Name && comps.Mandatory[i].Id == eduprogcomp.Id {
					log.Printf("EduprogcompController: %s", err)
					controllers.BadRequest(w, errors.New("eduprog component with this name already exists"))
					return
				}
			}
			eduprogcomp.BlockName = ""
			eduprogcomp.BlockNum = ""
		} else if eduprogcomp.Type == "ВБ" { //if educomp type is "VB"
			eduprogcomp.Category = BLOC
			for i := range comps.Selective {
				for _, comp := range comps.Selective[i].CompsInBlock {
					if comp.Name == eduprogcomp.Name && (comp.Id > eduprogcomp.Id || comp.Id < eduprogcomp.Id) {
						log.Printf("EduprogcompController: %s", err)
						controllers.BadRequest(w, errors.New("eduprog component with this name already exists"))
						return
					}
					if comp.BlockNum == eduprogcomp.BlockNum {
						maxCode++
					}
				}
			}
		}

		//Free credits check

		eduprog, err := c.eduprogService.FindById(eduprogcomp.EduprogId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		creditsDto, err := c.eduprogController.GetCreditsInfo(comps, eduprog.EducationLevel)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		if eduprogcomp.Type == "ОК" {
			if eduprogcomp.Credits+(creditsDto.MandatoryCredits-eduprogcompById.Credits) > creditsDto.MandatoryCreditsForLevel {
				log.Printf("EduprogcompController: %s", err)
				controllers.BadRequest(w, errors.New("too much credits"))
				return
			}
		} else if eduprogcomp.Type == "ВБ" {
			if eduprogcomp.Credits+(creditsDto.SelectiveCredits-eduprogcompById.Credits) > creditsDto.SelectiveCreditsForLevel {
				log.Printf("EduprogcompController: %s", err)
				controllers.BadRequest(w, errors.New("too much credits"))
				return
			}
		} else {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, errors.New(`wrong type, it can only be "ОК" or "ВБ"`))
			return
		}

		eduprogcomp.Id = id
		eduprogcomp.Code = eduprogcompById.Code
		eduprogcomp, err = c.eduprogcompService.Update(eduprogcomp, id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		if eduprogcomp.Type == "ВБ" {
			for i := range comps.Selective {
				for i2, comp := range comps.Selective[i].CompsInBlock {
					if comp.Id == eduprogcomp.Id {
						comps.Selective[i].CompsInBlock[i2].BlockNum = eduprogcomp.BlockNum
						comps.Selective[i].CompsInBlock[i2].BlockName = eduprogcomp.BlockName
						comps.Selective[i].CompsInBlock[i2].Name = eduprogcomp.Name
						comps.Selective[i].CompsInBlock[i2].ControlType = eduprogcomp.ControlType
						comps.Selective[i].CompsInBlock[i2].Credits = eduprogcomp.Credits
						comps.Selective[i].CompsInBlock[i2].Type = eduprogcomp.Type
						comps.Selective[i].CompsInBlock[i2].Code = strconv.Itoa(maxCode)
					} else {
						comps.Selective[i].CompsInBlock[i2].BlockNum = comps.Selective[i].BlockNum
						comps.Selective[i].CompsInBlock[i2].BlockName = comps.Selective[i].BlockName
					}
					_, _ = c.eduprogcompService.Update(comps.Selective[i].CompsInBlock[i2], comps.Selective[i].CompsInBlock[i2].Id)
				}
			}
			eduprogcomps, err := c.eduprogcompService.SortComponentsByMnS(eduprogcomp.EduprogId)
			if err != nil {
				log.Printf("EduprogcompController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
			for i, elem := range eduprogcomps.Selective {
				elem.BlockNum = strconv.Itoa(i + 1)
				eduprogcomps.Selective[i] = elem
			}
			for i := range eduprogcomps.Selective {
				for _, comp := range eduprogcomps.Selective[i].CompsInBlock {
					comp.BlockNum = eduprogcomps.Selective[i].BlockNum
					_, _ = c.eduprogcompService.Update(comp, comp.Id)
				}
			}
		}

		var eduprogcompDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompDto.DomainToDto(eduprogcomp))
	}
}

func (c EduprogcompController) ReplaceComp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		educompId, err := strconv.ParseUint(r.URL.Query().Get("edcompId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		putAfterId, err := strconv.ParseInt(r.URL.Query().Get("putAfter"), 10, 64) // Now its just code (OK)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		educompById, err := c.eduprogcompService.FindById(educompId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		var targetEdcompCode string
		if putAfterId == 0 {
			targetEdcompCode = strconv.FormatInt(putAfterId, 10)
		} else {
			targetEdcompCode = strconv.FormatInt(putAfterId+1, 10)
		}

		//var targetEducompById domain.Eduprogcomp
		//
		//if putAfterId != 0 {
		//	targetEducompById, err = c.eduprogcompService.FindById(putAfterId)
		//	if err != nil {
		//		log.Printf("EduprogcompController: %s", err)
		//		controllers.InternalServerError(w, err)
		//		return
		//	}
		//} else {
		//	targetEducompById.Code = "0"
		//	targetEducompById.Type = educompById.Type
		//}

		eduprogcomps, err := c.eduprogcompService.SortComponentsByMnS(educompById.EduprogId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		if educompById.Type == "ОК" && targetEdcompCode != educompById.Code {
			eduprogcomps.Mandatory = moveElement(eduprogcomps.Mandatory, educompById.Code, targetEdcompCode)
		}
		//else if educompById.Type == "ВБ" && targetEducompById.Type == "ВБ" {
		//	if educompById.BlockNum == targetEducompById.BlockNum {
		//		eduprogcomps.Selective = moveElement(eduprogcomps.Selective, educompById.Code, targetEducompById.Code)
		//	}
		//}

		for i := range eduprogcomps.Mandatory {
			_, _ = c.eduprogcompService.Update(eduprogcomps.Mandatory[i], eduprogcomps.Mandatory[i].Id)
			if err != nil {
				log.Printf("EduprogcompController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
		}
		//for i := range eduprogcomps.Selective {
		//	_, _ = c.eduprogcompService.Update(eduprogcomps.Selective[i], eduprogcomps.Selective[i].Id)
		//	if err != nil {
		//		log.Printf("EduprogcompController: %s", err)
		//		controllers.InternalServerError(w, err)
		//		return
		//	}
		//}

		var eduprogcompDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompDto.DomainToDtoWCompCollection(eduprogcomps, eduprogcomps.Selective))
	}
}

func moveElement(slice []domain.Eduprogcomp, code string, afterCode string) []domain.Eduprogcomp {
	var index int = -1
	for i, elem := range slice {
		if elem.Code == code {
			index = i
			break
		}
	}
	if index == -1 {
		// Element with given code not found, return the original slice.
		return slice
	}

	elem := slice[index]
	slice = append(slice[:index], slice[index+1:]...)

	var afterIndex int = -1
	for i, elem := range slice {
		if elem.Code == afterCode {
			afterIndex = i
			break
		}
	}

	var insertIndex int
	if afterIndex == -1 {
		insertIndex = 0
	} else {
		insertIndex = afterIndex + 1
	}

	slice = append(slice[:insertIndex], append([]domain.Eduprogcomp{elem}, slice[insertIndex:]...)...)

	for i, elem := range slice {
		elem.Code = strconv.Itoa(i + 1)
		slice[i] = elem
	}

	return slice
}

func (c EduprogcompController) UpdateMandatoryComps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogcomps, err := requests.Bind(r, requests.SendEduprogcompSliceRequest{}, []domain.Eduprogcomp{})
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		for i, elem := range eduprogcomps {
			elem.Code = strconv.Itoa(i + 1)
			eduprogcomps[i] = elem
		}

		for i := range eduprogcomps {
			_, _ = c.eduprogcompService.Update(eduprogcomps[i], eduprogcomps[i].Id)
			if err != nil {
				log.Printf("EduprogcompController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
		}

		var eduprogcompDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompDto.DomainToDtoCollection(eduprogcomps))
	}
}

//func moveAfterCode(s []domain.Eduprogcomp, putAfterCode int, chosenEduprog domain.Eduprogcomp) []domain.Eduprogcomp {
//	var insertIndex int = -1
//	for i, eduprog := range s {
//		if eduprog.Code == strconv.Itoa(putAfterCode) {
//			// found the position to move after
//			insertIndex = i + 1
//			continue
//		}
//	}
//	for i, eduprog := range s {
//		if eduprog.Id == chosenEduprog.Id {
//			// found the chosen eduprog to move
//			if insertIndex != -1 {
//				// remove chosen eduprog from the slice
//				s = append(s[:i], s[i+1:]...)
//				// insert chosen eduprog after the position
//				if i < insertIndex {
//					// shift insertIndex since we removed an element before it
//					insertIndex--
//				}
//				s = append(s[:insertIndex], append([]domain.Eduprogcomp{chosenEduprog}, s[insertIndex:]...)...)
//				return s
//			} else {
//				// the putAfterCode was not found
//				return s
//			}
//		}
//	}
//	// the chosen eduprog was not found
//	return s
//}

func (c EduprogcompController) UpdateVBName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epcId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomp, err := requests.Bind(r, requests.UpdateBlockName{}, domain.Eduprogcomp{})
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomps, err := c.eduprogcompService.SortComponentsByMnS(id)

		for i := range eduprogcomps.Selective {
			if eduprogcomp.BlockName == eduprogcomps.Selective[i].BlockName && eduprogcomp.BlockNum != eduprogcomps.Selective[i].BlockNum {
				log.Printf("EduprogcompController: %s", err)
				controllers.BadRequest(w, errors.New(`block with this name already exists`))
				return
			}
		}

		vbBlock, err := c.eduprogcompService.FindByBlockNum(id, eduprogcomp.BlockNum)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var result []domain.Eduprogcomp
		for i := range vbBlock {
			edcompById, err := c.eduprogcompService.FindById(vbBlock[i].Id)
			if err != nil {
				log.Printf("EduprogcompController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
			edcompById.BlockName = eduprogcomp.BlockName
			updEduprogcomp, err := c.eduprogcompService.Update(edcompById, edcompById.Id)
			if err != nil {
				log.Printf("EduprogcompController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}

			result = append(result, updEduprogcomp)
		}

		var eduprogcompsDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompsDto.DomainToDtoCollection(result))
	}
}

func (c EduprogcompController) GetVBBlocksInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epcId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomps, err := c.eduprogcompService.SortComponentsByMnS(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompsDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompsDto.BlockInfoToDtoCollection(eduprogcomps.Selective))
	}
}

func (c EduprogcompController) ShowList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogcomps, err := c.eduprogcompService.ShowList()
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompsDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompsDto.DomainToDtoCollection(eduprogcomps))
	}
}

func (c EduprogcompController) ShowListByEduprogId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epcId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomps, err := c.eduprogcompService.ShowListByEduprogId(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompsDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompsDto.DomainToDtoCollection(eduprogcomps))
	}
}

func (c EduprogcompController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epcId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var eduprogcompDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompDto.DomainToDto(eduprogcomp))
	}
}

func (c EduprogcompController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epcId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		err = c.eduprogcompService.Delete(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcomps, err := c.eduprogcompService.SortComponentsByMnS(eduprogcomp.EduprogId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for i, elem := range eduprogcomps.Mandatory {
			elem.Code = strconv.Itoa(i + 1)
			eduprogcomps.Mandatory[i] = elem
		}

		for i, elem := range eduprogcomps.Selective {
			elem.BlockNum = strconv.Itoa(i + 1)
			eduprogcomps.Selective[i] = elem
		}

		for i := range eduprogcomps.Mandatory {
			_, _ = c.eduprogcompService.Update(eduprogcomps.Mandatory[i], eduprogcomps.Mandatory[i].Id)
			if err != nil {
				log.Printf("EduprogcompController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
		}

		for i := range eduprogcomps.Selective {
			for _, comp := range eduprogcomps.Selective[i].CompsInBlock {
				comp.BlockNum = eduprogcomps.Selective[i].BlockNum
				_, _ = c.eduprogcompService.Update(comp, comp.Id)
			}
		}

		controllers.Ok(w)
	}
}
