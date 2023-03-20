package requests

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type AddEduprogresultToEduprogRequest struct {
	EduprogId  uint64 `json:"eduprog_id" validate:"required,number"`
	Definition string `json:"definition" validate:"required,alphanum,gte=1,max=500"`
}

func (r AddEduprogresultToEduprogRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogresult{
		EduprogId:  r.EduprogId,
		Definition: r.Definition,
	}, nil
}
