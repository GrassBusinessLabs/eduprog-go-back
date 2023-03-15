package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
)

type CompetenciesBaseService interface {
	ShowAllCompetencies() ([]domain.CompetenciesBase, error)
	ShowCompetenciesByType(ttype string) ([]domain.CompetenciesBase, error)
	FindById(id uint64) (domain.CompetenciesBase, error)
}

type competenciesBaseService struct {
	competenciesBaseRepo eduprog.CompetenciesBaseRepository
}

func NewCompetenciesBaseService(cb eduprog.CompetenciesBaseRepository) CompetenciesBaseService {
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

func (s competenciesBaseService) ShowCompetenciesByType(ttype string) ([]domain.CompetenciesBase, error) {
	e, err := s.competenciesBaseRepo.ShowCompetenciesByType(ttype)
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
