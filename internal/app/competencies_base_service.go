package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database"
	"log"
)

type CompetenciesBaseService interface {
	ShowAllCompetencies() ([]domain.CompetenciesBase, error)
	ShowZK() ([]domain.CompetenciesBase, error)
	ShowFK() ([]domain.CompetenciesBase, error)
	FindById(id uint64) (domain.CompetenciesBase, error)
}

type competenciesBaseService struct {
	competenciesBaseRepo database.CompetenciesBaseRepository
}

func NewCompetenciesBaseService(cb database.CompetenciesBaseRepository) CompetenciesBaseService {
	return competenciesBaseService{
		competenciesBaseRepo: cb,
	}
}

func (s competenciesBaseService) ShowAllCompetencies() ([]domain.CompetenciesBase, error) {
	e, err := s.competenciesBaseRepo.ShowAllCompetencies()
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return []domain.CompetenciesBase{}, err
	}
	return e, nil
}

func (s competenciesBaseService) ShowZK() ([]domain.CompetenciesBase, error) {
	e, err := s.competenciesBaseRepo.ShowZK()
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return []domain.CompetenciesBase{}, err
	}
	return e, nil
}

func (s competenciesBaseService) ShowFK() ([]domain.CompetenciesBase, error) {
	e, err := s.competenciesBaseRepo.ShowFK()
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return []domain.CompetenciesBase{}, err
	}
	return e, nil
}

func (s competenciesBaseService) FindById(id uint64) (domain.CompetenciesBase, error) {
	e, err := s.competenciesBaseRepo.FindById(id)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return domain.CompetenciesBase{}, err
	}
	return e, nil
}
