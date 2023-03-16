package resources

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type CompetenciesBaseDto struct {
	Id         uint64 `json:"id"`
	Type       string `json:"type"`
	Code       uint64 `json:"code"`
	Definition string `json:"definition"`
	Specialty  string `json:"specialty"`
}

func (d CompetenciesBaseDto) DomainToDto(competency domain.CompetenciesBase) CompetenciesBaseDto {
	return CompetenciesBaseDto{
		Id:         competency.Id,
		Type:       competency.Type,
		Code:       competency.Code,
		Definition: competency.Definition,
		Specialty:  competency.Specialty,
	}
}

func (d CompetenciesBaseDto) DomainToDtoCollection(competency []domain.CompetenciesBase) []CompetenciesBaseDto {
	result := make([]CompetenciesBaseDto, len(competency))

	for i := range competency {
		result[i] = d.DomainToDto(competency[i])
	}

	return result
}
