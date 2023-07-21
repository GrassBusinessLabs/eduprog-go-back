package requests

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
)

type CreateEduprogRequest struct {
	Name           string `json:"name" validate:"required,gte=1,max=50"`
	EducationLevel string `json:"education_level" validate:"required,gte=1,max=50"`
	SpecialityCode string `json:"speciality_code" validate:"required,gte=1,max=3"`
	ApprovalYear   int    `json:"approval_year" validate:"required,number"`
}

type UpdateEduprogRequest struct {
	Name           string `json:"name" validate:"required,gte=1,max=50"`
	EducationLevel string `json:"education_level" validate:"required,gte=1,max=50"`
	SpecialityCode string `json:"speciality_code" validate:"required,gte=1,max=3"`
	ApprovalYear   int    `json:"approval_year" validate:"number"`
}

type DuplicateEduprogRequest struct {
	Name         string `json:"name" validate:"required,gte=1,max=50"`
	ApprovalYear int    `json:"approval_year" validate:"number,required"`
}

func (r CreateEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprog{
		Name:           r.Name,
		EducationLevel: r.EducationLevel,
		SpecialtyCode:  r.SpecialityCode,
		ApprovalYear:   r.ApprovalYear,
	}, nil
}

func (r UpdateEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprog{
		Name:           r.Name,
		EducationLevel: r.EducationLevel,
		SpecialtyCode:  r.SpecialityCode,
		ApprovalYear:   r.ApprovalYear,
	}, nil
}

func (r DuplicateEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprog{
		Name:         r.Name,
		ApprovalYear: r.ApprovalYear,
	}, nil
}
