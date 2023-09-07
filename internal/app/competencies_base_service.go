package app

import (
	"fmt"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
	"strconv"
)

type CompetenciesBaseService interface {
	CreateCompetency(competency domain.CompetenciesBase) (domain.CompetenciesBase, error)
	UpdateCompetency(ref, req domain.CompetenciesBase) (domain.CompetenciesBase, error)
	ShowAllCompetencies() ([]domain.CompetenciesBase, error)
	ShowCompetenciesByType(ttype string, specialty string) ([]domain.CompetenciesBase, error)
	ShowCompetenciesByEduprogData(ttype string, eduprogId uint64) ([]domain.CompetenciesBase, error)
	FindById(id uint64) (domain.CompetenciesBase, error)
	Delete(id uint64) error
}

type competenciesBaseService struct {
	competenciesBaseRepo eduprog.CompetenciesBaseRepository
	eduprogService       EduprogService
}

func NewCompetenciesBaseService(cb eduprog.CompetenciesBaseRepository, es EduprogService) CompetenciesBaseService {
	return competenciesBaseService{
		competenciesBaseRepo: cb,
		eduprogService:       es,
	}
}

func (s competenciesBaseService) CreateCompetency(competencyBase domain.CompetenciesBase) (domain.CompetenciesBase, error) {
	if competencyBase.Type != "ЗК" && competencyBase.Type != "ФК" && competencyBase.Type != "ПР" {
		log.Printf("CompetenciesBaseService: %s", fmt.Errorf("only 'ЗК', 'ФК' and 'ПР'"))
		return domain.CompetenciesBase{}, fmt.Errorf("only 'ЗК', 'ФК' and 'ПР'")
	}

	specialty, err := strconv.ParseUint(competencyBase.Specialty, 10, 64)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return domain.CompetenciesBase{}, err
	}
	if specialty < 11 || specialty > 293 {
		log.Printf("CompetenciesBaseService: %s", err)
		return domain.CompetenciesBase{}, err
	}

	allCompetencies, err := s.ShowCompetenciesByType(competencyBase.Type, competencyBase.Specialty)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return domain.CompetenciesBase{}, err
	}

	var maxCode uint64 = 0

	for i := range allCompetencies {
		if i == 0 || allCompetencies[i].Code > maxCode {
			maxCode = allCompetencies[i].Code
		}
	}

	competencyBase.Code = maxCode + 1

	competencyBase, err = s.competenciesBaseRepo.CreateCompetency(competencyBase)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return domain.CompetenciesBase{}, err
	}

	return competencyBase, nil
}

func (s competenciesBaseService) UpdateCompetency(ref, req domain.CompetenciesBase) (domain.CompetenciesBase, error) {
	if req.Type != "" {
		if req.Type != "ЗК" && req.Type != "ФК" && req.Type != "ПР" {
			log.Printf("CompetenciesBaseService: %s", fmt.Errorf("only 'ЗК', 'ФК' and 'ПР'"))
			return domain.CompetenciesBase{}, fmt.Errorf("only 'ЗК', 'ФК' and 'ПР'")
		}
		ref.Type = req.Type
	}
	if req.Definition != "" {
		ref.Definition = req.Definition
	}
	if req.Specialty != "" {
		specialty, err := strconv.ParseUint(req.Specialty, 10, 64)
		if err != nil {
			log.Printf("CompetenciesBaseService: %s", err)
			return domain.CompetenciesBase{}, err
		}
		if specialty < 11 || specialty > 293 {
			log.Printf("CompetenciesBaseService: %s", fmt.Errorf("specialties can be from 11 to 293"))
			return domain.CompetenciesBase{}, fmt.Errorf("specialties can be from 11 to 293")
		}
	}
	if req.EducationLevel != "" {
		_, err := s.eduprogService.GetOPPLevelData(req.EducationLevel)
		if err != nil {
			log.Printf("CompetenciesBaseService: %s", fmt.Errorf("education level error: no such level in enumeration, can use only `Початковий рівень (короткий цикл)`, `Перший (бакалаврський) рівень`, `Другий (магістерський) рівень`, `Третій (освітньо-науковий/освітньо-творчий) рівень`; get this value from method `LevelsList`"))
			return domain.CompetenciesBase{}, fmt.Errorf("education level error: no such level in enumeration, can use only `Початковий рівень (короткий цикл)`, `Перший (бакалаврський) рівень`, `Другий (магістерський) рівень`, `Третій (освітньо-науковий/освітньо-творчий) рівень`; get this value from method `LevelsList`")
		}
	}

	e, err := s.competenciesBaseRepo.UpdateCompetency(ref)
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
	specialtyInt, err := strconv.ParseInt(specialty, 10, 64)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return []domain.CompetenciesBase{}, err
	}
	if specialtyInt < 11 || specialtyInt > 293 {
		log.Printf("CompetenciesBaseService: %s", fmt.Errorf("specialty can be from 11 to 293"))
		return []domain.CompetenciesBase{}, fmt.Errorf("specialty can be from 11 to 293")
	}

	e, err := s.competenciesBaseRepo.ShowCompetenciesByType(ttype, specialty)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return []domain.CompetenciesBase{}, err
	}
	return e, nil
}

func (s competenciesBaseService) ShowCompetenciesByEduprogData(ttype string, eduprogId uint64) ([]domain.CompetenciesBase, error) {
	e, _, err := s.eduprogService.FindById(eduprogId)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return []domain.CompetenciesBase{}, err
	}

	competencies, err := s.competenciesBaseRepo.ShowCompetenciesByEduprogData(ttype, e.SpecialtyCode, e.EducationLevel)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return []domain.CompetenciesBase{}, err
	}

	return competencies, nil
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
	competency, err := s.FindById(id)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return err
	}

	err = s.competenciesBaseRepo.Delete(id)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return err
	}

	allCompetencies, err := s.ShowCompetenciesByType(competency.Type, competency.Specialty)
	if err != nil {
		log.Printf("CompetenciesBaseService: %s", err)
		return err
	}

	for i := range allCompetencies {
		if allCompetencies[i].Code > competency.Code {
			allCompetencies[i].Code = allCompetencies[i].Code - 1
			_, err = s.UpdateCompetency(allCompetencies[i], allCompetencies[i])
			if err != nil {
				log.Printf("CompetenciesBaseService: %s", err)
				return err
			}
		}
	}

	return nil
}
