package resources

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type ResultsMatrixDto struct {
	EduprogId       uint64 `json:"eduprog_id"`
	ComponentId     uint64 `json:"component_id"`
	EduprogresultId uint64 `json:"eduprogresult_id"`
}

func (d ResultsMatrixDto) DomainToDto(relation domain.ResultsMatrix) ResultsMatrixDto {
	return ResultsMatrixDto{
		EduprogId:       relation.EduprogId,
		ComponentId:     relation.ComponentId,
		EduprogresultId: relation.EduprogresultId,
	}
}

func (d ResultsMatrixDto) DomainToDtoCollection(relation []domain.ResultsMatrix) []ResultsMatrixDto {
	result := make([]ResultsMatrixDto, len(relation))

	for i := range relation {
		result[i] = d.DomainToDto(relation[i])
	}

	return result
}
