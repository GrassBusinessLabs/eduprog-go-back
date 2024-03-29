package resources

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type EduprogcompetenciesDto struct {
	Id           uint64 `json:"id"`
	CompetencyId uint64 `json:"competency_id"`
	EduprogId    uint64 `json:"eduprog_id"`
	Code         uint64 `json:"code"`
	Type         string `json:"type"`
	Definition   string `json:"definition"`
}

func (d EduprogcompetenciesDto) DomainToDto(competency domain.Eduprogcompetencies) EduprogcompetenciesDto {
	return EduprogcompetenciesDto{
		Id:           competency.Id,
		CompetencyId: competency.CompetencyId,
		EduprogId:    competency.EduprogId,
		Code:         competency.Code,
		Type:         competency.Type,
		Definition:   competency.Definition,
	}
}

func (d EduprogcompetenciesDto) DomainToDtoCollection(competency []domain.Eduprogcompetencies) []EduprogcompetenciesDto {
	result := make([]EduprogcompetenciesDto, len(competency))

	for i := range competency {
		result[i] = d.DomainToDto(competency[i])
	}

	return result
}
