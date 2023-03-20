package requests

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type CreateSpecialtyRequest struct {
	Code           string `json:"code" validate:"required"`
	Name           string `json:"name" validate:"required"`
	KFCode         string `json:"kf_code" validate:"required"`
	KnowledgeField string `json:"knowledge_field" validate:"required"`
}

type UpdateSpecialtyRequest struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	KFCode         string `json:"kf_code"`
	KnowledgeField string `json:"knowledge_field"`
}

func (r CreateSpecialtyRequest) ToDomainModel() (interface{}, error) {
	return domain.Specialty{
		Code:           r.Code,
		Name:           r.Name,
		KFCode:         r.KFCode,
		KnowledgeField: r.KnowledgeField,
	}, nil
}

func (r UpdateSpecialtyRequest) ToDomainModel() (interface{}, error) {
	return domain.Specialty{
		Code:           r.Code,
		Name:           r.Name,
		KFCode:         r.KFCode,
		KnowledgeField: r.KnowledgeField,
	}, nil
}
