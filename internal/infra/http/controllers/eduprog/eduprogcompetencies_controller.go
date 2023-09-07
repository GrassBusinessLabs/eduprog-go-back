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

type EduprogcompetenciesController struct {
	eduprogcompetenciesService app.EduprogcompetenciesService
}

func NewEduprogcompetenciesController(ecc app.EduprogcompetenciesService) EduprogcompetenciesController {
	return EduprogcompetenciesController{
		eduprogcompetenciesService: ecc,
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

func (c EduprogcompetenciesController) UpdateCompetency() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "compId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		req, err := requests.Bind(r, requests.UpdateCompetencyRequest{}, domain.Eduprogcompetencies{})
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		ref, err := c.eduprogcompetenciesService.FindById(id)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		ref, err = c.eduprogcompetenciesService.UpdateCompetency(ref, req)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompetenciesDto resources.EduprogcompetenciesDto
		controllers.Created(w, eduprogcompetenciesDto.DomainToDto(ref))
	}
}

func (c EduprogcompetenciesController) AddCustomCompetencyToEduprog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogcompetency, err := requests.Bind(r, requests.AddCustomCompetencyToEduprogRequest{}, domain.Eduprogcompetencies{})
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcompetency, err = c.eduprogcompetenciesService.AddCustomCompetecyToEduprog(eduprogcompetency)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompetenciesDto resources.EduprogcompetenciesDto
		controllers.Created(w, eduprogcompetenciesDto.DomainToDto(eduprogcompetency))
	}
}

func (c EduprogcompetenciesController) AddAllCompetencies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		ttype := r.URL.Query().Get("type")
		if ttype != "ZK" && ttype != "FK" && ttype != "PR" {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, errors.New("only ZK, FK or PR"))
			return
		}

		eduprogcompetenciesList, err := c.eduprogcompetenciesService.AddAllCompetencies(id, ttype)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompetenciesDto resources.EduprogcompetenciesDto
		controllers.Success(w, eduprogcompetenciesDto.DomainToDtoCollection(eduprogcompetenciesList))

	}
}

func (c EduprogcompetenciesController) DeleteAllCompetencies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		ttype := r.URL.Query().Get("type")

		err = c.eduprogcompetenciesService.DeleteAllCompetencies(id, ttype)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		controllers.Ok(w)
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

		eduprogcompetencies, err := c.eduprogcompetenciesService.ShowCompetenciesByEduprogId(id)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompetenciesDto resources.EduprogcompetenciesDto
		controllers.Success(w, eduprogcompetenciesDto.DomainToDtoCollection(eduprogcompetencies))
	}
}

func (c EduprogcompetenciesController) ShowCompetenciesByType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		ttype := r.URL.Query().Get("type")
		if ttype != "ZK" && ttype != "FK" && ttype != "PR" && ttype != "VFK" && ttype != "VPR" {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, errors.New("only ZK, FK, PR, VFK or VPR"))
			return
		}

		eduprogcompetencies, err := c.eduprogcompetenciesService.ShowCompetenciesByType(id, ttype)
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

		eduprogcompetency, err := c.eduprogcompetenciesService.FindById(id)
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
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		err = c.eduprogcompetenciesService.Delete(id)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		controllers.Ok(w)
	}
}
