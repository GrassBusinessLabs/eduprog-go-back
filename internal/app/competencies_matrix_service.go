package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
)

type CompetenciesMatrixService interface {
	CreateRelation(relation domain.CompetenciesMatrix) (domain.CompetenciesMatrix, error)
	ShowByEduprogId(eduprog_id uint64) ([]domain.CompetenciesMatrix, error)
	Delete(componentId uint64, competencyId uint64) error
}

type competenciesMatrixService struct {
	competenciesMatrixRepo eduprog.CompetenciesMatrixRepository
}

func NewCompetenciesMatrixService(cm eduprog.CompetenciesMatrixRepository) CompetenciesMatrixService {
	return competenciesMatrixService{
		competenciesMatrixRepo: cm,
	}
}

func (s competenciesMatrixService) CreateRelation(relation domain.CompetenciesMatrix) (domain.CompetenciesMatrix, error) {
	e, err := s.competenciesMatrixRepo.CreateRelation(relation)
	if err != nil {
		log.Printf("CompetenciesMatrixService: %s", err)
		return domain.CompetenciesMatrix{}, err
	}
	return e, err
}

func (s competenciesMatrixService) ShowByEduprogId(eduprog_id uint64) ([]domain.CompetenciesMatrix, error) {
	e, err := s.competenciesMatrixRepo.ShowByEduprogId(eduprog_id)
	if err != nil {
		log.Printf("CompetenciesMatrixService: %s", err)
		return []domain.CompetenciesMatrix{}, err
	}
	return e, err
}

func (s competenciesMatrixService) Delete(componentId uint64, competencyId uint64) error {
	err := s.competenciesMatrixRepo.Delete(componentId, competencyId)
	if err != nil {
		log.Printf("CompetenciesMatrixService: %s", err)
		return err
	}
	return err
}
