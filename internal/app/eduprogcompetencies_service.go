package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
)

type EduprogcompetenciesService interface {
	AddCompetencyToEduprog(eduprogcompetency domain.Eduprogcompetencies) (domain.Eduprogcompetencies, error)
	UpdateCompetency(eduprogcompetency domain.Eduprogcompetencies, id uint64) (domain.Eduprogcompetencies, error)
	ShowCompetenciesByEduprogId(eduprogId uint64) ([]domain.Eduprogcompetencies, error)
	FindById(competencyId uint64) (domain.Eduprogcompetencies, error)
	Delete(competencyId uint64) error
}

type eduprogcompetenciesService struct {
	eduprogcompetenciesRepo eduprog.EduprogcompetenciesRepository
}

func NewEduprogcompetenciesService(cb eduprog.EduprogcompetenciesRepository) EduprogcompetenciesService {
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

func (s eduprogcompetenciesService) UpdateCompetency(eduprogcompetency domain.Eduprogcompetencies, id uint64) (domain.Eduprogcompetencies, error) {
	e, err := s.eduprogcompetenciesRepo.UpdateCompetency(eduprogcompetency, id)
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
