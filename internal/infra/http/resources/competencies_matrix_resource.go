package resources

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type CompetenciesMatrixDto struct {
	EduprogId    uint64 `json:"eduprog_id"`
	ComponentId  uint64 `json:"component_id"`
	CompetencyId uint64 `json:"competency_id"`
}

func (d CompetenciesMatrixDto) DomainToDto(relation domain.CompetenciesMatrix) CompetenciesMatrixDto {
	return CompetenciesMatrixDto{
		EduprogId:    relation.EduprogId,
		ComponentId:  relation.ComponentId,
		CompetencyId: relation.CompetencyId,
	}
}

func (d CompetenciesMatrixDto) DomainToDtoCollection(relation []domain.CompetenciesMatrix) []CompetenciesMatrixDto {
	result := make([]CompetenciesMatrixDto, len(relation))

	for i := range relation {
		result[i] = d.DomainToDto(relation[i])
	}

	return result
}
