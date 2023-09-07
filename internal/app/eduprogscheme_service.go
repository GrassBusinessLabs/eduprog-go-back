package app

import (
	"fmt"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
)

type EduprogschemeService interface {
	SetComponentToEdprogscheme(eduprogscheme domain.Eduprogscheme) (domain.Eduprogscheme, error)
	UpdateComponentInEduprogscheme(eduprogscheme domain.Eduprogscheme, id uint64) (domain.Eduprogscheme, error)
	ExpandOrShrinkEduprogschemeComponent(eduprogcomp_id uint64, semNum uint64, direction string) ([]domain.Eduprogscheme, error)
	FindById(id uint64) (domain.Eduprogscheme, error)
	FindBySemesterNum(semester_num uint16, eduprog_id uint64) ([]domain.Eduprogscheme, error)
	ShowSchemeByEduprogId(eduprog_id uint64) ([]domain.Eduprogscheme, error)
	Delete(id uint64) error
}

type eduprogschemeService struct {
	eduprogschemeRepo  eduprog.EduprogschemeRepository
	eduprogcompService EduprogcompService
}

func NewEduprogschemeService(er eduprog.EduprogschemeRepository, ess EduprogcompService) EduprogschemeService {
	return eduprogschemeService{
		eduprogschemeRepo:  er,
		eduprogcompService: ess,
	}
}

func (s eduprogschemeService) SetComponentToEdprogscheme(eduprogscheme domain.Eduprogscheme) (domain.Eduprogscheme, error) {
	e, err := s.eduprogschemeRepo.SetComponentToEdprogscheme(eduprogscheme)
	if err != nil {
		log.Printf("EduprogschemeService: %s", err)
		return domain.Eduprogscheme{}, err
	}
	return e, err
}

func (s eduprogschemeService) UpdateComponentInEduprogscheme(eduprogscheme domain.Eduprogscheme, id uint64) (domain.Eduprogscheme, error) {
	e, err := s.eduprogschemeRepo.UpdateComponentInEduprogscheme(eduprogscheme, id)
	if err != nil {
		log.Printf("EduprogschemeService: %s", err)
		return domain.Eduprogscheme{}, err
	}
	return e, err
}

