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

type EduprogresultsController struct {
	eduprogresultsService app.EduprogresultsService
}

func NewEduprogresultsController(ers app.EduprogresultsService) EduprogresultsController {
	return EduprogresultsController{
		eduprogresultsService: ers,
	}
}

func (c EduprogresultsController) AddEduprogresultToEduprog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogresult, err := requests.Bind(r, requests.AddEduprogresultToEduprogRequest{}, domain.Eduprogresult{})
		if err != nil {
			log.Printf("EduprogresultController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogresult.Type = "ПР"

		allEduprogresults, err := c.eduprogresultsService.ShowEduprogResultsByEduprogId(eduprogresult.EduprogId)
		if err != nil {
			log.Printf("EduprogresultController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var maxCode uint64 = 0

		for i := range allEduprogresults {
			if allEduprogresults[i].Type == eduprogresult.Type {
				if i == 0 || allEduprogresults[i].Code > maxCode {
					maxCode = allEduprogresults[i].Code
				}
			}

		}

		eduprogresult.Code = maxCode + 1

		eduprogresult, err = c.eduprogresultsService.AddEduprogresultToEduprog(eduprogresult)
		if err != nil {
			log.Printf("EduprogresultController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogresultsDto resources.EduprogresultsDto
		controllers.Created(w, eduprogresultsDto.DomainToDto(eduprogresult))
	}
}

func (c EduprogresultsController) ShowEduprogResultsByEduprogId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EduprogresultController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogresults, err := c.eduprogresultsService.ShowEduprogResultsByEduprogId(id)
		if err != nil {
			log.Printf("EduprogresultController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogresultsDto resources.EduprogresultsDto
		controllers.Success(w, eduprogresultsDto.DomainToDtoCollection(eduprogresults))
	}
}

func (c EduprogresultsController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "resId"), 10, 64)
		if err != nil {
			log.Printf("EduprogresultController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogresult, err := c.eduprogresultsService.FindById(id)
		if err != nil {
			log.Printf("EduprogresultController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogresultsDto resources.EduprogresultsDto
		controllers.Success(w, eduprogresultsDto.DomainToDto(eduprogresult))
	}
}

func (c EduprogresultsController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "resId"), 10, 64)
		if err != nil {
			log.Printf("EduprogresultController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogresult, err := c.eduprogresultsService.FindById(id)
		if err != nil {
			log.Printf("EduprogresultController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		err = c.eduprogresultsService.Delete(id)
		if err != nil {
			log.Printf("EduprogresultController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		allEduprogresults, _ := c.eduprogresultsService.ShowEduprogResultsByEduprogId(eduprogresult.EduprogId)
		if err != nil {
			log.Printf("EduprogresultController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for i := range allEduprogresults {
			if allEduprogresults[i].Type == eduprogresult.Type {
				if allEduprogresults[i].Code > eduprogresult.Code {
					allEduprogresults[i].Code = allEduprogresults[i].Code - 1
					_, _ = c.eduprogresultsService.UpdateEduprogresult(allEduprogresults[i], allEduprogresults[i].Id)
					if err != nil {
						log.Printf("EduprogresultController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
				}
			}

		}

		controllers.Ok(w)
	}
}
