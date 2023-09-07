package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
)

type DisciplineService interface {
	Save(discipline domain.Discipline) (domain.Discipline, error)
	Update(ref, req domain.Discipline) (domain.Discipline, error)
	ShowDisciplinesByEduprogId(eduprog_id uint64) ([]domain.Discipline, error)
	FindById(id uint64) (domain.Discipline, error)
	AddRow(disciplineId uint64) (domain.Discipline, error)
	Delete(id uint64) error
}

type disciplineService struct {
	disciplineRepo eduprog.DisciplineRepository
}

func NewDisciplineService(dr eduprog.DisciplineRepository) DisciplineService {
	return disciplineService{
		disciplineRepo: dr,
	}
}

func (s disciplineService) Save(discipline domain.Discipline) (domain.Discipline, error) {
	e, err := s.disciplineRepo.Save(discipline)
	if err != nil {
		log.Printf("DisciplineService: %s", err)
		return domain.Discipline{}, err
	}
	return e, err
}

func (s disciplineService) Update(ref, req domain.Discipline) (domain.Discipline, error) {
	if req.Name != "" {
		ref.Name = req.Name
	}
	if req.EduprogId != 0 {
		ref.EduprogId = req.EduprogId
	}

	return s.disciplineRepo.Update(ref)
}

func (s disciplineService) ShowDisciplinesByEduprogId(eduprog_id uint64) ([]domain.Discipline, error) {
	e, err := s.disciplineRepo.ShowDisciplinesByEduprogId(eduprog_id)
	if err != nil {
		log.Printf("DisciplineService: %s", err)
		return []domain.Discipline{}, err
	}
	return e, nil
}

func (s disciplineService) FindById(id uint64) (domain.Discipline, error) {
	e, err := s.disciplineRepo.FindById(id)
	if err != nil {
		log.Printf("DisciplineService: %s", err)
		return domain.Discipline{}, err
	}
	return e, nil
}

func (s disciplineService) AddRow(disciplineId uint64) (domain.Discipline, error) {
	discipline, err := s.FindById(disciplineId)
	if err != nil {
		log.Printf("DisciplineService: %s", err)
		return domain.Discipline{}, err
	}

	req := discipline
	req.Rows = req.Rows + 1

	discipline, err = s.Update(discipline, req)
	if err != nil {
		log.Printf("DisciplineService: %s", err)
		return domain.Discipline{}, err
	}
	return discipline, nil
}

func (s disciplineService) Delete(id uint64) error {
	err := s.disciplineRepo.Delete(id)
	if err != nil {
		log.Printf("DisciplineService: %s", err)
		return err
	}
	return nil
}
