package requests

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type CreateDisciplineRequest struct {
	Name      string `json:"name" validate:"required,gte=1,max=50"`
	EduprogId uint64 `json:"eduprog_id" validate:"required,number"`
}

type UpdateDisciplineRequest struct {
	Name      string `json:"name"`
	EduprogId uint64 `json:"eduprog_id"`
}

func (r CreateDisciplineRequest) ToDomainModel() (interface{}, error) {
	return domain.Discipline{
		Name:      r.Name,
		EduprogId: r.EduprogId,
	}, nil
}

func (r UpdateDisciplineRequest) ToDomainModel() (interface{}, error) {
	return domain.Discipline{
		Name:      r.Name,
		EduprogId: r.EduprogId,
	}, nil
}
