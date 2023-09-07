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

type DisciplineController struct {
	disciplineService app.DisciplineService
}

func NewDisciplineController(ds app.DisciplineService) DisciplineController {
	return DisciplineController{
		disciplineService: ds,
	}
}

func (c DisciplineController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		discipline, err := requests.Bind(r, requests.CreateDisciplineRequest{}, domain.Discipline{})
		if err != nil {
			log.Printf("DisciplineController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		discipline, err = c.disciplineService.Save(discipline)
		if err != nil {
			log.Printf("DisciplineController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var disciplineDto resources.DisciplineDto
		controllers.Created(w, disciplineDto.DomainToDto(discipline))
	}
}

func (c DisciplineController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("DisciplineController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		req, err := requests.Bind(r, requests.UpdateDisciplineRequest{}, domain.Discipline{})
		if err != nil {
			log.Printf("DisciplineController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		ref, err := c.disciplineService.FindById(id)
		if err != nil {
			log.Printf("DisciplineController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		discipline, err := c.disciplineService.Update(ref, req)
		if err != nil {
			log.Printf("DisciplineController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var disciplineDto resources.DisciplineDto
		controllers.Success(w, disciplineDto.DomainToDto(discipline))
	}
}

func (c DisciplineController) AddRow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("DisciplineController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		discipline, err := c.disciplineService.AddRow(id)
		if err != nil {
			log.Printf("DisciplineController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var disciplineDto resources.DisciplineDto
		controllers.Success(w, disciplineDto.DomainToDto(discipline))
	}
}

func (c DisciplineController) ShowDisciplinesByEduprogId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("DisciplineController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		disciplines, err := c.disciplineService.ShowDisciplinesByEduprogId(id)
		if err != nil {
			log.Printf("DisciplineController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var disciplineDto resources.DisciplineDto
		controllers.Success(w, disciplineDto.DomainToDtoCollection(disciplines))
	}
}

func (c DisciplineController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("DisciplineController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		discipline, err := c.disciplineService.FindById(id)
		if err != nil {
			log.Printf("DisciplineController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var disciplineDto resources.DisciplineDto
		controllers.Success(w, disciplineDto.DomainToDto(discipline))
	}
}

func (c DisciplineController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("DisciplineController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		err = c.disciplineService.Delete(id)
		if err != nil {
			log.Printf("DisciplineController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		controllers.Ok(w)
	}
}
