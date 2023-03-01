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

type EduprogcompetenciesController struct {
	eduprogcompetenciesService app.EduprogcompetenciesService
	competenciesBaseService    app.CompetenciesBaseService
}

func NewEduprogcompetenciesController(ecc app.EduprogcompetenciesService, cbs app.CompetenciesBaseService) EduprogcompetenciesController {
	return EduprogcompetenciesController{
		eduprogcompetenciesService: ecc,
		competenciesBaseService:    cbs,
	}
}

func (c EduprogcompetenciesController) AddCompetencyToEduprog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogcompetency, err := requests.Bind(r, requests.AddCompetencyToEduprogRequest{}, domain.Eduprogcompetencies{})
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		if eduprogcompetency.Redefinition == "" {
			competency, err := c.competenciesBaseService.FindById(eduprogcompetency.CompetencyId)
			if err != nil {
				log.Printf("EduprogcompetenciesController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
			eduprogcompetency.Redefinition = competency.Definition
		}

		eduprogcompetency, err = c.eduprogcompetenciesService.AddCompetencyToEduprog(eduprogcompetency)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompetenciesDto resources.EduprogcompetenciesDto
		controllers.Created(w, eduprogcompetenciesDto.DomainToDto(eduprogcompetency))
	}
}

func (c EduprogcompetenciesController) ShowCompetenciesByEduprogId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcompetencies, _ := c.eduprogcompetenciesService.ShowCompetenciesByEduprogId(id)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompetenciesDto resources.EduprogcompetenciesDto
		controllers.Success(w, eduprogcompetenciesDto.DomainToDtoCollection(eduprogcompetencies))
	}
}

func (c EduprogcompetenciesController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "compId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcompetency, _ := c.eduprogcompetenciesService.FindById(id)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompetenciesDto resources.EduprogcompetenciesDto
		controllers.Success(w, eduprogcompetenciesDto.DomainToDto(eduprogcompetency))
	}
}

func (c EduprogcompetenciesController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "compId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		_ = c.eduprogcompetenciesService.Delete(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

	}
}
