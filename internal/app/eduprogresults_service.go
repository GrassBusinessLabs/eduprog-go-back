package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
)

type EduprogresultsService interface {
	AddEduprogresultToEduprog(eduprogresult domain.Eduprogresult) (domain.Eduprogresult, error)
	UpdateEduprogresult(eduprogresult domain.Eduprogresult, id uint64) (domain.Eduprogresult, error)
	ShowEduprogResultsByEduprogId(eduprogId uint64) ([]domain.Eduprogresult, error)
	FindById(eduprogresultId uint64) (domain.Eduprogresult, error)
	Delete(eduprogresultId uint64) error
}

type eduprogresultsService struct {
	eduprogresultsRepo eduprog.EduprogresultsRepository
}

func NewEduprogresultsService(err eduprog.EduprogresultsRepository) EduprogresultsService {
	return eduprogresultsService{
		eduprogresultsRepo: err,
	}
}

func (s eduprogresultsService) AddEduprogresultToEduprog(eduprogresult domain.Eduprogresult) (domain.Eduprogresult, error) {
	e, err := s.eduprogresultsRepo.AddEduprogresultToEduprog(eduprogresult)
	if err != nil {
		log.Printf("Eduprogresult service: %s", err)
		return domain.Eduprogresult{}, err
	}
	return e, err
}

func (s eduprogresultsService) UpdateEduprogresult(eduprogresult domain.Eduprogresult, id uint64) (domain.Eduprogresult, error) {
	e, err := s.eduprogresultsRepo.UpdateEduprogresult(eduprogresult, id)
	if err != nil {
		log.Printf("Eduprogresult service: %s", err)
		return domain.Eduprogresult{}, err
	}
	return e, err
}

func (s eduprogresultsService) ShowEduprogResultsByEduprogId(eduprogId uint64) ([]domain.Eduprogresult, error) {
	e, err := s.eduprogresultsRepo.ShowEduprogResultsByEduprogId(eduprogId)
	if err != nil {
		log.Printf("Eduprogresult service: %s", err)
		return []domain.Eduprogresult{}, err
	}
	return e, err
}

func (s eduprogresultsService) FindById(eduprogresultId uint64) (domain.Eduprogresult, error) {
	e, err := s.eduprogresultsRepo.FindById(eduprogresultId)
	if err != nil {
		log.Printf("Eduprogresult service: %s", err)
		return domain.Eduprogresult{}, err
	}
	return e, err
}

func (s eduprogresultsService) Delete(eduprogresultId uint64) error {
	err := s.eduprogresultsRepo.Delete(eduprogresultId)
	if err != nil {
		log.Printf("Eduprogresult service: %s", err)
		return err
	}
	return err
}
