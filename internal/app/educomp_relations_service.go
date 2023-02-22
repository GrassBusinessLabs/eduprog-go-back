package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database"
	"log"
)

type EducompRelationsService interface {
	CreateRelation(relation domain.Educomp_relations) (domain.Educomp_relations, error)
	ShowByEduprogId(eduprog_id uint64) ([]domain.Educomp_relations, error)
	DeleteByBaseCompId(base_comp_id uint64) error
}

type educompRelationsService struct {
	educompRelationsRepo database.EducompRelationsRepository
}

func NewEducompRelationsService(err database.EducompRelationsRepository) EducompRelationsService {
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

func (s educompRelationsService) DeleteByBaseCompId(base_comp_id uint64) error {
	err := s.educompRelationsRepo.DeleteByBaseCompId(base_comp_id)
	if err != nil {
		log.Printf("EducompRelationsService: %s", err)
		return err
	}
	return nil
}
