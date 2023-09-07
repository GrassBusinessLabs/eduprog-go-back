package resources

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type EducompRelationsDto struct {
	EduprogId   uint64 `json:"eduprog_id"`
	BaseCompId  uint64 `json:"base_comp_id"`
	ChildCompId uint64 `json:"child_comp_id"`
}

func (d EducompRelationsDto) DomainToDto(relation domain.EducompRelations) EducompRelationsDto {
	return EducompRelationsDto{
		EduprogId:   relation.EduprogId,
		BaseCompId:  relation.BaseCompId,
		ChildCompId: relation.ChildCompId,
	}
}

func (d EducompRelationsDto) DomainToDtoCollection(relation []domain.EducompRelations) []EducompRelationsDto {
	result := make([]EducompRelationsDto, len(relation))

	for i := range relation {
		result[i] = d.DomainToDto(relation[i])
	}

	return result
}
