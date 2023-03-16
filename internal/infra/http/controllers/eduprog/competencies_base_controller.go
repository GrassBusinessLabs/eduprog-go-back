package eduprog

import (
	"errors"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/app"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers"
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
		competencies, err := c.competenciesBaseService.ShowCompetenciesByType(ttype)
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
