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
}

func NewEducompRelationsController(ecrs app.EducompRelationsService) EducompRelationsController {
	return EducompRelationsController{
		educompRelationsService: ecrs,
	}
}

func (c EducompRelationsController) ShowPossibleRelations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogId, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		educompPossibleRelations, err := c.educompRelationsService.ShowPossibleRelations(eduprogId)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var educompWithPossibleRelationsDto resources.EducompWithPossibleRelationsDto
		controllers.Success(w, educompWithPossibleRelationsDto.DomainToDtoCollection(educompPossibleRelations))
	}
}

func (c EducompRelationsController) CreateRelation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		relation, err := requests.Bind(r, requests.CreateEducompRelationRequest{}, domain.EducompRelations{})
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.BadRequest(w, errors.New("invalid request body"))
			return
		} else if relation.ChildCompId == relation.BaseCompId {
			log.Printf("EducompRelationsController: %s", err)
			controllers.BadRequest(w, errors.New("base comp cannot be equal to child comp"))
			return
		}

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

		var educompRelationsDto resources.EducompRelationsDto
		controllers.Success(w, educompRelationsDto.DomainToDtoCollection(relations))
	}
}

func (c EducompRelationsController) ShowPossibleRelationsForComp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogId, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.BadRequest(w, err)
			return
		}
		eduprogcompId, err := strconv.ParseUint(chi.URLParam(r, "compId"), 10, 64)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		possibleEduprogcomps, err := c.educompRelationsService.ShowPossibleRelationsForComp(eduprogId, eduprogcompId)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompDto.DomainToDtoCollection(possibleEduprogcomps))
	}
}

func (c EducompRelationsController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseCompId, err := strconv.ParseUint(chi.URLParam(r, "baseId"), 10, 64)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		childCompId, err := strconv.ParseUint(chi.URLParam(r, "childId"), 10, 64)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		err = c.educompRelationsService.Delete(baseCompId, childCompId)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		controllers.Ok(w)
	}
}
