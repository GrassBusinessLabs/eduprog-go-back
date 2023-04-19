package requests

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
)

type CreateEduprogRequest struct {
	Name           string `json:"name" validate:"required,gte=1,max=50"`
	EducationLevel string `json:"education_level" validate:"required,gte=1,max=50"`
	SpecialityCode string `json:"speciality_code" validate:"required,gte=1,max=3"`
	ApprovalYear   int    `json:"approval_year" validate:"required,number"`
	//ChildOf        uint64 `json:"child_of"`
}

type UpdateEduprogRequest struct {
	Name           string `json:"name" validate:"required,gte=1,max=50"`
	EducationLevel string `json:"education_level" validate:"required,gte=1,max=50"`
	SpecialityCode string `json:"speciality_code" validate:"required,gte=1,max=3"`
	ApprovalYear   int    `json:"approval_year" validate:"number"`
	//ChildOf        uint64 `json:"child_of"`
}

func (r CreateEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprog{
		Name:           r.Name,
		EducationLevel: r.EducationLevel,
		SpecialtyCode:  r.SpecialityCode,
		ApprovalYear:   r.ApprovalYear,
		//ChildOf:        r.ChildOf,
	}, nil
}

func (r UpdateEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprog{
		Name:           r.Name,
		EducationLevel: r.EducationLevel,
		SpecialtyCode:  r.SpecialityCode,
		ApprovalYear:   r.ApprovalYear,
		//ChildOf:        r.ChildOf,
	}, nil
}

//type UpdateEduprogRequest struct {
//
//}
