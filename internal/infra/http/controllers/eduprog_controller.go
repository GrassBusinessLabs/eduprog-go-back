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

type EduprogController struct {
	eduprogService     app.EduprogService
	eduprogcompService app.EduprogcompService
}

func NewEduprogController(es app.EduprogService, ecs app.EduprogcompService) EduprogController {
	return EduprogController{
		eduprogService:     es,
		eduprogcompService: ecs,
	}
}

func (c EduprogController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		eduprog, err := requests.Bind(r, requests.CreateEduprogRequest{}, domain.Eduprog{})
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
		}
		u := r.Context().Value(UserKey).(domain.User)
		eduprog.UserId = u.Id
		eduprog, err = c.eduprogService.Save(eduprog)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
			return
		}

		var eduprogDto resources.EduprogDto
		Created(w, eduprogDto.DomainToDto(eduprog))
	}
}

func (c EduprogController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
			return
		}

		eduprog, err := requests.Bind(r, requests.UpdateEduprogRequest{}, domain.Eduprog{})
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
			return
		}
		u := r.Context().Value(UserKey).(domain.User)
		eduprog.UserId = u.Id
		eduprog, err = c.eduprogService.Update(eduprog, id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			InternalServerError(w, err)
			return
		}

		var eduprogDto resources.EduprogDto
		Success(w, eduprogDto.DomainToDto(eduprog))
	}
}

func (c EduprogController) ShowList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogs, err := c.eduprogService.ShowList()
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

		var eduprogsDto resources.EduprogDto
		Success(w, eduprogsDto.DomainToDtoCollection(eduprogs))
	}
}

func (c EduprogController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
			return
		}

		eduprog, _ := c.eduprogService.FindById(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
			return
		}

		comps, err := c.eduprogcompService.SortComponentsByMnS(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			InternalServerError(w, err)
			return
		}

		var eduprogDto resources.EduprogDto
		var compsDto resources.EduprogcompDto
		Success(w, eduprogDto.DomainToDtoWithComps(eduprog, compsDto.DomainToDtoWCompCollection(comps)))
	}
}

func (c EduprogController) CreditsInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
			return
		}

		comps, err := c.eduprogcompService.SortComponentsByMnS(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			InternalServerError(w, err)
			return
		}

		var creditsDto resources.CreditsDto

		for _, comp := range comps.Selective {
			creditsDto.SelectiveCredits += comp.Credits
		}
		for _, comp := range comps.Mandatory {
			creditsDto.MandatoryCredits += comp.Credits
		}
		creditsDto.TotalCredits = creditsDto.SelectiveCredits + creditsDto.MandatoryCredits
		creditsDto.TotalFreeCredits = 240 - creditsDto.TotalCredits
		creditsDto.MandatoryFreeCredits = 180 - creditsDto.MandatoryCredits
		creditsDto.SelectiveFreeCredits = 60 - creditsDto.SelectiveCredits

		Success(w, creditsDto)
	}
}

func (c EduprogController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			BadRequest(w, err)
			return
		}

		err = c.eduprogService.Delete(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			InternalServerError(w, err)
			return
		}
		Ok(w)
	}
}
