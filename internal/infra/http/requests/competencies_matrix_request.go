package requests

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type CreateCompetenciesMatrixRelationRequest struct {
	EduprogId    uint64 `json:"eduprog_id" validate:"required"`
	ComponentId  uint64 `json:"component_id" validate:"numeric,required"`
	CompetencyId uint64 `json:"competency_id" validate:"numeric,required"`
}

func (r CreateCompetenciesMatrixRelationRequest) ToDomainModel() (interface{}, error) {
	return domain.CompetenciesMatrix{
		EduprogId:    r.EduprogId,
		ComponentId:  r.ComponentId,
		CompetencyId: r.CompetencyId,
	}, nil
}
