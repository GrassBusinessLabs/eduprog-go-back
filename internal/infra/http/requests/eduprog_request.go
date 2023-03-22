package requests

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
)

type CreateEduprogRequest struct {
	Name           string `json:"name" validate:"required,gte=1,max=50"`
	EducationLevel string `json:"education_level" validate:"required,gte=1,max=50"`
	SpecialityCode string `json:"speciality_code" validate:"required,gte=1,max=3"`
}

type UpdateEduprogRequest struct {
	Name           string `json:"name" validate:"required,gte=1,max=50"`
	EducationLevel string `json:"education_level" validate:"required,gte=1,max=50"`
	SpecialityCode string `json:"speciality_code" validate:"required,gte=1,max=3"`
}

func (r CreateEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprog{
		Name:           r.Name,
		EducationLevel: r.EducationLevel,
		SpecialtyCode:  r.SpecialityCode,
	}, nil
}

func (r UpdateEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprog{
		Name:           r.Name,
		EducationLevel: r.EducationLevel,
		SpecialtyCode:  r.SpecialityCode,
	}, nil
}

//type UpdateEduprogRequest struct {
//
//}
