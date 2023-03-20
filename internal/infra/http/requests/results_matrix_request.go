package requests

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type CreateResultsMatrixRelationRequest struct {
	EduprogId       uint64 `json:"eduprog_id" validate:"required,number"`
	ComponentId     uint64 `json:"component_id" validate:"number,required"`
	EduprogresultId uint64 `json:"eduprogresult_id" validate:"number,required"`
}

func (r CreateResultsMatrixRelationRequest) ToDomainModel() (interface{}, error) {
	return domain.ResultsMatrix{
		EduprogId:       r.EduprogId,
		ComponentId:     r.ComponentId,
		EduprogresultId: r.EduprogresultId,
	}, nil
}
