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
			return
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

func (c EduprogschemeController) ShowFreeComponents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "sNum"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			BadRequest(w, err)
			return
		}

		eduprogcomps, err := c.eduprogcompService.ShowList()
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			InternalServerError(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.ShowSchemeByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			InternalServerError(w, err)
			return
		}

		//var result []domain.Eduprogcomp
		var escIds []uint64 // 3 4 4 2 3 3
		for i := range eduprogscheme {
			escIds = append(escIds, eduprogscheme[i].EduprogcompId)
		}
		for i := range eduprogcomps {
			for i2 := range escIds {
				if eduprogcomps[i].Id == escIds[i2] {
					remove(eduprogcomps, uint64(i))
				}
			}
		}

		uniqes := unique(eduprogcomps) //not working

		var eduprogcompDto resources.EduprogcompDto
		Success(w, eduprogcompDto.DomainToDtoCollection2(uniqes))
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

func remove(s []domain.Eduprogcomp, i uint64) []domain.Eduprogcomp {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func unique(compSlice []domain.Eduprogcomp) []domain.Eduprogcomp {
	keys := make(map[domain.Eduprogcomp]bool)
	list := []domain.Eduprogcomp{}
	for _, entry := range compSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
