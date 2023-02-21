package controllers

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/app"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
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
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
		}

		discipline, err = c.disciplineService.Save(discipline)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
			return
		}

		var disciplineDto resources.DisciplineDto
		Created(w, disciplineDto.DomainToDto(discipline))
	}
}

func (c DisciplineController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
			return
		}

		discipline, err := requests.Bind(r, requests.UpdateDisciplineRequest{}, domain.Discipline{})
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
			return
		}

		discipline, err = c.disciplineService.Update(discipline, id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			InternalServerError(w, err)
			return
		}

		var disciplineDto resources.DisciplineDto
		Created(w, disciplineDto.DomainToDto(discipline))
	}
}

func (c DisciplineController) ShowDisciplinesByEduprogId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
			return
		}

		disciplines, err := c.disciplineService.ShowDisciplinesByEduprogId(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			InternalServerError(w, err)
			return
		}

		//comps, err := c.eduprogcompService.SortComponentsByMnS()
		//if err != nil {
		//	log.Printf("EduprogController: %s", err)
		//	InternalServerError(w, err)
		//	return
		//}
		var disciplineDto resources.DisciplineDto
		Created(w, disciplineDto.DomainToDtoCollection(disciplines))
	}
}

func (c DisciplineController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
			return
		}

		discipline, _ := c.disciplineService.FindById(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
			return
		}

		var disciplineDto resources.DisciplineDto
		Created(w, disciplineDto.DomainToDto(discipline))
	}
}

func (c DisciplineController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
			return
		}

		err = c.disciplineService.Delete(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			InternalServerError(w, err)
			return
		}
		Ok(w)
	}
}