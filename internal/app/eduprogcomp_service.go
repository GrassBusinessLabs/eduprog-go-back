package app

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
	"reflect"
	"sort"
	"strconv"
)

type EduprogcompService interface {
	Save(eduprogcomp domain.Eduprogcomp) (domain.Eduprogcomp, error)
	Update(eduprogcomp domain.Eduprogcomp, id uint64) (domain.Eduprogcomp, error)
	ShowList() ([]domain.Eduprogcomp, error)
	FindById(id uint64) (domain.Eduprogcomp, error)
	FindByWODeleteDate(id uint64) (domain.Eduprogcomp, error)
	FindByBlockNum(id uint64, blockNum string) ([]domain.Eduprogcomp, error)
	SortComponentsByMnS(eduprog_id uint64) (domain.Components, error)
	ShowListByEduprogId(eduprog_id uint64) ([]domain.Eduprogcomp, error)
	ShowListByEduprogIdWithType(eduprog_id uint64, _type string) ([]domain.Eduprogcomp, error)
	Delete(id uint64) error
	GetVBBlocksDomain(eduprogcomps []domain.Eduprogcomp) []domain.BlockInfo
}

type eduprogcompService struct {
	eduprogcompRepo eduprog.EduprogcompRepository
}

func NewEduprogcompService(er eduprog.EduprogcompRepository) EduprogcompService {
	return eduprogcompService{
		eduprogcompRepo: er,
	}
}

func (s eduprogcompService) GetVBBlocksDomain(eduprogcomps []domain.Eduprogcomp) []domain.BlockInfo {
	var blockInfo []domain.BlockInfo
	for i := range eduprogcomps {
		var temp domain.BlockInfo
		temp.BlockNum = eduprogcomps[i].BlockNum
		temp.BlockName = eduprogcomps[i].BlockName
		blockInfo = append(blockInfo, temp)
	}
	blockInfo = RemoveDuplicatesByField(blockInfo, "BlockNum")
	for i := range blockInfo {
		for i2 := range eduprogcomps {
			if blockInfo[i].BlockNum == eduprogcomps[i2].BlockNum {
				blockInfo[i].CompsInBlock = append(blockInfo[i].CompsInBlock, eduprogcomps[i])
			}
		}
		blockInfo[i].CompsInBlock = sortByCode(blockInfo[i].CompsInBlock)
	}
	sortBlocks(blockInfo)
	return blockInfo
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

func (s eduprogcompService) FindByBlockNum(id uint64, blockNum string) ([]domain.Eduprogcomp, error) {
	e, err := s.eduprogcompRepo.FindByBlockNum(id, blockNum)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return []domain.Eduprogcomp{}, err
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

func (s eduprogcompService) ShowListByEduprogIdWithType(eduprogId uint64, _type string) ([]domain.Eduprogcomp, error) {
	e, err := s.eduprogcompRepo.ShowListByEduprogIdWithType(eduprogId, _type)
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

func sortBlocks(blocks []domain.BlockInfo) {
	sort.Slice(blocks, func(i, j int) bool {
		blockNumI, errI := strconv.Atoi(blocks[i].BlockNum)
		blockNumJ, errJ := strconv.Atoi(blocks[j].BlockNum)
		if errI != nil || errJ != nil {
			// handle error cases where blockNum is not an integer
			return false
		}
		return blockNumI < blockNumJ
	})
}

func sortByCode(eduprogcomps []domain.Eduprogcomp) []domain.Eduprogcomp {
	sort.Slice(eduprogcomps, func(i, j int) bool {
		// Parse the Code field as integers and compare them
		codeI, errI := strconv.ParseUint(eduprogcomps[i].Code, 10, 64)
		codeJ, errJ := strconv.ParseUint(eduprogcomps[j].Code, 10, 64)
		if errI != nil || errJ != nil {
			return eduprogcomps[i].Code < eduprogcomps[j].Code
		}
		return codeI < codeJ
	})
	return eduprogcomps
}

func RemoveDuplicatesByField(mySlice []domain.BlockInfo, fieldName string) []domain.BlockInfo {
	unique := make(map[string]bool)
	result := make([]domain.BlockInfo, 0)
	for _, v := range mySlice {
		fieldValue := reflect.ValueOf(v).FieldByName(fieldName).String()
		if !unique[fieldValue] {
			unique[fieldValue] = true
			result = append(result, v)
		}
	}
	return result
}
