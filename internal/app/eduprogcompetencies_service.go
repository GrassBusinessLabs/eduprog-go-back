package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database"
	"log"
)

type EduprogcompetenciesService interface {
	AddCompetencyToEduprog(eduprogcompetency domain.Eduprogcompetencies) (domain.Eduprogcompetencies, error)
	ShowCompetenciesByEduprogId(eduprogId uint64) ([]domain.Eduprogcompetencies, error)
	FindById(competencyId uint64) (domain.Eduprogcompetencies, error)
	Delete(competencyId uint64) error
}

type eduprogcompetenciesService struct {
	eduprogcompetenciesRepo database.EduprogcompetenciesRepository
}

func NewEduprogcompetenciesService(cb database.EduprogcompetenciesRepository) EduprogcompetenciesService {
	return eduprogcompetenciesService{
		eduprogcompetenciesRepo: cb,
	}
}

func (s eduprogcompetenciesService) AddCompetencyToEduprog(eduprogcompetency domain.Eduprogcompetencies) (domain.Eduprogcompetencies, error) {
	e, err := s.eduprogcompetenciesRepo.AddCompetencyToEduprog(eduprogcompetency)
	if err != nil {
		log.Printf("Eduprogcompetency service: %s", err)
		return domain.Eduprogcompetencies{}, err
	}
	return e, err
}

func (s eduprogcompetenciesService) ShowCompetenciesByEduprogId(eduprogId uint64) ([]domain.Eduprogcompetencies, error) {
	e, err := s.eduprogcompetenciesRepo.ShowCompetenciesByEduprogId(eduprogId)
	if err != nil {
		log.Printf("Eduprogcompetency service: %s", err)
		return []domain.Eduprogcompetencies{}, err
	}
	return e, err
}

func (s eduprogcompetenciesService) FindById(competencyId uint64) (domain.Eduprogcompetencies, error) {
	e, err := s.eduprogcompetenciesRepo.FindById(competencyId)
	if err != nil {
		log.Printf("Eduprogcompetency service: %s", err)
		return domain.Eduprogcompetencies{}, err
	}
	return e, err
}

func (s eduprogcompetenciesService) Delete(competencyId uint64) error {
	err := s.eduprogcompetenciesRepo.Delete(competencyId)
	if err != nil {
		log.Printf("Eduprogcompetency service: %s", err)
		return err
	}
	return err
}
