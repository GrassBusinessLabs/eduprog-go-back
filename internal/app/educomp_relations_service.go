package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
)

type EducompRelationsService interface {
	CreateRelation(relation domain.EducompRelations) (domain.EducompRelations, error)
	ShowByEduprogId(eduprogId uint64) ([]domain.EducompRelations, error)
	ShowPossibleRelationsForComp(eduprogId, eduprogcompId uint64) ([]domain.Eduprogcomp, error)
	Delete(baseCompId uint64, childCompId uint64) error
}

type educompRelationsService struct {
	educompRelationsRepo eduprog.EducompRelationsRepository
	eduprogschemeService EduprogschemeService
	eduprogcompService   EduprogcompService
}

func NewEducompRelationsService(err eduprog.EducompRelationsRepository, ess EduprogschemeService, ecs EduprogcompService) EducompRelationsService {
	return educompRelationsService{
		educompRelationsRepo: err,
		eduprogschemeService: ess,
		eduprogcompService:   ecs,
	}
}

func (s educompRelationsService) CreateRelation(relation domain.EducompRelations) (domain.EducompRelations, error) {
	e, err := s.educompRelationsRepo.CreateRelation(relation)
	if err != nil {
		log.Printf("EducompRelationsService: %s", err)
		return domain.EducompRelations{}, err
	}
	return e, err
}

func (s educompRelationsService) ShowByEduprogId(eduprogId uint64) ([]domain.EducompRelations, error) {
	e, err := s.educompRelationsRepo.ShowByEduprogId(eduprogId)
	if err != nil {
		log.Printf("EducompRelationsService: %s", err)
		return []domain.EducompRelations{}, err
	}
	return e, err
}

func (s educompRelationsService) ShowPossibleRelationsForComp(eduprogId, eduprogcompId uint64) ([]domain.Eduprogcomp, error) {
	eduprogscheme, err := s.eduprogschemeService.ShowSchemeByEduprogId(eduprogId)
	if err != nil {
		log.Printf("EducompRelationsService: %s", err)
		return []domain.Eduprogcomp{}, err
	}

	eduprogcomps, err := s.eduprogcompService.ShowListByEduprogId(eduprogcompId)
	if err != nil {
		log.Printf("EducompRelationsService: %s", err)
		return []domain.Eduprogcomp{}, err
	}

	var result []domain.Eduprogcomp
	var maxCompSemester uint64 = 0
	for i := range eduprogscheme {
		if eduprogscheme[i].EduprogcompId == eduprogcompId {
			if maxCompSemester < eduprogscheme[i].SemesterNum {
				maxCompSemester = eduprogscheme[i].SemesterNum
			}
		}
	}

	for i := range eduprogscheme {
		if eduprogscheme[i].EduprogcompId == eduprogcompId {

			for i2 := range eduprogscheme {
				if eduprogscheme[i2].SemesterNum > maxCompSemester {
					for i3 := range eduprogcomps {
						if eduprogcomps[i3].Id == eduprogscheme[i2].EduprogcompId {
							result = append(result, eduprogcomps[i3])
						}
					}
				}
			}
		}
	}

	uniqes := s.unique(result)

	return uniqes, nil
}

func (s educompRelationsService) Delete(baseCompId uint64, childCompId uint64) error {
	err := s.educompRelationsRepo.Delete(baseCompId, childCompId)
	if err != nil {
		log.Printf("EducompRelationsService: %s", err)
		return err
	}
	return nil
}

func (s educompRelationsService) unique(compSlice []domain.Eduprogcomp) []domain.Eduprogcomp {
	keys := make(map[domain.Eduprogcomp]bool)
	var list []domain.Eduprogcomp
	for _, entry := range compSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
