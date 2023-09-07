package eduprog

import (
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

type CompetenciesMatrixController struct {
	competenciesMatrixService app.CompetenciesMatrixService
}

func NewCompetenciesMatrixController(cms app.CompetenciesMatrixService) CompetenciesMatrixController {
	return CompetenciesMatrixController{
		competenciesMatrixService: cms,
	}
}

func (c CompetenciesMatrixController) CreateRelation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		relation, err := requests.Bind(r, requests.CreateCompetenciesMatrixRelationRequest{}, domain.CompetenciesMatrix{})
		if err != nil {
			log.Printf("CompetenciesMatrixController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		relation, err = c.competenciesMatrixService.CreateRelation(relation)
		if err != nil {
			log.Printf("CompetenciesMatrixController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var competencyMatrixDto resources.CompetenciesMatrixDto
		controllers.Created(w, competencyMatrixDto.DomainToDto(relation))
	}
}

func (c CompetenciesMatrixController) ShowByEduprogId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("CompetenciesMatrixController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		relations, err := c.competenciesMatrixService.ShowByEduprogId(id)
		if err != nil {
			log.Printf("CompetenciesMatrixController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var competencyMatrixDto resources.CompetenciesMatrixDto
		controllers.Created(w, competencyMatrixDto.DomainToDtoCollection(relations))
	}
}

func (c CompetenciesMatrixController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		component_id, err := strconv.ParseUint(chi.URLParam(r, "componentId"), 10, 64)
		if err != nil {
			log.Printf("CompetenciesMatrixController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		competency_id, err := strconv.ParseUint(chi.URLParam(r, "competencyId"), 10, 64)
		if err != nil {
			log.Printf("CompetenciesMatrixController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		err = c.competenciesMatrixService.Delete(component_id, competency_id)
		if err != nil {
			log.Printf("CompetenciesMatrixController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		controllers.Ok(w)
	}
}
