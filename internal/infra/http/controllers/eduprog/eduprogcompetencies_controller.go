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
	competenciesBaseService    app.CompetenciesBaseService
	eduprogService             app.EduprogService
}

func NewEduprogcompetenciesController(ecc app.EduprogcompetenciesService, cbs app.CompetenciesBaseService, es app.EduprogService) EduprogcompetenciesController {
	return EduprogcompetenciesController{
		eduprogcompetenciesService: ecc,
		competenciesBaseService:    cbs,
		eduprogService:             es,
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

		baseCompetency, err := c.competenciesBaseService.FindById(eduprogcompetency.CompetencyId)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		if eduprogcompetency.Definition == "" {
			eduprogcompetency.Definition = baseCompetency.Definition
		}

		eduprogcompetency.Type = baseCompetency.Type

		allEdpcompetencies, err := c.eduprogcompetenciesService.ShowCompetenciesByType(eduprogcompetency.EduprogId, eduprogcompetency.Type)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var maxCode uint64 = 0

		for i := range allEdpcompetencies {
			if allEdpcompetencies[i].CompetencyId == eduprogcompetency.CompetencyId {
				log.Printf("EduprogcompetenciesController: %s", err)
				controllers.InternalServerError(w, errors.New("competency is in this eduprog already"))
				return
			}
			if i == 0 || allEdpcompetencies[i].Code > maxCode {
				maxCode = allEdpcompetencies[i].Code
			}
		}

		eduprogcompetency.Code = maxCode + 1

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

		eduprogcompetencyID, err := c.eduprogcompetenciesService.FindById(id)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcompetency, err := requests.Bind(r, requests.UpdateCompetencyRequest{}, domain.Eduprogcompetencies{})
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcompetency.Id = eduprogcompetencyID.Id
		eduprogcompetency.CompetencyId = eduprogcompetencyID.CompetencyId
		eduprogcompetency.Type = eduprogcompetencyID.Type
		eduprogcompetency.Code = eduprogcompetencyID.Code
		eduprogcompetency.EduprogId = eduprogcompetencyID.EduprogId

		eduprogcompetency, err = c.eduprogcompetenciesService.UpdateCompetency(eduprogcompetency, id)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompetenciesDto resources.EduprogcompetenciesDto
		controllers.Created(w, eduprogcompetenciesDto.DomainToDto(eduprogcompetency))
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

		eduprogcompetency.CompetencyId = 65

		allEdpcompetencies, err := c.eduprogcompetenciesService.ShowCompetenciesByType(eduprogcompetency.EduprogId, eduprogcompetency.Type)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var maxCode uint64 = 0

		for i := range allEdpcompetencies {
			if allEdpcompetencies[i].EduprogId == eduprogcompetency.EduprogId {
				if i == 0 || allEdpcompetencies[i].Code > maxCode {
					maxCode = allEdpcompetencies[i].Code
				}
			}
		}

		eduprogcompetency.Code = maxCode + 1

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

		eduprog, err := c.eduprogService.FindById(id)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcompetencies, err := c.eduprogcompetenciesService.ShowCompetenciesByType(id, ttype)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for i := range eduprogcompetencies {
			err := c.eduprogcompetenciesService.Delete(eduprogcompetencies[i].Id)
			if err != nil {
				log.Printf("EduprogcompetenciesController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
		}

		baseCompetencies, err := c.competenciesBaseService.ShowCompetenciesByType(ttype, eduprog.SpecialtyCode)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompetenciesList []domain.Eduprogcompetencies

		for i := range baseCompetencies {
			var eduprogcompetency domain.Eduprogcompetencies

			eduprogcompetency.CompetencyId = baseCompetencies[i].Id
			eduprogcompetency.EduprogId = id
			eduprogcompetency.Type = baseCompetencies[i].Type
			eduprogcompetency.Code = baseCompetencies[i].Code
			eduprogcompetency.Definition = baseCompetencies[i].Definition

			eduprogcompetency, err = c.eduprogcompetenciesService.AddCompetencyToEduprog(eduprogcompetency)
			if err != nil {
				log.Printf("EduprogcompetenciesController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
			//competencyId,err = c.eduprogcompetenciesService.FindById(eduprogcompetency.Id)

			eduprogcompetenciesList = append(eduprogcompetenciesList, eduprogcompetency)
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

		if ttype != "ZK" && ttype != "FK" && ttype != "PR" {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.BadRequest(w, errors.New("only ZK, FK or PR"))
			return
		}

		eduprogcompetencies, err := c.eduprogcompetenciesService.ShowCompetenciesByType(id, ttype)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for i := range eduprogcompetencies {
			err := c.eduprogcompetenciesService.Delete(eduprogcompetencies[i].Id)
			if err != nil {
				log.Printf("EduprogcompetenciesController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
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

		competency, err := c.eduprogcompetenciesService.FindById(id)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		err = c.eduprogcompetenciesService.Delete(id)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		allEdpcompetencies, err := c.eduprogcompetenciesService.ShowCompetenciesByType(competency.EduprogId, competency.Type)
		if err != nil {
			log.Printf("EduprogcompetenciesController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for i := range allEdpcompetencies {
			if allEdpcompetencies[i].Code > competency.Code {
				allEdpcompetencies[i].Code = allEdpcompetencies[i].Code - 1
				_, _ = c.eduprogcompetenciesService.UpdateCompetency(allEdpcompetencies[i], allEdpcompetencies[i].Id)
				if err != nil {
					log.Printf("EduprogcompetenciesController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
			}

		}

		controllers.Ok(w)
	}
}
