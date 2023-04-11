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
				if comps.Selective[i].Name == eduprogcomp.Name {
					log.Printf("EduprogcompController: %s", err)
					controllers.BadRequest(w, errors.New("eduprog component with this name already exists"))
					return
				}
				if comps.Selective[i].BlockNum == eduprogcomp.BlockNum {
					temp, _ := strconv.ParseUint(comps.Selective[i].Code, 10, 64)
					if i == 0 || temp > maxCode {
						maxCode = temp
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

		comps, _ := c.eduprogcompService.SortComponentsByMnS(eduprogcomp.EduprogId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

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
				if comps.Selective[i].Name == eduprogcomp.Name && comps.Selective[i].Id == eduprogcomp.Id {
					log.Printf("EduprogcompController: %s", err)
					controllers.BadRequest(w, errors.New("eduprog component with this name already exists"))
					return
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
		putAfterId, err := strconv.ParseUint(r.URL.Query().Get("putAfterId"), 10, 64)
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

		var targetEducompById domain.Eduprogcomp

		if putAfterId != 0 {
			targetEducompById, err = c.eduprogcompService.FindById(putAfterId)
			if err != nil {
				log.Printf("EduprogcompController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
		} else {
			targetEducompById.Code = "0"
			targetEducompById.Type = educompById.Type
		}

		eduprogcomps, err := c.eduprogcompService.SortComponentsByMnS(educompById.EduprogId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		if educompById.Type == "ОК" && targetEducompById.Type == "ОК" {
			eduprogcomps.Mandatory = moveElement(eduprogcomps.Mandatory, educompById.Code, targetEducompById.Code)
		} else if educompById.Type == "ВБ" && targetEducompById.Type == "ВБ" {
			if educompById.BlockNum == targetEducompById.BlockNum {
				eduprogcomps.Selective = moveElement(eduprogcomps.Selective, educompById.Code, targetEducompById.Code)
			}
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
			_, _ = c.eduprogcompService.Update(eduprogcomps.Selective[i], eduprogcomps.Selective[i].Id)
			if err != nil {
				log.Printf("EduprogcompController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
		}

		var eduprogcompDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompDto.DomainToDtoWCompCollection(eduprogcomps))
	}
}

func moveElement(slice []domain.Eduprogcomp, code string, afterCode string) []domain.Eduprogcomp {
	// Find the index of the element with the given code.
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

	// Remove the element with the given code from the slice.
	elem := slice[index]
	slice = append(slice[:index], slice[index+1:]...)

	// Find the index of the element with the given afterCode.
	var afterIndex int = -1
	for i, elem := range slice {
		if elem.Code == afterCode {
			afterIndex = i
			break
		}
	}

	// Determine the index to insert the element.
	var insertIndex int
	if afterIndex == -1 {
		// Element with given afterCode not found, insert at the beginning.
		insertIndex = 0
	} else {
		// Insert the element after the element with the given afterCode.
		insertIndex = afterIndex + 1
	}

	// Insert the element at the determined index.
	slice = append(slice[:insertIndex], append([]domain.Eduprogcomp{elem}, slice[insertIndex:]...)...)

	// Reassign codes to ensure they are sequential starting from 1.
	for i, elem := range slice {
		elem.Code = strconv.Itoa(i + 1)
		slice[i] = elem
	}

	// Return the updated slice.
	return slice
}

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

		blockInfo := c.eduprogcompService.GetVBBlocksDomain(eduprogcomps)

		var eduprogcompsDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompsDto.BlockInfoToDtoCollection(blockInfo))
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

		eduprogcomps, err := c.eduprogcompService.ShowListByEduprogId(eduprogcomp.EduprogId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		if eduprogcomp.Type == "ОК" {
			for i := range eduprogcomps {
				if eduprogcomps[i].Type == eduprogcomp.Type {
					educompsCode, err := strconv.ParseUint(eduprogcomps[i].Code, 10, 64)
					if err != nil {
						log.Printf("EduprogcompController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
					educompCode, _ := strconv.ParseUint(eduprogcomp.Code, 10, 64)
					if err != nil {
						log.Printf("EduprogcompController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
					if educompsCode > educompCode {
						eduprogcomps[i].Code = strconv.FormatUint(educompsCode-1, 10)
						_, _ = c.eduprogcompService.Update(eduprogcomps[i], eduprogcomps[i].Id)
						if err != nil {
							log.Printf("EduprogcompController: %s", err)
							controllers.InternalServerError(w, err)
							return
						}
					}
				}
			}
		} else if eduprogcomp.Type == "ВБ" {
			for i := range eduprogcomps {
				educompsCode, err := strconv.ParseUint(eduprogcomps[i].Code, 10, 64)
				if err != nil {
					log.Printf("EduprogcompController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
				educompCode, _ := strconv.ParseUint(eduprogcomp.Code, 10, 64)
				if err != nil {
					log.Printf("EduprogcompController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}

				if educompsCode > educompCode && eduprogcomp.BlockNum == eduprogcomps[i].BlockNum {
					eduprogcomps[i].Code = strconv.FormatUint(educompsCode-1, 10)
					_, _ = c.eduprogcompService.Update(eduprogcomps[i], eduprogcomps[i].Id)
					if err != nil {
						log.Printf("EduprogcompController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
				}
			}
		}

		controllers.Ok(w)
	}
}
