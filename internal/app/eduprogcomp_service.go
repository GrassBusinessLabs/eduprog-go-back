package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database"
	"log"
)

type EduprogcompService interface {
	Save(eduprogcomp domain.Eduprogcomp) (domain.Eduprogcomp, error)
	Update(eduprogcomp domain.Eduprogcomp, id uint64) (domain.Eduprogcomp, error)
	ShowList() ([]domain.Eduprogcomp, error)
	FindById(id uint64) (domain.Eduprogcomp, error)
	FindByWODeleteDate(id uint64) (domain.Eduprogcomp, error)
	SortComponentsByMnS(eduprog_id uint64) (domain.Components, error)
	ShowListByEduprogId(eduprog_id uint64) ([]domain.Eduprogcomp, error)
	Delete(id uint64) error
}

type eduprogcompService struct {
	eduprogcompRepo database.EduprogcompRepository
}

func NewEduprogcompService(er database.EduprogcompRepository) EduprogcompService {
	return eduprogcompService{
		eduprogcompRepo: er,
	}
}

func (s eduprogcompService) Save(eduprogcomp domain.Eduprogcomp) (domain.Eduprogcomp, error) {
	e, err := s.eduprogcompRepo.Save(eduprogcomp)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, err
	}
	return e, err
}

func (s eduprogcompService) Update(eduprogcomp domain.Eduprogcomp, id uint64) (domain.Eduprogcomp, error) {
	e, err := s.eduprogcompRepo.Update(eduprogcomp, id)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, err
	}
	return e, err
}

func (s eduprogcompService) ShowList() ([]domain.Eduprogcomp, error) {
	e, err := s.eduprogcompRepo.ShowList()
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return []domain.Eduprogcomp{}, err
	}
	return e, nil
}

func (s eduprogcompService) FindById(id uint64) (domain.Eduprogcomp, error) {
	e, err := s.eduprogcompRepo.FindById(id)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, err
	}
	return e, nil
}
func (s eduprogcompService) FindByWODeleteDate(id uint64) (domain.Eduprogcomp, error) {
	e, err := s.eduprogcompRepo.FindByWODeleteDate(id)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, err
	}
	return e, nil
}

func (s eduprogcompService) SortComponentsByMnS(eduprog_id uint64) (domain.Components, error) {
	e, err := s.eduprogcompRepo.SortComponentsByMnS(eduprog_id)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Components{}, err
	}
	return e, nil
}

func (s eduprogcompService) ShowListByEduprogId(eduprog_id uint64) ([]domain.Eduprogcomp, error) {
	e, err := s.eduprogcompRepo.ShowListByEduprogId(eduprog_id)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return []domain.Eduprogcomp{}, err
	}
	return e, nil
}

func (s eduprogcompService) Delete(id uint64) error {
	err := s.eduprogcompRepo.Delete(id)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return err
	}
	return nil
}
