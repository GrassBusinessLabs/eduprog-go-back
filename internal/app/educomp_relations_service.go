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
	ShowPossibleRelations(eduprogId uint64) ([]domain.EducompWithPossibleRelations, error)
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

func (s educompRelationsService) ShowPossibleRelations(eduprogId uint64) ([]domain.EducompWithPossibleRelations, error) {
	eduprogscheme, err := s.eduprogschemeService.ShowSchemeByEduprogId(eduprogId)
	if err != nil {
		log.Printf("EducompRelationsService: %s", err)
		return []domain.EducompWithPossibleRelations{}, err
	}

	eduprogcomps, err := s.eduprogcompService.ShowListByEduprogId(eduprogId)
	if err != nil {
		log.Printf("EducompRelationsService: %s", err)
		return []domain.EducompWithPossibleRelations{}, err
	}

	educompsWithPossibleRelations := make([]domain.EducompWithPossibleRelations, len(eduprogcomps))

	for i, eduprogcomp := range eduprogcomps {
		educompsWithPossibleRelations[i].Id = eduprogcomp.Id
		educompsWithPossibleRelations[i].Name = eduprogcomp.Name
		educompsWithPossibleRelations[i].Code = eduprogcomp.Code
		educompsWithPossibleRelations[i].Type = eduprogcomp.Type
		educompsWithPossibleRelations[i].ControlType = eduprogcomp.ControlType
		educompsWithPossibleRelations[i].Credits = eduprogcomp.Credits
		educompsWithPossibleRelations[i].BlockName = eduprogcomp.BlockName
		educompsWithPossibleRelations[i].BlockNum = eduprogcomp.BlockNum
		educompsWithPossibleRelations[i].EduprogId = eduprogcomp.EduprogId
		for _, schemecomp := range eduprogscheme {
			if eduprogcomp.Id == schemecomp.EduprogcompId {
				for _, schemecomp2 := range eduprogscheme {
					if schemecomp.SemesterNum < schemecomp2.SemesterNum {
						eduprogcompById, err := s.eduprogcompService.FindById(schemecomp2.EduprogcompId)
						if err != nil {
							log.Printf("EducompRelationsService: %s", err)
							return []domain.EducompWithPossibleRelations{}, err
						}

						educompsWithPossibleRelations[i].PossibleRelations = append(educompsWithPossibleRelations[i].PossibleRelations, eduprogcompById)
					}
				}
			}
		}
	}

	return educompsWithPossibleRelations, nil
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

	possibleRelations := make([]domain.Eduprogcomp, len(eduprogcomps))

	for _, schemecomp := range eduprogscheme {
		if schemecomp.EduprogcompId == eduprogcompId {
			for _, schemecomp2 := range eduprogscheme {
				if schemecomp.SemesterNum < schemecomp2.SemesterNum {
					eduprogcompById, err := s.eduprogcompService.FindById(schemecomp2.EduprogcompId)
					if err != nil {
						log.Printf("EducompRelationsService: %s", err)
						return []domain.Eduprogcomp{}, err
					}

					possibleRelations = append(possibleRelations, eduprogcompById)
				}
			}
		}
	}

	return possibleRelations, nil
}

func (s educompRelationsService) Delete(baseCompId uint64, childCompId uint64) error {
	err := s.educompRelationsRepo.Delete(baseCompId, childCompId)
	if err != nil {
		log.Printf("EducompRelationsService: %s", err)
		return err
	}
	return nil
}
