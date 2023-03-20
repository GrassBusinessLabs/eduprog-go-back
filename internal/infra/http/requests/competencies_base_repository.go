package requests

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type CreateCompetencyBaseRequest struct {
	Type       string `json:"type" validate:"required,gte=1,max=10"`
	Definition string `json:"definition" validate:"required,gte=1,max=500"`
	Specialty  string `json:"specialty" validate:"required,gte=1,max=3"`
}

type UpdateCompetencyBaseRequest struct {
	Type       string `json:"type" validate:"gte=1,max=10"`
	Definition string `json:"definition" validate:"gte=1,max=500"`
	Specialty  string `json:"specialty" validate:"gte=1,max=3"`
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
