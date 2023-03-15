package requests

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type AddCompetencyToEduprogRequest struct {
	CompetencyId uint64 `json:"competency_id" validate:"required"`
	EduprogId    uint64 `json:"eduprog_id" validate:"required"`
	Redefinition string `json:"redefinition" `
}

type AddCustomCompetencyToEduprogRequest struct {
	EduprogId    uint64 `json:"eduprog_id" validate:"required"`
	Type         string `json:"type" validate:"required"`
	Redefinition string `json:"redefinition"`
}

func (r AddCompetencyToEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogcompetencies{
		CompetencyId: r.CompetencyId,
		EduprogId:    r.EduprogId,
		Redefinition: r.Redefinition,
	}, nil
}

func (r AddCustomCompetencyToEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogcompetencies{
		EduprogId:    r.EduprogId,
		Type:         r.Type,
		Redefinition: r.Redefinition,
	}, nil
}