func (s eduprogschemeService) ExpandOrShrinkEduprogschemeComponent(eduprogcompId uint64, semNum uint64, direction string) ([]domain.Eduprogscheme, error) {
	if semNum < 1 && semNum > 8 {
		err := fmt.Errorf("invalid semester num, must be from 1 to 8")
		log.Printf("EduprogschemeService: %s", err)
		return []domain.Eduprogscheme{}, err
	}

	eduprogcomp, err := s.eduprogcompService.FindById(eduprogcompId)
	if err != nil {
		log.Printf("EduprogschemeService: %s", err)
		return []domain.Eduprogscheme{}, err
	}

	eduprogscheme, err := s.ShowSchemeByEduprogId(eduprogcomp.EduprogId)
	if err != nil {
		log.Printf("EduprogschemeService: %s", err)
		return []domain.Eduprogscheme{}, err
	}

	var schemeComponentsList []domain.Eduprogscheme

	for i := range eduprogscheme {
		if eduprogscheme[i].EduprogcompId == eduprogcomp.Id {
			schemeComponentsList = append(schemeComponentsList, eduprogscheme[i])
		}
	}
	if len(schemeComponentsList) == 0 {
		err = fmt.Errorf("eduprogcomp is not in scheme")
		log.Printf("EduprogschemeService: %s", err)
		return []domain.Eduprogscheme{}, err
	}

	startingSchemeCompsListLen := len(schemeComponentsList)

	minSemesterNum := schemeComponentsList[0].SemesterNum
	maxSemesterNum := schemeComponentsList[len(schemeComponentsList)-1].SemesterNum
	var creditsPerSemester float64

	if semNum < minSemesterNum {
		schemeComponentToCreate := schemeComponentsList[0]
		for i := semNum; i < minSemesterNum; i++ {
			schemeComponentToCreate.SemesterNum = schemeComponentToCreate.SemesterNum - 1
			schemeComponentToCreate.Id = 0
			schemeComponentsList = append(schemeComponentsList, schemeComponentToCreate)
		}

		creditsPerSemester = eduprogcomp.Credits / float64(len(schemeComponentsList))

		if creditsPerSemester < 3 {
			err = fmt.Errorf("cannot expand component, must be 3 credits in each semester")
			log.Printf("EduprogschemeService: %s", err)
			return []domain.Eduprogscheme{}, err
		}

		for i := startingSchemeCompsListLen; i < len(schemeComponentsList); i++ {
			_, err = s.SetComponentToEdprogscheme(schemeComponentsList[i])
			if err != nil {
				log.Printf("EduprogschemeService: %s", err)
				return []domain.Eduprogscheme{}, err
			}
		}
	} else if semNum > maxSemesterNum {
		schemeComponentToCreate := schemeComponentsList[0]
		for i := maxSemesterNum; i < semNum; i++ {
			schemeComponentToCreate.SemesterNum = schemeComponentToCreate.SemesterNum + 1
			schemeComponentToCreate.Id = 0
			schemeComponentsList = append(schemeComponentsList, schemeComponentToCreate)
		}

		creditsPerSemester = eduprogcomp.Credits / float64(len(schemeComponentsList))

		if creditsPerSemester < 3 {
			err = fmt.Errorf("cannot expand component, must be 3 credits in each semester")
			log.Printf("EduprogschemeService: %s", err)
			return []domain.Eduprogscheme{}, err
		}

		for i := startingSchemeCompsListLen; i < len(schemeComponentsList); i++ {
			schemeComponentsList[i].CreditsPerSemester = creditsPerSemester
			_, err = s.SetComponentToEdprogscheme(schemeComponentsList[i])
			if err != nil {
				log.Printf("EduprogschemeService: %s", err)
				return []domain.Eduprogscheme{}, err
			}
		}
	} else if semNum >= minSemesterNum && semNum <= maxSemesterNum {
		if direction == "RIGHT" {
			for i := range schemeComponentsList {
				if schemeComponentsList[i].SemesterNum < semNum {
					err = s.Delete(schemeComponentsList[i].Id)
					if err != nil {
						log.Printf("EduprogschemeService: %s", err)
						return []domain.Eduprogscheme{}, err
					}
				}
			}
		} else if direction == "LEFT" {
			for i := range schemeComponentsList {
				if schemeComponentsList[i].SemesterNum > semNum {
					err = s.Delete(schemeComponentsList[i].Id)
					if err != nil {
						log.Printf("EduprogschemeService: %s", err)
						return []domain.Eduprogscheme{}, err
					}
				}
			}
		}

	}

	eduprogscheme, err = s.ShowSchemeByEduprogId(eduprogcomp.EduprogId)
	if err != nil {
		log.Printf("EduprogschemeService: %s", err)
		return []domain.Eduprogscheme{}, err
	}

	schemeComponentsList = []domain.Eduprogscheme{}

	for i := range eduprogscheme {
		if eduprogscheme[i].EduprogcompId == eduprogcomp.Id {
			if minSemesterNum > eduprogscheme[i].SemesterNum {
				minSemesterNum = eduprogscheme[i].SemesterNum
			}
			schemeComponentsList = append(schemeComponentsList, eduprogscheme[i])
		}
	}

	creditsPerSemester = eduprogcomp.Credits / float64(len(schemeComponentsList))
	for i := range schemeComponentsList {
		schemeComponentsList[i].CreditsPerSemester = creditsPerSemester
		if semNum >= minSemesterNum && semNum <= maxSemesterNum {
			if direction == "RIGHT" {
				schemeComponentsList[i].SemesterNum = semNum
				semNum++
			} else if direction == "LEFT" {
				schemeComponentsList[i].SemesterNum = minSemesterNum
				minSemesterNum++
			}
		} else {
			schemeComponentsList[i].SemesterNum = minSemesterNum
			minSemesterNum++
		}
		_, err = s.UpdateComponentInEduprogscheme(schemeComponentsList[i], schemeComponentsList[i].Id)
		if err != nil {
			log.Printf("EduprogschemeService: %s", err)
			return []domain.Eduprogscheme{}, err
		}
	}

	eduprogschemeToShow, err := s.ShowSchemeByEduprogId(eduprogcomp.EduprogId)
	if err != nil {
		log.Printf("EduprogschemeService: %s", err)
		return []domain.Eduprogscheme{}, err
	}

	return eduprogschemeToShow, nil
}

func (s eduprogschemeService) FindById(id uint64) (domain.Eduprogscheme, error) {
	e, err := s.eduprogschemeRepo.FindById(id)
	if err != nil {
		log.Printf("EduprogschemeService: %s", err)
		return domain.Eduprogscheme{}, err
	}
	return e, nil
}

func (s eduprogschemeService) FindBySemesterNum(semesterNum uint16, eduprogId uint64) ([]domain.Eduprogscheme, error) {
	e, err := s.eduprogschemeRepo.FindBySemesterNum(semesterNum, eduprogId)
	if err != nil {
		log.Printf("EduprogschemeService: %s", err)
		return []domain.Eduprogscheme{}, err
	}
	return e, nil
}

func (s eduprogschemeService) ShowSchemeByEduprogId(eduprog_id uint64) ([]domain.Eduprogscheme, error) {
	e, err := s.eduprogschemeRepo.ShowSchemeByEduprogId(eduprog_id)
	if err != nil {
		log.Printf("EduprogschemeService: %s", err)
		return []domain.Eduprogscheme{}, err
	}
	return e, nil
}

func (s eduprogschemeService) Delete(id uint64) error {
	err := s.eduprogschemeRepo.Delete(id)
	if err != nil {
		log.Printf("EduprogschemeService: %s", err)
		return err
	}
	return nil
}
