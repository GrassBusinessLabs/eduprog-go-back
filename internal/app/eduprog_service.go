package app

import (
	"errors"
	"fmt"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
	"time"
)

type EduprogService interface {
	Save(eduprog domain.Eduprog, userId uint64) (domain.Eduprog, error)
	Update(eduprog domain.Eduprog, id uint64) (domain.Eduprog, error)
	ShowList() (domain.Eduprogs, error)
	FindById(id uint64) (domain.Eduprog, error)
	GetOPPLevelsList() ([]domain.OPPLevelStruct, error)
	GetOPPLevelData(level string) (domain.OPPLevelStruct, error)
	Delete(id uint64) error
}

type eduprogService struct {
	eduprogRepo        eduprog.EduprogRepository
	specialtiesService SpecialtiesService
}

func NewEduprogService(er eduprog.EduprogRepository, ss SpecialtiesService) EduprogService {
	return eduprogService{
		eduprogRepo:        er,
		specialtiesService: ss,
	}
}

func (s eduprogService) Save(eduprog domain.Eduprog, userId uint64) (domain.Eduprog, error) {
	var err error

	maxYear := time.Now().Year() + 10
	if eduprog.ApprovalYear <= 1990 || eduprog.ApprovalYear > maxYear {
		log.Printf("EduprogService: %s", fmt.Errorf("approval year cant be less then 1990 and greater than %d", maxYear))
		return domain.Eduprog{}, fmt.Errorf("approval year cant be less then 1990 and greater than %d", maxYear)
	}

	eduprog.UserId = userId

	levelData, err := s.GetOPPLevelData(eduprog.EducationLevel)
	if err != nil {
		log.Printf("EduprogService: %s", errors.New("eduprog level error: no such level in enumeration, can use only `Початковий рівень (короткий цикл)`, `Перший (бакалаврський) рівень`, `Другий (магістерський) рівень`, `Третій (освітньо-науковий/освітньо-творчий) рівень`; get this value from method `LevelsList`"))
		return domain.Eduprog{}, errors.New("eduprog level error: no such level in enumeration, can use only `Початковий рівень (короткий цикл)`, `Перший (бакалаврський) рівень`, `Другий (магістерський) рівень`, `Третій (освітньо-науковий/освітньо-творчий) рівень`; get this value from method `LevelsList`")
	}

	eduprog.EducationLevel = levelData.Level
	eduprog.Stage = levelData.Stage

	allSpecialties, err := s.specialtiesService.ShowAllSpecialties()
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	check := false
	for i := range allSpecialties {
		if allSpecialties[i].Code == eduprog.SpecialtyCode {
			check = true
			eduprog.Speciality = allSpecialties[i].Name
			eduprog.KFCode = allSpecialties[i].KFCode
			eduprog.KnowledgeField = allSpecialties[i].KnowledgeField
		}
	}

	if !check {
		log.Printf("EduprogService: %s", errors.New("there is no such specialty in enum, only values from `ShowAllSpecialties` can be used"))
		return domain.Eduprog{}, errors.New("there is no such specialty in enum, only values from `ShowAllSpecialties` can be used")
	}

	eduprog, err = s.eduprogRepo.Save(eduprog)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	return eduprog, nil
}

func (s eduprogService) Update(eduprog domain.Eduprog, id uint64) (domain.Eduprog, error) {
	e, err := s.eduprogRepo.Update(eduprog, id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}
	return e, err
}

func (s eduprogService) ShowList() (domain.Eduprogs, error) {
	e, err := s.eduprogRepo.ShowList()
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprogs{}, err
	}
	return e, nil
}

func (s eduprogService) FindById(id uint64) (domain.Eduprog, error) {
	e, err := s.eduprogRepo.FindById(id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}
	return e, nil
}

func (s eduprogService) GetOPPLevelsList() ([]domain.OPPLevelStruct, error) {
	e, err := s.eduprogRepo.GetOPPLevelsList()
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return []domain.OPPLevelStruct{}, err
	}
	return e, nil
}

func (s eduprogService) GetOPPLevelData(level string) (domain.OPPLevelStruct, error) {
	e, err := s.eduprogRepo.GetOPPLevelData(level)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.OPPLevelStruct{}, err
	}
	return e, nil
}

func (s eduprogService) Delete(id uint64) error {
	err := s.eduprogRepo.Delete(id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return err
	}
	return nil
}
