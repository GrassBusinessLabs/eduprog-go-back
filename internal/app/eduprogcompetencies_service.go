package app

import (
	"fmt"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
)

type EduprogcompetenciesService interface {
	AddCompetencyToEduprog(eduprogcompetency domain.Eduprogcompetencies) (domain.Eduprogcompetencies, error)
	AddCustomCompetecyToEduprog(eduprogcompetency domain.Eduprogcompetencies) (domain.Eduprogcompetencies, error)
	AddAllCompetencies(eduprogId uint64, ttype string) ([]domain.Eduprogcompetencies, error)
	UpdateCompetency(ref, req domain.Eduprogcompetencies) (domain.Eduprogcompetencies, error)
	ShowCompetenciesByEduprogId(eduprogId uint64) ([]domain.Eduprogcompetencies, error)
	ShowCompetenciesByType(eduprogId uint64, ttype string) ([]domain.Eduprogcompetencies, error)
	FindById(competencyId uint64) (domain.Eduprogcompetencies, error)
	Delete(competencyId uint64) error
	DeleteAllCompetencies(eduprogId uint64, ttype string) error
	SetCompetenciesBaseService(competenciesBaseService *CompetenciesBaseService)
}

type eduprogcompetenciesService struct {
	eduprogcompetenciesRepo eduprog.EduprogcompetenciesRepository
	competenciesBaseService CompetenciesBaseService
}

func NewEduprogcompetenciesService(ecr eduprog.EduprogcompetenciesRepository, cbs CompetenciesBaseService) EduprogcompetenciesService {
	return &eduprogcompetenciesService{
		eduprogcompetenciesRepo: ecr,
		competenciesBaseService: cbs,
	}
}

func (s *eduprogcompetenciesService) SetCompetenciesBaseService(competenciesBaseService *CompetenciesBaseService) {
	if competenciesBaseService != nil {
		s.competenciesBaseService = *competenciesBaseService
	}
}

func (s *eduprogcompetenciesService) AddAllCompetencies(eduprogId uint64, ttype string) ([]domain.Eduprogcompetencies, error) {
	eduprogcompetencies, err := s.ShowCompetenciesByType(eduprogId, ttype)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return []domain.Eduprogcompetencies{}, err
	}

	for i := range eduprogcompetencies {
		err = s.Delete(eduprogcompetencies[i].Id)
		if err != nil {
			log.Printf("EduprogcompetenciesService: %s", err)
			return []domain.Eduprogcompetencies{}, err
		}
	}

	baseCompetencies, err := s.competenciesBaseService.ShowCompetenciesByEduprogData(ttype, eduprogId)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return []domain.Eduprogcompetencies{}, err
	}

	var eduprogcompetenciesList []domain.Eduprogcompetencies

	for i := range baseCompetencies {
		var eduprogcompetency domain.Eduprogcompetencies

		eduprogcompetency.CompetencyId = baseCompetencies[i].Id
		eduprogcompetency.EduprogId = eduprogId
		eduprogcompetency.Type = baseCompetencies[i].Type
		eduprogcompetency.Code = baseCompetencies[i].Code
		eduprogcompetency.Definition = baseCompetencies[i].Definition

		eduprogcompetency, err = s.AddCompetencyToEduprog(eduprogcompetency)
		if err != nil {
			log.Printf("EduprogcompetenciesService: %s", err)
			return []domain.Eduprogcompetencies{}, err
		}
		//competencyId,err = c.eduprogcompetenciesService.FindById(eduprogcompetency.Id)

		eduprogcompetenciesList = append(eduprogcompetenciesList, eduprogcompetency)
	}

	return eduprogcompetenciesList, nil
}

func (s *eduprogcompetenciesService) AddCustomCompetecyToEduprog(eduprogcompetency domain.Eduprogcompetencies) (domain.Eduprogcompetencies, error) {
	eduprogcompetency.CompetencyId = 65

	allEdpcompetencies, err := s.ShowCompetenciesByType(eduprogcompetency.EduprogId, eduprogcompetency.Type)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return domain.Eduprogcompetencies{}, err
	}

	var maxCode uint64 = 0

	for i := range allEdpcompetencies {
		if allEdpcompetencies[i].EduprogId == eduprogcompetency.EduprogId {
			if i == 0 || allEdpcompetencies[i].Code > maxCode {
				maxCode = allEdpcompetencies[i].Code
			}
		}
	}

	eduprogcompetency.Code = maxCode + 1

	eduprogcompetency, err = s.AddCompetencyToEduprog(eduprogcompetency)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return domain.Eduprogcompetencies{}, err
	}

	return eduprogcompetency, nil
}

