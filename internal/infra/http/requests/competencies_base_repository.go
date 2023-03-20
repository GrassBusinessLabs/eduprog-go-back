package requests

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type CreateCompetencyBaseRequest struct {
	Type       string `json:"type" validate:"required"`
	Definition string `json:"definition" validate:"required"`
	Specialty  string `json:"specialty" validate:"required"`
}

type UpdateCompetencyBaseRequest struct {
	Type       string `json:"type" validate:"required"`
	Definition string `json:"definition" validate:"required"`
	Specialty  string `json:"specialty" validate:"required"`
}

func (r CreateCompetencyBaseRequest) ToDomainModel() (interface{}, error) {
	return domain.CompetenciesBase{
		Type:       r.Type,
		Definition: r.Definition,
		Specialty:  r.Specialty,
	}, nil
}

func (r UpdateCompetencyBaseRequest) ToDomainModel() (interface{}, error) {
	return domain.CompetenciesBase{
		Type:       r.Type,
		Definition: r.Definition,
		Specialty:  r.Specialty,
	}, nil
}
