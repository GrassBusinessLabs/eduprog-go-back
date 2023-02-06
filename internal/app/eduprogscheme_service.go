package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database"
	"log"
)

type EduprogschemeService interface {
	SetComponentToEdprogscheme(eduprogscheme domain.Eduprogscheme) (domain.Eduprogscheme, error)
	UpdateComponentInEduprogscheme(eduprogscheme domain.Eduprogscheme, id uint64) (domain.Eduprogscheme, error)
	FindById(id uint64) (domain.Eduprogscheme, error)
	FindBySemesterNum(semester_num uint16) ([]domain.Eduprogscheme, error)
	ShowSchemeByEduprogId(eduprog_id uint64) ([]domain.Eduprogscheme, error)
	Delete(id uint64) error
}

type eduprogschemeService struct {
	eduprogschemeRepo database.EduprogschemeRepository
}

func NewEduprogschemeService(er database.EduprogschemeRepository) EduprogschemeService {
	return eduprogschemeService{
		eduprogschemeRepo: er,
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

func (s eduprogschemeService) FindById(id uint64) (domain.Eduprogscheme, error) {
	e, err := s.eduprogschemeRepo.FindById(id)
	if err != nil {
		log.Printf("EduprogschemeService: %s", err)
		return domain.Eduprogscheme{}, err
	}
	return e, nil
}

func (s eduprogschemeService) FindBySemesterNum(semester_num uint16) ([]domain.Eduprogscheme, error) {
	e, err := s.eduprogschemeRepo.FindBySemesterNum(semester_num)
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
