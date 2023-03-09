package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
)

type EducompRelationsService interface {
	CreateRelation(relation domain.Educomp_relations) (domain.Educomp_relations, error)
	ShowByEduprogId(eduprog_id uint64) ([]domain.Educomp_relations, error)
	Delete(base_comp_id uint64, child_comp_id uint64) error
}

type educompRelationsService struct {
	educompRelationsRepo eduprog.EducompRelationsRepository
}

func NewEducompRelationsService(err eduprog.EducompRelationsRepository) EducompRelationsService {
	return educompRelationsService{
		educompRelationsRepo: err,
	}
}

func (s educompRelationsService) CreateRelation(relation domain.Educomp_relations) (domain.Educomp_relations, error) {
	e, err := s.educompRelationsRepo.CreateRelation(relation)
	if err != nil {
		log.Printf("EducompRelationsService: %s", err)
		return domain.Educomp_relations{}, err
	}
	return e, err
}

func (s educompRelationsService) ShowByEduprogId(eduprog_id uint64) ([]domain.Educomp_relations, error) {
	e, err := s.educompRelationsRepo.ShowByEduprogId(eduprog_id)
	if err != nil {
		log.Printf("EducompRelationsService: %s", err)
		return []domain.Educomp_relations{}, err
	}
	return e, err
}

func (s educompRelationsService) Delete(base_comp_id uint64, child_comp_id uint64) error {
	err := s.educompRelationsRepo.Delete(base_comp_id, child_comp_id)
	if err != nil {
		log.Printf("EducompRelationsService: %s", err)
		return err
	}
	return nil
}
