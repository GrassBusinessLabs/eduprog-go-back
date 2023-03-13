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

type ResultsMatrixController struct {
	resultsMatrixService app.ResultsMatrixService
}

func NewResultsMatrixController(rms app.ResultsMatrixService) ResultsMatrixController {
	return ResultsMatrixController{
		resultsMatrixService: rms,
	}
}

func (c ResultsMatrixController) CreateRelation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		relation, err := requests.Bind(r, requests.CreateResultsMatrixRelationRequest{}, domain.ResultsMatrix{})
		if err != nil {
			log.Printf("ResultsMatrixController: %s", err)
			controllers.BadRequest(w, errors.New("invalid request body"))
			return
		}

		relation, err = c.resultsMatrixService.CreateRelation(relation)
		if err != nil {
			log.Printf("ResultsMatrixController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var resultMatrixDto resources.ResultsMatrixDto
		controllers.Created(w, resultMatrixDto.DomainToDto(relation))
	}
}

func (c ResultsMatrixController) ShowByEduprogId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("ResultsMatrixController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		relations, err := c.resultsMatrixService.ShowByEduprogId(id)
		if err != nil {
			log.Printf("ResultsMatrixController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var resultMatrixDto resources.ResultsMatrixDto
		controllers.Created(w, resultMatrixDto.DomainToDtoCollection(relations))
	}
}

func (c ResultsMatrixController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		component_id, err := strconv.ParseUint(chi.URLParam(r, "componentId"), 10, 64)
		if err != nil {
			log.Printf("ResultsMatrixController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogresult_id, err := strconv.ParseUint(chi.URLParam(r, "edresultId"), 10, 64)
		if err != nil {
			log.Printf("ResultsMatrixController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		err = c.resultsMatrixService.Delete(component_id, eduprogresult_id)
		if err != nil {
			log.Printf("ResultsMatrixController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		controllers.Ok(w)
	}
}
