package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
)

type ResultsMatrixService interface {
	CreateRelation(relation domain.ResultsMatrix) (domain.ResultsMatrix, error)
	ShowByEduprogId(eduprog_id uint64) ([]domain.ResultsMatrix, error)
	Delete(componentId uint64, eduprogresultId uint64) error
}

type resultsMatrixService struct {
	resultsMatrixRepo eduprog.ResultsMatrixRepository
}

func NewResultsMatrixService(rm eduprog.ResultsMatrixRepository) ResultsMatrixService {
	return resultsMatrixService{
		resultsMatrixRepo: rm,
	}
}

func (s resultsMatrixService) CreateRelation(relation domain.ResultsMatrix) (domain.ResultsMatrix, error) {
	e, err := s.resultsMatrixRepo.CreateRelation(relation)
	if err != nil {
		log.Printf("ResultsMatrixService: %s", err)
		return domain.ResultsMatrix{}, err
	}
	return e, err
}

func (s resultsMatrixService) ShowByEduprogId(eduprog_id uint64) ([]domain.ResultsMatrix, error) {
	e, err := s.resultsMatrixRepo.ShowByEduprogId(eduprog_id)
	if err != nil {
		log.Printf("ResultsMatrixService: %s", err)
		return []domain.ResultsMatrix{}, err
	}
	return e, err
}

func (s resultsMatrixService) Delete(componentId uint64, eduprogresultId uint64) error {
	err := s.resultsMatrixRepo.Delete(componentId, eduprogresultId)
	if err != nil {
		log.Printf("ResultsMatrixService: %s", err)
		return err
	}
	return err
}
