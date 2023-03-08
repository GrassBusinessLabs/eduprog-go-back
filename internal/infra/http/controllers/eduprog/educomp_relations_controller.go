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

type EducompRelationsController struct {
	educompRelationsService app.EducompRelationsService
	eduprogschemeService    app.EduprogschemeService
	eduprogcompService      app.EduprogcompService
}

func NewEducompRelationsController(ecrs app.EducompRelationsService, epss app.EduprogschemeService, epcs app.EduprogcompService) EducompRelationsController {
	return EducompRelationsController{
		educompRelationsService: ecrs,
		eduprogschemeService:    epss,
		eduprogcompService:      epcs,
	}
}

func (c EducompRelationsController) CreateRelation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		relation, err := requests.Bind(r, requests.CreateEducompRelationRequest{}, domain.Educomp_relations{})
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.BadRequest(w, errors.New("invalid request body"))
			return
		} else if relation.ChildCompId == relation.BaseCompId {
			log.Printf("EducompRelationsController: %s", err)
			controllers.BadRequest(w, errors.New("Base comp cannot be equal to child comp"))
			return
		}

		//baseComp, _ := c.eduprogcompService.FindByWODeleteDate(relation.BaseCompId)
		//childComp, _ := c.eduprogcompService.FindByWODeleteDate(relation.ChildCompId)
		//if baseComp.DeletedDate != nil || childComp.DeletedDate != nil {
		//	log.Printf("EducompRelationsController: %s", err)
		//	controllers.BadRequest(w, errors.New("Base or child comp dont exist"))
		//	return
		//}

		relation, err = c.educompRelationsService.CreateRelation(relation)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var educompRelationsDto resources.EducompRelationsDto
		controllers.Created(w, educompRelationsDto.DomainToDto(relation))
	}
}

func (c EducompRelationsController) ShowByEduprogId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		relations, err := c.educompRelationsService.ShowByEduprogId(id)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		//var result []domain.Educomp_relations
		//for i := range relations {
		//	baseComp, _ := c.eduprogcompService.FindByWODeleteDate(relations[i].BaseCompId)
		//	childComp, _ := c.eduprogcompService.FindByWODeleteDate(relations[i].ChildCompId)
		//	if baseComp.DeletedDate == nil && childComp.DeletedDate == nil {
		//		result = append(result, relations[i])
		//	}
		//}

		var educompRelationsDto resources.EducompRelationsDto
		controllers.Success(w, educompRelationsDto.DomainToDtoCollection(relations))
	}
}

func (c EducompRelationsController) ShowPossibleRelationsForComp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		edId, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}
		compId, _ := strconv.ParseUint(chi.URLParam(r, "compId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogscheme, _ := c.eduprogschemeService.ShowSchemeByEduprogId(edId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomps, _ := c.eduprogcompService.ShowListByEduprogId(edId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var result []domain.Eduprogcomp
		var maxCompSemester uint64 = 0
		for i := range eduprogscheme {
			if eduprogscheme[i].EduprogcompId == compId {
				if maxCompSemester < eduprogscheme[i].SemesterNum {
					maxCompSemester = eduprogscheme[i].SemesterNum
				}
			}
		}

		for i := range eduprogscheme {
			if eduprogscheme[i].EduprogcompId == compId {

				for i2 := range eduprogscheme {
					if eduprogscheme[i2].SemesterNum > maxCompSemester {
						for i3 := range eduprogcomps {
							if eduprogcomps[i3].Id == eduprogscheme[i2].EduprogcompId {
								result = append(result, eduprogcomps[i3])
							}
						}
					}
				}
			}
		}

		uniqes := unique(result)

		var eduprogcompDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompDto.DomainToDtoCollection(uniqes))
	}
}

func (c EducompRelationsController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		base_id, err := strconv.ParseUint(chi.URLParam(r, "baseId"), 10, 64)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		child_id, err := strconv.ParseUint(chi.URLParam(r, "childId"), 10, 64)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		err = c.educompRelationsService.Delete(base_id, child_id)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		controllers.Ok(w)
	}
}
