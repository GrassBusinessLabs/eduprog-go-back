package controllers

import (
	"errors"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/app"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
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

func (c EducompRelationsController) CreateRelation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		relation, err := requests.Bind(r, requests.CreateEducompRelationRequest{}, domain.Educomp_relations{})
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			BadRequest(w, errors.New("invalid request body"))
			return
		} else if relation.ChildCompId == relation.BaseCompId {
			log.Printf("EducompRelationsController: %s", err)
			BadRequest(w, errors.New("Base comp cannot be equal to child comp"))
			return
		}

		relation, err = c.educompRelationsService.CreateRelation(relation)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			InternalServerError(w, err)
			return
		}

		var educompRelationsDto resources.EducompRelationsDto
		Created(w, educompRelationsDto.DomainToDto(relation))
	}
}

func (c EducompRelationsController) ShowByEduprogId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			InternalServerError(w, err)
			return
		}

		relations, err := c.educompRelationsService.ShowByEduprogId(id)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			InternalServerError(w, err)
			return
		}

		var educompRelationsDto resources.EducompRelationsDto
		Success(w, educompRelationsDto.DomainToDtoCollection(relations))
	}
}

func (c EducompRelationsController) DeleteByBaseId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			InternalServerError(w, err)
			return
		}

		err = c.educompRelationsService.DeleteByBaseCompId(id)
		if err != nil {
			log.Printf("EducompRelationsController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
