package requests

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type CreateEducompRelationRequest struct {
	EduprogId   uint64 `json:"eduprog_id" validate:"required,number"`
	BaseCompId  uint64 `json:"base_comp_id" validate:"number,required"`
	ChildCompId uint64 `json:"child_comp_id" validate:"number,required"`
}

func (r CreateEducompRelationRequest) ToDomainModel() (interface{}, error) {
	return domain.Educomp_relations{
		EduprogId:   r.EduprogId,
		BaseCompId:  r.BaseCompId,
		ChildCompId: r.ChildCompId,
	}, nil
}
