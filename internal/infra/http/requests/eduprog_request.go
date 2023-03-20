package requests

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
)

type CreateEduprogRequest struct {
	Name           string `json:"name" validate:"required,gte=1,max=40"`
	EducationLevel string `json:"education_level" validate:"required,gte=1,max=40"`
	Speciality     string `json:"speciality" validate:"required,gte=1,max=40"`
	KnowledgeField string `json:"knowledge_field" validate:"required,gte=1,max=40"`
}

type UpdateEduprogRequest struct {
	Name           string `json:"name" validate:"required,gte=1,max=40"`
	EducationLevel string `json:"education_level" validate:"required,gte=1,max=40"`
	Speciality     string `json:"speciality" validate:"required,gte=1,max=40"`
	KnowledgeField string `json:"knowledge_field" validate:"required,gte=1,max=40"`
}

func (r CreateEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprog{
		Name:           r.Name,
		EducationLevel: r.EducationLevel,
		Speciality:     r.Speciality,
		KnowledgeField: r.KnowledgeField,
	}, nil
}

func (r UpdateEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprog{
		Name:           r.Name,
		EducationLevel: r.EducationLevel,
		Speciality:     r.Speciality,
		KnowledgeField: r.KnowledgeField,
	}, nil
}

//type UpdateEduprogRequest struct {
//
//}
