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
			controllers.BadRequest(w, err)
			return
		}

		eduprogschemes, err := c.eduprogschemeService.ShowSchemeByEduprogId(eduprogscheme.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogscheme.EduprogcompId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		totalCompCredits := eduprogscheme.CreditsPerSemester

		for i := range eduprogschemes {
			if eduprogschemes[i].EduprogcompId == eduprogscheme.EduprogcompId {
				totalCompCredits = totalCompCredits + eduprogschemes[i].CreditsPerSemester
				if eduprogschemes[i].SemesterNum == eduprogscheme.SemesterNum {
					log.Printf("EduprogschemeController: %s", err)
					controllers.BadRequest(w, errors.New("this component already exists in this semester"))
					return
				}
			}
		}

		if totalCompCredits > eduprogcomp.Credits {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, errors.New("too much credits per semester, free credits to use left: "))
			return
		}

		eduprogscheme, err = c.eduprogschemeService.SetComponentToEdprogscheme(eduprogscheme)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Created(w, eduprogschemeDto.DomainToDto(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) UpdateComponentInEduprogscheme() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogscheme, err := requests.Bind(r, requests.UpdateComponentInEduprogschemeRequest{}, domain.Eduprogscheme{})
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogscheme, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(eduprogscheme, id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogscheme.EduprogcompId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDto(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.FindById(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogscheme.EduprogcompId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDto(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) FindBySemesterNum() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sNum, err := strconv.ParseUint(chi.URLParam(r, "sNum"), 10, 16)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.FindBySemesterNum(uint16(sNum), id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.ShowListByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) ShowSchemeByEduprogId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "sNum"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.ShowSchemeByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.ShowListByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) ShowFreeComponents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "sNum"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomps, err := c.eduprogcompService.ShowList()
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.ShowSchemeByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var escIds []uint64
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

		uniqes := unique(eduprogcomps)

		var eduprogcompDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompDto.DomainToDtoCollection(uniqes))
	}
}

func (c EduprogschemeController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		err = c.eduprogschemeService.Delete(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		controllers.Ok(w)
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
