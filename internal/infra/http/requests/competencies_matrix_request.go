package requests

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type CreateCompetenciesMatrixRelationRequest struct {
	EduprogId    uint64 `json:"eduprog_id" validate:"required,number"`
	ComponentId  uint64 `json:"component_id" validate:"required,number"`
	CompetencyId uint64 `json:"competency_id" validate:"required,number"`
}

func (r CreateCompetenciesMatrixRelationRequest) ToDomainModel() (interface{}, error) {
	return domain.CompetenciesMatrix{
		EduprogId:    r.EduprogId,
		ComponentId:  r.ComponentId,
		CompetencyId: r.CompetencyId,
	}, nil
}
