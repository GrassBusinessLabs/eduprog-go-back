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
			controllers.BadRequest(w, err)
			return
		}

		u := r.Context().Value(controllers.UserKey).(domain.User)

		eduprog, err = c.eduprogService.Save(eduprog, u.Id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogDto resources.EduprogDto
		controllers.Created(w, eduprogDto.DomainToDto(eduprog))
	}
}

func (c EduprogController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		req, err := requests.Bind(r, requests.UpdateEduprogRequest{}, domain.Eduprog{})
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		ref, _, err := c.eduprogService.FindById(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		ref, err = c.eduprogService.Update(ref, req)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogDto resources.EduprogDto
		controllers.Success(w, eduprogDto.DomainToDto(ref))
	}
}

func (c EduprogController) GetOPPLevelsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		levels, err := c.eduprogService.GetOPPLevelsList()
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogDto resources.EduprogDto
		controllers.Success(w, eduprogDto.OPPLevelDomainToDtoCollection(levels))
	}
}

func (c EduprogController) ShowList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogs, err := c.eduprogService.ShowList()
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogsDto resources.EduprogDto
		controllers.Success(w, eduprogsDto.DomainToDtoCollection(eduprogs))
	}
}

func (c EduprogController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprog, comps, err := c.eduprogService.FindById(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		comps.Mandatory = c.eduprogService.SortByCode(comps.Mandatory)

		var eduprogDto resources.EduprogDto
		controllers.Success(w, eduprogDto.DomainToDtoWithComps(eduprog, comps, comps.Selective))
	}
}

func (c EduprogController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		err = c.eduprogService.Delete(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		controllers.Ok(w)
	}
}

func (c EduprogController) CreditsInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		creditsDto, err := c.eduprogService.GetCreditsInfo(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		controllers.Success(w, creditsDto)
	}
}

func (c EduprogController) CreateDuplicateOf() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogAdditionalData, err := requests.Bind(r, requests.DuplicateEduprogRequest{}, domain.Eduprog{})
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		u := r.Context().Value(controllers.UserKey).(domain.User)

		eduprog, err := c.eduprogService.CreateDuplicateOf(id, u.Id, eduprogAdditionalData.Name, eduprogAdditionalData.ApprovalYear)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogDto resources.EduprogDto
		controllers.Created(w, eduprogDto.DomainToDto(eduprog))
	}
}
