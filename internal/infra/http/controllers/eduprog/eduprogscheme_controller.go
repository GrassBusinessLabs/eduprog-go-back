package eduprog

import (
	"errors"
	"fmt"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/app"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/requests"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
	"sort"

	"log"
	"net/http"
	"strconv"
)

type EduprogschemeController struct {
	eduprogschemeService app.EduprogschemeService
	eduprogcompService   app.EduprogcompService
	disciplineService    app.DisciplineService
}

func NewEduprogschemeController(ess app.EduprogschemeService, ecs app.EduprogcompService, ds app.DisciplineService) EduprogschemeController {
	return EduprogschemeController{
		eduprogschemeService: ess,
		eduprogcompService:   ecs,
		disciplineService:    ds,
	}
}

func (c EduprogschemeController) SetComponentToEdprogscheme() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		eduprogscheme, err := requests.Bind(r, requests.SetComponentToEdprogschemeRequest{}, domain.Eduprogscheme{})
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogschemes, err := c.eduprogschemeService.ShowSchemeByEduprogId(eduprogscheme.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogscheme.EduprogcompId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		discipline, err := c.disciplineService.FindById(eduprogscheme.DisciplineId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		totalCompCredits := eduprogscheme.CreditsPerSemester

		for i := range eduprogschemes {
			if eduprogschemes[i].EduprogcompId == eduprogscheme.EduprogcompId {
				totalCompCredits = totalCompCredits + eduprogschemes[i].CreditsPerSemester
				if eduprogschemes[i].SemesterNum == eduprogscheme.SemesterNum {
					log.Printf("EduprogschemeController: %s", err)
					controllers.BadRequest(w, errors.New("this component already exists in this semester"))
					return
				}
			}
		}

		if discipline.Rows < eduprogscheme.Row || eduprogscheme.Row == 0 {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, fmt.Errorf("max row num for this discipline is %d, and it cant be zero", discipline.Rows))
			return
		}

		if totalCompCredits > eduprogcomp.Credits {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, errors.New("too much credits per semester"))
			return
		}

		eduprogscheme, err = c.eduprogschemeService.SetComponentToEdprogscheme(eduprogscheme)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Created(w, eduprogschemeDto.DomainToDto(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) UpdateComponentInEduprogscheme() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogscheme, err := requests.Bind(r, requests.UpdateComponentInEduprogschemeRequest{}, domain.Eduprogscheme{})
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		//eduprogschemes, err := c.eduprogschemeService.ShowSchemeByEduprogId(eduprogscheme.EduprogId)
		//if err != nil {
		//	log.Printf("EduprogschemeController: %s", err)
		//	controllers.BadRequest(w, err)
		//	return
		//}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogscheme.EduprogcompId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		discipline, err := c.disciplineService.FindById(eduprogscheme.DisciplineId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		totalCompCredits := eduprogscheme.CreditsPerSemester

		//for i := range eduprogschemes {
		//	if eduprogschemes[i].EduprogcompId == eduprogscheme.EduprogcompId {
		//		//totalCompCredits = totalCompCredits + eduprogschemes[i].CreditsPerSemester
		//		if eduprogschemes[i].SemesterNum == eduprogscheme.SemesterNum {
		//			log.Printf("EduprogschemeController: %s", err)
		//			controllers.BadRequest(w, errors.New("this component already exists in this semester"))
		//			return
		//		}
		//	}
		//}

		if discipline.Rows < eduprogscheme.Row || eduprogscheme.Row == 0 {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, fmt.Errorf("max row num for this discipline is %d, and it cant be zero", discipline.Rows))
			return
		}

		if totalCompCredits > eduprogcomp.Credits {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, errors.New("too much credits per semester"))
			return
		}

		eduprogscheme.Id = id
		eduprogscheme, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(eduprogscheme, id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDto(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) ExpandComponentInEduprogscheme() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		expandRequest, err := requests.Bind(r, requests.ExpandComponentInEduprogschemeRequest{}, domain.ExpandEduprogScheme{})
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogschemeComponent, err := c.eduprogschemeService.FindById(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.ShowSchemeByEduprogId(eduprogschemeComponent.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogschemeComponent.EduprogcompId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		eduprogcompMaxCredits := eduprogcomp.Credits
		creditsInScheme := 0.0

		var schemeComponentsList []domain.Eduprogscheme

		for i := range eduprogscheme {
			if eduprogscheme[i].EduprogcompId == eduprogschemeComponent.EduprogcompId {
				schemeComponentsList = append(schemeComponentsList, eduprogscheme[i])
				//creditsInScheme = eduprogscheme[i].CreditsPerSemester - expandRequest.CreditsPerSemester
			}
		}

		if creditsInScheme+expandRequest.CreditsPerSemester > eduprogcompMaxCredits {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, errors.New("too much credits"))
			return
		}

		var createdSchemeComponent domain.Eduprogscheme
		var componentWithMaxSemester domain.Eduprogscheme
		var componentWithMinSemester domain.Eduprogscheme

		if expandRequest.ExpandTo == "LEFT" { // EXPANDING SCHEME COMP TO LEFT (+semester)
			if len(schemeComponentsList) < 2 {
				if schemeComponentsList[0].SemesterNum > 1 {
					createdSchemeComponent = schemeComponentsList[0]
					createdSchemeComponent.CreditsPerSemester = expandRequest.CreditsPerSemester
					createdSchemeComponent.SemesterNum = createdSchemeComponent.SemesterNum - 1
					schemeComponentsList[0].CreditsPerSemester = schemeComponentsList[0].CreditsPerSemester - expandRequest.CreditsPerSemester
					if schemeComponentsList[0].CreditsPerSemester <= 0 {
						log.Printf("EduprogschemeController: %s", err)
						controllers.BadRequest(w, errors.New("too much credits"))
						return
					}

					_, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(schemeComponentsList[0], schemeComponentsList[0].Id)
					if err != nil {
						log.Printf("EduprogschemeController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
					createdSchemeComponent, err = c.eduprogschemeService.SetComponentToEdprogscheme(createdSchemeComponent)
					if err != nil {
						log.Printf("EduprogschemeController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
				}
			} else if len(schemeComponentsList) > 1 {
				var minSemester uint64 = 8
				for i := range schemeComponentsList {
					if minSemester > schemeComponentsList[i].SemesterNum {
						minSemester = schemeComponentsList[i].SemesterNum
						componentWithMinSemester = schemeComponentsList[i]
					}
				}
				if componentWithMinSemester.SemesterNum > 1 {
					createdSchemeComponent = componentWithMinSemester
					createdSchemeComponent.CreditsPerSemester = expandRequest.CreditsPerSemester
					createdSchemeComponent.SemesterNum = createdSchemeComponent.SemesterNum - 1
					componentWithMinSemester.CreditsPerSemester = componentWithMinSemester.CreditsPerSemester - expandRequest.CreditsPerSemester
					if componentWithMinSemester.CreditsPerSemester <= 0 {
						log.Printf("EduprogschemeController: %s", err)
						controllers.BadRequest(w, errors.New("too much credits"))
						return
					}
					_, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(componentWithMinSemester, componentWithMinSemester.Id)
					if err != nil {
						log.Printf("EduprogschemeController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
					createdSchemeComponent, err = c.eduprogschemeService.SetComponentToEdprogscheme(createdSchemeComponent)
					if err != nil {
						log.Printf("EduprogschemeController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
				}
			}
		} else if expandRequest.ExpandTo == "RIGHT" { // EXPANDING SCHEME COMP TO LEFT (-semester)
			if len(schemeComponentsList) < 2 {
				if schemeComponentsList[0].SemesterNum < 8 {
					createdSchemeComponent = schemeComponentsList[0]
					createdSchemeComponent.CreditsPerSemester = expandRequest.CreditsPerSemester
					createdSchemeComponent.SemesterNum = createdSchemeComponent.SemesterNum + 1
					schemeComponentsList[0].CreditsPerSemester = schemeComponentsList[0].CreditsPerSemester - expandRequest.CreditsPerSemester
					if schemeComponentsList[0].CreditsPerSemester <= 0 {
						log.Printf("EduprogschemeController: %s", err)
						controllers.BadRequest(w, errors.New("too much credits"))
						return
					}

					_, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(schemeComponentsList[0], schemeComponentsList[0].Id)
					if err != nil {
						log.Printf("EduprogschemeController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
					createdSchemeComponent, err = c.eduprogschemeService.SetComponentToEdprogscheme(createdSchemeComponent)
					if err != nil {
						log.Printf("EduprogschemeController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
				}
			} else if len(schemeComponentsList) > 1 {
				var maxSemester uint64 = 0
				for i := range schemeComponentsList {
					if maxSemester < schemeComponentsList[i].SemesterNum {
						maxSemester = schemeComponentsList[i].SemesterNum
						componentWithMaxSemester = schemeComponentsList[i]
					}
				}
				if componentWithMaxSemester.SemesterNum < 8 {
					createdSchemeComponent = componentWithMaxSemester
					createdSchemeComponent.CreditsPerSemester = expandRequest.CreditsPerSemester
					createdSchemeComponent.SemesterNum = createdSchemeComponent.SemesterNum + 1
					componentWithMaxSemester.CreditsPerSemester = componentWithMaxSemester.CreditsPerSemester - expandRequest.CreditsPerSemester
					if componentWithMaxSemester.CreditsPerSemester <= 0 {
						log.Printf("EduprogschemeController: %s", err)
						controllers.BadRequest(w, errors.New("too much credits"))
						return
					}
					_, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(componentWithMaxSemester, componentWithMaxSemester.Id)
					if err != nil {
						log.Printf("EduprogschemeController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
					createdSchemeComponent, err = c.eduprogschemeService.SetComponentToEdprogscheme(createdSchemeComponent)
					if err != nil {
						log.Printf("EduprogschemeController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
				}
			}
		}

		eduprogschemeToShow, err := c.eduprogschemeService.ShowSchemeByEduprogId(eduprogcomp.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcompsToShow, err := c.eduprogcompService.ShowListByEduprogId(eduprogcomp.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogschemeToShow, eduprogcompsToShow))
	}
}

func (c EduprogschemeController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.FindById(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogscheme.EduprogcompId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDto(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) FindBySemesterNum() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sNum, err := strconv.ParseUint(chi.URLParam(r, "sNum"), 10, 16)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.FindBySemesterNum(uint16(sNum), id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.ShowListByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) ShowSchemeByEduprogId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "sNum"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.ShowSchemeByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.ShowListByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) ShowFreeComponents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "sNum"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		sortOrder := r.URL.Query().Get("order")
		if sortOrder != "Az" && sortOrder != "Za" {
			controllers.BadRequest(w, errors.New("only Az (alphabetic) or Za (naoborot)"))
			return
		}

		eduprogcomps, err := c.eduprogcompService.ShowListByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.ShowSchemeByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for i := range eduprogcomps {
			eduprogcomps[i].FreeCredits = eduprogcomps[i].Credits
			for i2 := range eduprogscheme {
				if eduprogcomps[i].Id == eduprogscheme[i2].EduprogcompId {
					eduprogcomps[i].FreeCredits = eduprogcomps[i].FreeCredits - eduprogscheme[i2].CreditsPerSemester
				}
			}
		}

		var result []domain.Eduprogcomp

		for i := range eduprogcomps {
			if eduprogcomps[i].FreeCredits > 0 {
				result = append(result, eduprogcomps[i])
			}
		}

		//result = sortByCode(result)
		if sortOrder == "Az" {
			sort.Slice(result, func(i, j int) bool {
				return result[i].Name < result[j].Name
			})
		} else if sortOrder == "Za" {
			sort.Slice(result, func(i, j int) bool {
				return result[i].Name > result[j].Name
			})
		}

		var eduprogcompDto resources.EduprogcompDtoWithFreeCredits
		controllers.Success(w, eduprogcompDto.DomainToDtoCollectionWithFreeCredits(result))
	}
}

func (c EduprogschemeController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		err = c.eduprogschemeService.Delete(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		controllers.Ok(w)
	}
}

func unique(compSlice []domain.Eduprogcomp) []domain.Eduprogcomp {
	keys := make(map[domain.Eduprogcomp]bool)
	var list []domain.Eduprogcomp
	for _, entry := range compSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