func (s *eduprogcompetenciesService) AddCompetencyToEduprog(eduprogcompetency domain.Eduprogcompetencies) (domain.Eduprogcompetencies, error) {
	baseCompetency, err := s.competenciesBaseService.FindById(eduprogcompetency.CompetencyId)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return domain.Eduprogcompetencies{}, err
	}

	if eduprogcompetency.Definition == "" {
		eduprogcompetency.Definition = baseCompetency.Definition
	}

	eduprogcompetency.Type = baseCompetency.Type

	allEdpcompetencies, err := s.ShowCompetenciesByType(eduprogcompetency.EduprogId, eduprogcompetency.Type)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return domain.Eduprogcompetencies{}, err
	}

	var maxCode uint64 = 0

	for i := range allEdpcompetencies {
		if allEdpcompetencies[i].CompetencyId == eduprogcompetency.CompetencyId {
			log.Printf("EduprogcompetenciesService: %s", err)
			return domain.Eduprogcompetencies{}, fmt.Errorf("competency is already in eduprog")
		}
		if i == 0 || allEdpcompetencies[i].Code > maxCode {
			maxCode = allEdpcompetencies[i].Code
		}
	}

	eduprogcompetency.Code = maxCode + 1

	e, err := s.eduprogcompetenciesRepo.AddCompetencyToEduprog(eduprogcompetency)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return domain.Eduprogcompetencies{}, err
	}
	return e, err
}

func (s *eduprogcompetenciesService) UpdateCompetency(ref, req domain.Eduprogcompetencies) (domain.Eduprogcompetencies, error) {
	if req.Definition != "" {
		ref.Definition = req.Definition
	}

	e, err := s.eduprogcompetenciesRepo.UpdateCompetency(ref)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return domain.Eduprogcompetencies{}, err
	}
	return e, err
}

func (s *eduprogcompetenciesService) ShowCompetenciesByEduprogId(eduprogId uint64) ([]domain.Eduprogcompetencies, error) {
	e, err := s.eduprogcompetenciesRepo.ShowCompetenciesByEduprogId(eduprogId)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return []domain.Eduprogcompetencies{}, err
	}
	return e, err
}

func (s *eduprogcompetenciesService) ShowCompetenciesByType(eduprogId uint64, ttype string) ([]domain.Eduprogcompetencies, error) {
	e, err := s.eduprogcompetenciesRepo.ShowCompetenciesByType(eduprogId, ttype)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return []domain.Eduprogcompetencies{}, err
	}
	return e, err
}

func (s *eduprogcompetenciesService) FindById(competencyId uint64) (domain.Eduprogcompetencies, error) {
	e, err := s.eduprogcompetenciesRepo.FindById(competencyId)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return domain.Eduprogcompetencies{}, err
	}
	return e, err
}

func (s *eduprogcompetenciesService) Delete(competencyId uint64) error {
	competency, err := s.FindById(competencyId)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return err
	}

	err = s.Delete(competencyId)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return err
	}

	allEdpcompetencies, err := s.ShowCompetenciesByType(competency.EduprogId, competency.Type)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return err
	}

	for i := range allEdpcompetencies {
		if allEdpcompetencies[i].Code > competency.Code {
			allEdpcompetencies[i].Code = allEdpcompetencies[i].Code - 1
			_, err = s.UpdateCompetency(allEdpcompetencies[i], allEdpcompetencies[i])
			if err != nil {
				log.Printf("EduprogcompetenciesService: %s", err)
				return err
			}
		}
	}

	return err
}

func (s *eduprogcompetenciesService) DeleteAllCompetencies(eduprogId uint64, ttype string) error {
	if ttype != "ZK" && ttype != "FK" && ttype != "PR" {
		log.Printf("EduprogcompetenciesService: %s", fmt.Errorf("only ZK, FK, PR"))
		return fmt.Errorf("only ZK, FK, PR")
	}

	eduprogcompetencies, err := s.ShowCompetenciesByType(eduprogId, ttype)
	if err != nil {
		log.Printf("EduprogcompetenciesService: %s", err)
		return err
	}

	for i := range eduprogcompetencies {
		err = s.Delete(eduprogcompetencies[i].Id)
		if err != nil {
			log.Printf("EduprogcompetenciesService: %s", err)
			return err
		}
	}

	return nil
}
