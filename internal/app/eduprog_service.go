package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database"
	"log"
)

type EduprogService interface {
	Save(eduprog domain.Eduprog) (domain.Eduprog, error)
	Update(eduprog domain.Eduprog, id uint64) (domain.Eduprog, error)
	Delete(id uint64) error
}

type eduprogService struct {
	eduprogRepo database.EduprogRepository
}

func NewEduprogService(er database.EduprogRepository) EduprogService {
	return eduprogService{
		eduprogRepo: er,
	}
}

func (s eduprogService) Save(eduprog domain.Eduprog) (domain.Eduprog, error) {
	e, err := s.eduprogRepo.Save(eduprog)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}
	return e, err
}

func (s eduprogService) Update(eduprog domain.Eduprog, id uint64) (domain.Eduprog, error) {
	e, err := s.eduprogRepo.Update(eduprog, id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}
	return e, err
}

func (s eduprogService) Delete(id uint64) error {
	err := s.eduprogRepo.Delete(id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return err
	}
	return nil
}
