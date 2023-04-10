package resources

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
)

type DisciplineDto struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	Rows      uint64 `json:"rows"`
	EduprogId uint64 `json:"eduprog_id"`
}

func (d DisciplineDto) DomainToDto(discipline domain.Discipline) DisciplineDto {
	return DisciplineDto{
		Id:        discipline.Id,
		Name:      discipline.Name,
		Rows:      discipline.Rows,
		EduprogId: discipline.EduprogId,
	}
}

func (d DisciplineDto) DomainToDtoCollection(discipline []domain.Discipline) []DisciplineDto {
	result := make([]DisciplineDto, len(discipline))

	for i := range discipline {
		result[i] = d.DomainToDto(discipline[i])
	}

	return result
}
