package requests

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type CreateCompetencyBaseRequest struct {
	Type           string `json:"type" validate:"required,gte=1,max=10"`
	Definition     string `json:"definition" validate:"required,gte=1,max=500"`
	Specialty      string `json:"specialty" validate:"required,gte=1,max=3"`
	EducationLevel string `json:"education_level"`
}

type UpdateCompetencyBaseRequest struct {
	Type           string `json:"type" validate:"gte=1,max=10"`
	Definition     string `json:"definition" validate:"gte=1,max=500"`
	Specialty      string `json:"specialty" validate:"gte=1,max=3"`
	EducationLevel string `json:"education_level"`
}

func (r CreateCompetencyBaseRequest) ToDomainModel() (interface{}, error) {
	return domain.CompetenciesBase{
		Type:           r.Type,
		Definition:     r.Definition,
		Specialty:      r.Specialty,
		EducationLevel: r.EducationLevel,
	}, nil
}

func (r UpdateCompetencyBaseRequest) ToDomainModel() (interface{}, error) {
	return domain.CompetenciesBase{
		Type:           r.Type,
		Definition:     r.Definition,
		Specialty:      r.Specialty,
		EducationLevel: r.EducationLevel,
	}, nil
}
