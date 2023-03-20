package requests

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type CreateSpecialtyRequest struct {
	Code           string `json:"code" validate:"required,numeric"`
	Name           string `json:"name" validate:"required,gte=1,max=100"`
	KFCode         string `json:"kf_code" validate:"required,numeric"`
	KnowledgeField string `json:"knowledge_field" validate:"required,gte=1,max=100"`
}

type UpdateSpecialtyRequest struct {
	Code           string `json:"code" validate:"numeric"`
	Name           string `json:"name" validate:"gte=1,max=100"`
	KFCode         string `json:"kf_code" validate:"numeric"`
	KnowledgeField string `json:"knowledge_field" validate:"gte=1,max=100"`
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
