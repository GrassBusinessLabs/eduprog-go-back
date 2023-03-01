package eduprog

import (
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
			log.Printf("CompetenciesBase controller: %s", err)
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
			log.Printf("DisciplineController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		competency, _ := c.competenciesBaseService.FindById(id)
		if err != nil {
			log.Printf("CompetenciesBase controller: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		var competenciesDto resources.CompetenciesBaseDto
		controllers.Success(w, competenciesDto.DomainToDto(competency))
	}
}
