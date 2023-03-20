package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
)

type SpecialtiesService interface {
	CreateSpecialty(specialty domain.Specialty) (domain.Specialty, error)
	UpdateSpecialty(specialty domain.Specialty, code string) (domain.Specialty, error)
	ShowAllSpecialties() ([]domain.Specialty, error)
	ShowByKFCode(kfCode string) ([]domain.Specialty, error)
	FindByCode(code string) (domain.Specialty, error)
	Delete(code string) error
}

type specialtiesService struct {
	specialtiesRepo eduprog.SpecialtiesRepository
}

func NewSpecialtiesService(sr eduprog.SpecialtiesRepository) SpecialtiesService {
	return specialtiesService{
		specialtiesRepo: sr,
	}
}

func (s specialtiesService) CreateSpecialty(specialty domain.Specialty) (domain.Specialty, error) {
	e, err := s.specialtiesRepo.CreateSpecialty(specialty)
	if err != nil {
		log.Printf("SpecialtiesService: %s", err)
		return domain.Specialty{}, err
	}
	return e, err
}

func (s specialtiesService) UpdateSpecialty(specialty domain.Specialty, code string) (domain.Specialty, error) {
	e, err := s.specialtiesRepo.UpdateSpecialty(specialty, code)
	if err != nil {
		log.Printf("SpecialtiesService: %s", err)
		return domain.Specialty{}, err
	}
	return e, err
}

func (s specialtiesService) ShowAllSpecialties() ([]domain.Specialty, error) {
	e, err := s.specialtiesRepo.ShowAllSpecialties()
	if err != nil {
		log.Printf("SpecialtiesService: %s", err)
		return []domain.Specialty{}, err
	}
	return e, err
}

func (s specialtiesService) ShowByKFCode(kfCode string) ([]domain.Specialty, error) {
	e, err := s.specialtiesRepo.ShowByKFCode(kfCode)
	if err != nil {
		log.Printf("SpecialtiesService: %s", err)
		return []domain.Specialty{}, err
	}
	return e, err
}

func (s specialtiesService) FindByCode(code string) (domain.Specialty, error) {
	e, err := s.specialtiesRepo.FindByCode(code)
	if err != nil {
		log.Printf("SpecialtiesService: %s", err)
		return domain.Specialty{}, err
	}
	return e, err
}

func (s specialtiesService) Delete(code string) error {
	err := s.specialtiesRepo.Delete(code)
	if err != nil {
		log.Printf("SpecialtiesService: %s", err)
		return err
	}
	return err
}
