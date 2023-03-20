package requests

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type AddCompetencyToEduprogRequest struct {
	CompetencyId uint64 `json:"competency_id" validate:"required,number"`
	EduprogId    uint64 `json:"eduprog_id" validate:"required,number"`
	Definition   string `json:"definition" validate:"alphanum,gte=1,max=500"`
}

type UpdateCompetencyRequest struct {
	Definition string `json:"definition" validate:"alphanum,gte=1,max=500"`
}

type AddCustomCompetencyToEduprogRequest struct {
	EduprogId  uint64 `json:"eduprog_id" validate:"required,number"`
	Type       string `json:"type" validate:"required,alpha"`
	Definition string `json:"definition" validate:"alphanum,gte=1,max=500"`
}

func (r AddCompetencyToEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogcompetencies{
		CompetencyId: r.CompetencyId,
		EduprogId:    r.EduprogId,
		Definition:   r.Definition,
	}, nil
}

func (r UpdateCompetencyRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogcompetencies{
		Definition: r.Definition,
	}, nil
}

func (r AddCustomCompetencyToEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogcompetencies{
		EduprogId:  r.EduprogId,
		Type:       r.Type,
		Definition: r.Definition,
	}, nil
}
