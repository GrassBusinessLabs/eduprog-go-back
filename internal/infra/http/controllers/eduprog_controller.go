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
	eduprogService app.EduprogService
}

func NewEduprogController(es app.EduprogService) EduprogController {
	return EduprogController{
		eduprogService: es,
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
