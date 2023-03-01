package requests

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type AddCompetencyToEduprogRequest struct {
	CompetencyId uint64 `json:"competency_id" validate:"required"`
	EduprogId    uint64 `json:"eduprog_id" validate:"required"`
	Code         uint64 `json:"code" validate:"required"`
	Redefiniton  string `json:"redefiniton" `
}

func (r AddCompetencyToEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogcompetencies{
		CompetencyId: r.CompetencyId,
		EduprogId:    r.EduprogId,
		Code:         r.Code,
		Redefinition: r.Redefiniton,
	}, nil
}
