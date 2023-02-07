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

type EduprogschemeController struct {
	eduprogschemeService app.EduprogschemeService
	eduprogcompService   app.EduprogcompService
}

func NewEduprogschemeController(ess app.EduprogschemeService, ecs app.EduprogcompService) EduprogschemeController {
	return EduprogschemeController{
		eduprogschemeService: ess,
		eduprogcompService:   ecs,
	}
}

func (c EduprogschemeController) SetComponentToEdprogscheme() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		eduprogscheme, err := requests.Bind(r, requests.SetComponentToEdprogschemeRequest{}, domain.Eduprogscheme{})
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			BadRequest(w, err)
		}

		eduprogscheme, err = c.eduprogschemeService.SetComponentToEdprogscheme(eduprogscheme)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			BadRequest(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogscheme.EduprogcompId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			BadRequest(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		Created(w, eduprogschemeDto.DomainToDto(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) UpdateComponentInEduprogscheme() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			BadRequest(w, err)
			return
		}

		eduprogscheme, err := requests.Bind(r, requests.UpdateComponentInEduprogschemeRequest{}, domain.Eduprogscheme{})
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			BadRequest(w, err)
			return
		}

		eduprogscheme, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(eduprogscheme, id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			InternalServerError(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogscheme.EduprogcompId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			BadRequest(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		Success(w, eduprogschemeDto.DomainToDto(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			BadRequest(w, err)
			return
		}

		eduprogscheme, _ := c.eduprogschemeService.FindById(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			BadRequest(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogscheme.EduprogcompId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			BadRequest(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		Success(w, eduprogschemeDto.DomainToDto(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) FindBySemesterNum() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "sNum"), 10, 16)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			BadRequest(w, err)
			return
		}

		eduprogscheme, _ := c.eduprogschemeService.FindBySemesterNum(uint16(id))
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			BadRequest(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.ShowListByEduprogId(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			BadRequest(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) ShowSchemeByEduprogId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "sNum"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			BadRequest(w, err)
			return
		}

		eduprogscheme, _ := c.eduprogschemeService.ShowSchemeByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			BadRequest(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.ShowListByEduprogId(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			BadRequest(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			BadRequest(w, err)
			return
		}

		err = c.eduprogschemeService.Delete(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
