package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
)

type CompetenciesBaseService interface {
	CreateCompetency(competency domain.CompetenciesBase) (domain.CompetenciesBase, error)
	UpdateCompetency(competency domain.CompetenciesBase, id uint64) (domain.CompetenciesBase, error)
	ShowAllCompetencies() ([]domain.CompetenciesBase, error)
	ShowCompetenciesByType(ttype string, specialty string) ([]domain.CompetenciesBase, error)
	ShowCompetenciesByEduprogData(ttype string, specialty string, edLevel string) ([]domain.CompetenciesBase, error)
	FindById(id uint64) (domain.CompetenciesBase, error)
	Delete(id uint64) error
}

type competenciesBaseService struct {
	competenciesBaseRepo eduprog.CompetenciesBaseRepository
}

func NewCompetenciesBaseService(cb eduprog.CompetenciesBaseRepository) CompetenciesBaseService {
	return competenciesBaseService{
		competenciesBaseRepo: cb,
	}
}

func (s competenciesBaseService) CreateCompetency(competency domain.CompetenciesBase) (domain.CompetenciesBase, error) {
	e, err := s.competenciesBaseRepo.CreateCompetency(competency)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return domain.CompetenciesBase{}, err
	}
	return e, nil
}

func (s competenciesBaseService) UpdateCompetency(competency domain.CompetenciesBase, id uint64) (domain.CompetenciesBase, error) {
	e, err := s.competenciesBaseRepo.UpdateCompetency(competency, id)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return domain.CompetenciesBase{}, err
	}
	return e, nil
}

func (s competenciesBaseService) ShowAllCompetencies() ([]domain.CompetenciesBase, error) {
	e, err := s.competenciesBaseRepo.ShowAllCompetencies()
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return []domain.CompetenciesBase{}, err
	}
	return e, nil
}

func (s competenciesBaseService) ShowCompetenciesByType(ttype string, specialty string) ([]domain.CompetenciesBase, error) {
	e, err := s.competenciesBaseRepo.ShowCompetenciesByType(ttype, specialty)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return []domain.CompetenciesBase{}, err
	}
	return e, nil
}

func (s competenciesBaseService) ShowCompetenciesByEduprogData(ttype string, specialty string, edLevel string) ([]domain.CompetenciesBase, error) {
	e, err := s.competenciesBaseRepo.ShowCompetenciesByEduprogData(ttype, specialty, edLevel)
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

func (s competenciesBaseService) Delete(id uint64) error {
	err := s.competenciesBaseRepo.Delete(id)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return err
	}
	return nil
}
