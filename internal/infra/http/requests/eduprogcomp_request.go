package requests

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
)

type CreateEduprogcompRequest struct {
	Name        string `json:"name" validate:"required,gte=1,max=40"`
	Credits     uint64 `json:"credits" validate:"required"`
	ControlType string `json:"control_type" validate:"required,gte=1,max=40"`
	Type        string `json:"type" validate:"required,gte=1,max=40"`
	SubType     string `json:"sub_type" validate:"required,gte=1,max=40"`
	Category    string `json:"category" validate:"required,gte=1,max=40"`
	EduprogId   uint64 `json:"eduprog_id" validate:"required"`
}

type UpdateEduprogcompRequest struct {
	Name        string `json:"name" validate:"required,gte=1,max=40"`
	Credits     uint64 `json:"credits" validate:"required"`
	ControlType string `json:"control_type" validate:"required,gte=1,max=40"`
	Type        string `json:"type" validate:"required,gte=1,max=40"`
	SubType     string `json:"sub_type" validate:"required,gte=1,max=40"`
	Category    string `json:"category" validate:"required,gte=1,max=40"`
	EduprogId   uint64 `json:"eduprog_id" validate:"required"`
}

func (r CreateEduprogcompRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogcomp{

		Name:        r.Name,
		Credits:     r.Credits,
		ControlType: r.ControlType,
		Type:        r.Type,
		SubType:     r.SubType,
		Category:    r.Category,
		EduprogId:   r.EduprogId,
	}, nil
}

func (r UpdateEduprogcompRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogcomp{

		Name:        r.Name,
		Credits:     r.Credits,
		ControlType: r.ControlType,
		Type:        r.Type,
		SubType:     r.SubType,
		Category:    r.Category,
		EduprogId:   r.EduprogId,
	}, nil
}
