package eduprog

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/app"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/requests"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/resources"
	"log"
	"net/http"
	"reflect"
)

type SpecialtyController struct {
	specialtiesService app.SpecialtiesService
}

func NewSpecialtiesController(ss app.SpecialtiesService) SpecialtyController {
	return SpecialtyController{
		specialtiesService: ss,
	}
}

func (c SpecialtyController) CreateSpecialty() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		specialty, err := requests.Bind(r, requests.CreateSpecialtyRequest{}, domain.Specialty{})
		if err != nil {
			log.Printf("SpecialtyController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		specialty, err = c.specialtiesService.CreateSpecialty(specialty)
		if err != nil {
			log.Printf("SpecialtyController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var specialtyDto resources.SpecialtyDto
		controllers.Created(w, specialtyDto.DomainToDto(specialty))
	}
}

func (c SpecialtyController) UpdateSpecialty() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("sCode")

		specialty, err := requests.Bind(r, requests.UpdateSpecialtyRequest{}, domain.Specialty{})
		if err != nil {
			log.Printf("SpecialtyController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		specialty, err = c.specialtiesService.UpdateSpecialty(specialty, code)
		if err != nil {
			log.Printf("SpecialtyController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		var specialtyDto resources.SpecialtyDto
		controllers.Success(w, specialtyDto.DomainToDto(specialty))
	}
}

func (c SpecialtyController) ShowAllSpecialties() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		specialties, err := c.specialtiesService.ShowAllSpecialties()
		if err != nil {
			log.Printf("SpecialtyController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var specialtyDto resources.SpecialtyDto
		controllers.Success(w, specialtyDto.DomainToDtoCollection(specialties))
	}
}

func (c SpecialtyController) ShowAllKFs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		specialties, err := c.specialtiesService.ShowAllSpecialties()
		if err != nil {
			log.Printf("SpecialtyController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var specialtyDto resources.SpecialtyDto
		controllers.Success(w, specialtyDto.KnowledgeFieldToDtoCollection(removeDuplicatesByField(specialties, "KFCode")))
	}
}

func removeDuplicatesByField(mySlice []domain.Specialty, fieldName string) []domain.Specialty {
	unique := make(map[string]bool)
	result := make([]domain.Specialty, 0)
	for _, v := range mySlice {
		fieldValue := reflect.ValueOf(v).FieldByName(fieldName).String()
		if !unique[fieldValue] {
			unique[fieldValue] = true
			result = append(result, v)
		}
	}
	return result
}

func (c SpecialtyController) ShowByKFCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kfCode := r.URL.Query().Get("kfCode")

		specialties, err := c.specialtiesService.ShowByKFCode(kfCode)
		if err != nil {
			log.Printf("SpecialtyController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var specialtyDto resources.SpecialtyDto
		controllers.Success(w, specialtyDto.DomainToDtoCollection(specialties))
	}
}

func (c SpecialtyController) FindByCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		specialty, err := c.specialtiesService.FindByCode(code)
		if err != nil {
			log.Printf("SpecialtyController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var specialtyDto resources.SpecialtyDto
		controllers.Success(w, specialtyDto.DomainToDto(specialty))
	}
}

func (c SpecialtyController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		err := c.specialtiesService.Delete(code)
		if err != nil {
			log.Printf("SpecialtyController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		controllers.Ok(w)
	}
}
