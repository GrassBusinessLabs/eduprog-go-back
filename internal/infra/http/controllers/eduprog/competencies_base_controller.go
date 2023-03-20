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

type CompetenciesBaseController struct {
	competenciesBaseService app.CompetenciesBaseService
}

func NewCompetenciesBaseController(cbs app.CompetenciesBaseService) CompetenciesBaseController {
	return CompetenciesBaseController{
		competenciesBaseService: cbs,
	}
}

func (c CompetenciesBaseController) CreateCompetencyBase() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		competencyBase, err := requests.Bind(r, requests.CreateCompetencyBaseRequest{}, domain.CompetenciesBase{})
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		if competencyBase.Type != "ЗК" && competencyBase.Type != "ФК" && competencyBase.Type != "ПР" {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.BadRequest(w, errors.New("only 'ЗК', 'ФК' or 'ПР'"))
			return
		}

		specialty, err := strconv.ParseUint(competencyBase.Specialty, 10, 64)
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.BadRequest(w, err)
			return
		}
		if specialty < 11 || specialty > 293 {
			controllers.BadRequest(w, errors.New("from 11 to 293"))
			return
		}

		allCompetencies, err := c.competenciesBaseService.ShowCompetenciesByType(competencyBase.Type, competencyBase.Specialty)
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var maxCode uint64 = 0

		for i := range allCompetencies {
			if i == 0 || allCompetencies[i].Code > maxCode {
				maxCode = allCompetencies[i].Code
			}
		}

		competencyBase.Code = maxCode + 1

		competencyBase, err = c.competenciesBaseService.CreateCompetency(competencyBase)
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var competenciesBaseDto resources.CompetenciesBaseDto
		controllers.Created(w, competenciesBaseDto.DomainToDto(competencyBase))
	}
}

func (c CompetenciesBaseController) UpdateCompetencyBase() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "cbId"), 10, 64)
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		competencyBase, err := requests.Bind(r, requests.CreateCompetencyBaseRequest{}, domain.CompetenciesBase{})
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		if competencyBase.Type != "ЗК" && competencyBase.Type != "ФК" && competencyBase.Type != "ПР" {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.BadRequest(w, errors.New("only 'ЗК', 'ФК' or 'ПР'"))
			return
		}

		specialty, err := strconv.ParseUint(competencyBase.Specialty, 10, 64)
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.BadRequest(w, err)
			return
		}
		if specialty < 11 || specialty > 293 {
			controllers.BadRequest(w, errors.New("from 11 to 293"))
			return
		}

		competencyBase, err = c.competenciesBaseService.UpdateCompetency(competencyBase, id)
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		competencyBase.Id = id
		var competenciesBaseDto resources.CompetenciesBaseDto
		controllers.Created(w, competenciesBaseDto.DomainToDto(competencyBase))
	}
}

func (c CompetenciesBaseController) ShowAllCompetencies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		competencies, err := c.competenciesBaseService.ShowAllCompetencies()
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var competenciesDto resources.CompetenciesBaseDto
		controllers.Success(w, competenciesDto.DomainToDtoCollection(competencies))
	}
}

func (c CompetenciesBaseController) ShowCompetenciesByType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ttype := r.URL.Query().Get("type")
		if ttype != "ZK" && ttype != "FK" && ttype != "PR" {
			controllers.BadRequest(w, errors.New("only ZK, FK or PR"))
			return
		}

		specialty, err := strconv.ParseInt(r.URL.Query().Get("specialty"), 10, 64)
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.BadRequest(w, err)
			return
		}
		if specialty < 11 || specialty > 293 {
			controllers.BadRequest(w, errors.New("only ZK, FK or PR"))
			return
		}

		competencies, err := c.competenciesBaseService.ShowCompetenciesByType(ttype, strconv.FormatInt(specialty, 10))
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var competenciesDto resources.CompetenciesBaseDto
		controllers.Success(w, competenciesDto.DomainToDtoCollection(competencies))
	}
}

func (c CompetenciesBaseController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "cbId"), 10, 64)
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		competency, err := c.competenciesBaseService.FindById(id)
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		var competenciesDto resources.CompetenciesBaseDto
		controllers.Success(w, competenciesDto.DomainToDto(competency))
	}
}

func (c CompetenciesBaseController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "cbId"), 10, 64)
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		competency, err := c.competenciesBaseService.FindById(id)
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		err = c.competenciesBaseService.Delete(id)
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		allCompetencies, err := c.competenciesBaseService.ShowCompetenciesByType(competency.Type, competency.Specialty)
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for i := range allCompetencies {
			if allCompetencies[i].Code > competency.Code {
				allCompetencies[i].Code = allCompetencies[i].Code - 1
				_, _ = c.competenciesBaseService.UpdateCompetency(allCompetencies[i], allCompetencies[i].Id)
				if err != nil {
					log.Printf("CompetenciesBaseController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
			}
		}

		controllers.Ok(w)
	}
}
