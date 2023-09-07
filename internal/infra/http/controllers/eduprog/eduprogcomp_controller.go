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

type EduprogcompController struct {
	eduprogcompService app.EduprogcompService
}

func NewEduprogcompController(es app.EduprogcompService) EduprogcompController {
	return EduprogcompController{
		eduprogcompService: es,
	}
}

func (c EduprogcompController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogcomp, err := requests.Bind(r, requests.CreateEduprogcompRequest{}, domain.Eduprogcomp{})
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomp, err = c.eduprogcompService.Save(eduprogcomp)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var eduprogcompDto resources.EduprogcompDto
		controllers.Created(w, eduprogcompDto.DomainToDto(eduprogcomp))
	}
}

func (c EduprogcompController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epcId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		req, err := requests.Bind(r, requests.UpdateEduprogcompRequest{}, domain.Eduprogcomp{})
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		ref, err := c.eduprogcompService.FindById(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		req, err = c.eduprogcompService.Update(ref, req)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompDto.DomainToDto(req))
	}
}

func (c EduprogcompController) ReplaceComp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		educompId, err := strconv.ParseUint(r.URL.Query().Get("edcompId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, errors.New("problems in parsing 'edcompId' query parameter"))
			return
		}

		putAfterCode, err := strconv.ParseUint(r.URL.Query().Get("putAfter"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, errors.New("problems in parsing 'putAfter' query parameter"))
			return
		}

		eduprogcomps, err := c.eduprogcompService.ReplaceOK(educompId, putAfterCode)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, errors.New("problems in parsing 'putAfter' query parameter"))
			return
		}

		var eduprogcompDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompDto.DomainToDtoWCompCollection(eduprogcomps, eduprogcomps.Selective))
	}
}

func (c EduprogcompController) ReplaceVB() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		educompId, err := strconv.ParseUint(r.URL.Query().Get("edcompId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		blockNum, err := strconv.ParseUint(r.URL.Query().Get("blockNum"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		putAfterCode, err := strconv.ParseUint(r.URL.Query().Get("putAfter"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomps, err := c.eduprogcompService.ReplaceVB(educompId, blockNum, putAfterCode)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var eduprogcompDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompDto.DomainToDtoWCompCollection(eduprogcomps, eduprogcomps.Selective))
	}
}

func (c EduprogcompController) ReplaceCompsBlock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		firstEducompId, err := strconv.ParseUint(r.URL.Query().Get("firstEducompId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		putAfterCode, err := strconv.ParseUint(r.URL.Query().Get("putAfter"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, errors.New("error while parsing 'putAfter' query param"))
			return
		}

		eduprogcomps, err := c.eduprogcompService.ReplaceVBBlock(firstEducompId, putAfterCode)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompDto.DomainToDtoWCompCollection(eduprogcomps, eduprogcomps.Selective))
	}
}

func (c EduprogcompController) ReplaceCompBySendingSlice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogcomps, err := requests.Bind(r, requests.SendEduprogcompSliceRequest{}, []domain.Eduprogcomp{})
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		for i := range eduprogcomps {
			eduprogcomps[i], err = c.eduprogcompService.FindById(eduprogcomps[i].Id)
			if err != nil {
				log.Printf("EduprogcompController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
			eduprogcomps[i].Code = strconv.Itoa(i + 1)
			_, err = c.eduprogcompService.Update(eduprogcomps[i], eduprogcomps[i])
			if err != nil {
				log.Printf("EduprogcompController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
		}

		var eduprogcompDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompDto.DomainToDtoCollection(eduprogcomps))
	}
}

func (c EduprogcompController) UpdateVBName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogId, err := strconv.ParseUint(chi.URLParam(r, "epcId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcompReq, err := requests.Bind(r, requests.UpdateBlockName{}, domain.Eduprogcomp{})
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		result, err := c.eduprogcompService.UpdateVBName(eduprogId, eduprogcompReq)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, errors.New("problems in parsing 'putAfter' query parameter"))
			return
		}

		var eduprogcompsDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompsDto.DomainToDtoCollection(result))
	}
}

func (c EduprogcompController) GetVBBlocksInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epcId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomps, err := c.eduprogcompService.SortComponentsByMnS(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompsDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompsDto.BlockInfoToDtoCollection(eduprogcomps.Selective))
	}
}

func (c EduprogcompController) ShowListByEduprogId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epcId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomps, err := c.eduprogcompService.ShowListByEduprogId(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogcompsDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompsDto.DomainToDtoCollection(eduprogcomps))
	}
}

func (c EduprogcompController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epcId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var eduprogcompDto resources.EduprogcompDto
		controllers.Success(w, eduprogcompDto.DomainToDto(eduprogcomp))
	}
}

func (c EduprogcompController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epcId"), 10, 64)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		err = c.eduprogcompService.Delete(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, errors.New("problems in parsing 'putAfter' query parameter"))
			return
		}

		controllers.Ok(w)
	}
}
