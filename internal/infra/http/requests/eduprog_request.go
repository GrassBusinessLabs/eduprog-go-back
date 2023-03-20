package requests

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
)

type CreateEduprogRequest struct {
	Name           string `json:"name" validate:"required,alphanum,gte=1,max=50"`
	EducationLevel string `json:"education_level" validate:"required,alphanum,gte=1,max=50"`
	Speciality     string `json:"speciality" validate:"required,gte=1,max=100"`
}

type UpdateEduprogRequest struct {
	Name           string `json:"name" validate:"required,alphanum,gte=1,max=50"`
	EducationLevel string `json:"education_level" validate:"required,alphanum,gte=1,max=50"`
	Speciality     string `json:"speciality" validate:"required,gte=1,max=100"`
}

func (r CreateEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprog{
		Name:           r.Name,
		EducationLevel: r.EducationLevel,
		Speciality:     r.Speciality,
	}, nil
}

func (r UpdateEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprog{
		Name:           r.Name,
		EducationLevel: r.EducationLevel,
		Speciality:     r.Speciality,
	}, nil
}

//type UpdateEduprogRequest struct {
//
//}
