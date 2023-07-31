package app

import (
	"errors"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"log"
	"sort"
	"strconv"
)

type EduprogcompService interface {
	Save(eduprogcomp domain.Eduprogcomp) (domain.Eduprogcomp, error)
	Update(eduprogcomp domain.Eduprogcomp) (domain.Eduprogcomp, error)
	FindById(id uint64) (domain.Eduprogcomp, error)
	FindByBlockNum(id uint64, blockNum string) ([]domain.Eduprogcomp, error)
	SortComponentsByMnS(eduprogId uint64) (domain.Components, error)
	ShowListByEduprogId(eduprogId uint64) ([]domain.Eduprogcomp, error)
	ShowListByEduprogIdWithType(eduprogId uint64, _type string) ([]domain.Eduprogcomp, error)
	UpdateVBName(eduprogId uint64, eduprogcompReq domain.Eduprogcomp) ([]domain.Eduprogcomp, error)
	Delete(id uint64) error
	ReplaceOK(eduprogcompId, putAfter uint64) (domain.Components, error)
	ReplaceVBBlock(firstCompId uint64, putAfter uint64) (domain.Components, error)
	ReplaceVB(eduprogcompId, blockNum, putAfter uint64) (domain.Components, error)
}

type eduprogcompService struct {
	eduprogcompRepo eduprog.EduprogcompRepository
}

func NewEduprogcompService(er eduprog.EduprogcompRepository) EduprogcompService {
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

func (s eduprogcompService) Update(eduprogcomp domain.Eduprogcomp) (domain.Eduprogcomp, error) {
	e, err := s.eduprogcompRepo.Update(eduprogcomp)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, err
	}
	return e, err
}

func (s eduprogcompService) FindById(id uint64) (domain.Eduprogcomp, error) {
	e, err := s.eduprogcompRepo.FindById(id)
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

func (s eduprogcompService) SortComponentsByMnS(eduprogId uint64) (domain.Components, error) {
	e, err := s.eduprogcompRepo.SortComponentsByMnS(eduprogId)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Components{}, err
	}
	return e, nil
}

func (s eduprogcompService) ShowListByEduprogId(eduprogId uint64) ([]domain.Eduprogcomp, error) {
	e, err := s.eduprogcompRepo.ShowListByEduprogId(eduprogId)
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

func (s eduprogcompService) UpdateVBName(eduprogId uint64, eduprogcompReq domain.Eduprogcomp) ([]domain.Eduprogcomp, error) {
	eduprogcomps, err := s.SortComponentsByMnS(eduprogId)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return []domain.Eduprogcomp{}, err
	}

	for i := range eduprogcomps.Selective {
		if eduprogcompReq.BlockName == eduprogcomps.Selective[i].BlockName && eduprogcompReq.BlockNum != eduprogcomps.Selective[i].BlockNum {
			log.Printf("EduprogcompService: %s", err)
			return []domain.Eduprogcomp{}, errors.New("block with this name already exists")
		}
	}

	vbBlock, err := s.FindByBlockNum(eduprogId, eduprogcompReq.BlockNum)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return []domain.Eduprogcomp{}, err
	}

	var result []domain.Eduprogcomp
	for i := range vbBlock {
		edcompById, err := s.FindById(vbBlock[i].Id)
		if err != nil {
			log.Printf("EduprogcompService: %s", err)
			return []domain.Eduprogcomp{}, err
		}
		edcompById.BlockName = eduprogcompReq.BlockName
		updEduprogcomp, err := s.Update(edcompById)
		if err != nil {
			log.Printf("EduprogcompService: %s", err)
			return []domain.Eduprogcomp{}, err
		}

		result = append(result, updEduprogcomp)
	}

	return result, nil
}

func (s eduprogcompService) Delete(id uint64) error {
	eduprogcomp, err := s.FindById(id)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return err
	}

	err = s.eduprogcompRepo.Delete(id)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return err
	}

	eduprogcomps, err := s.SortComponentsByMnS(eduprogcomp.EduprogId)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return err
	}

	for i, elem := range eduprogcomps.Mandatory {
		elem.Code = strconv.Itoa(i + 1)
		eduprogcomps.Mandatory[i] = elem
	}

	for i, elem := range eduprogcomps.Selective {
		elem.BlockNum = strconv.Itoa(i + 1)
		eduprogcomps.Selective[i] = elem
	}

	for i := range eduprogcomps.Mandatory {
		_, err = s.Update(eduprogcomps.Mandatory[i])
		if err != nil {
			log.Printf("EduprogcompService: %s", err)
			return err
		}

	}

	for i := range eduprogcomps.Selective {
		for _, comp := range eduprogcomps.Selective[i].CompsInBlock {
			comp.BlockNum = eduprogcomps.Selective[i].BlockNum
			_, err = s.Update(comp)
			if err != nil {
				log.Printf("EduprogcompService: %s", err)
				return err
			}

		}
	}

	return nil
}

func (s eduprogcompService) ReplaceOK(eduprogcompId, putAfter uint64) (domain.Components, error) {
	educompById, err := s.FindById(eduprogcompId)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Components{}, err
	}

	eduprogcomps, err := s.SortComponentsByMnS(educompById.EduprogId)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Components{}, err
	}

	fromIndex, err := strconv.Atoi(educompById.Code)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Components{}, err
	}

	eduprogcomps.Mandatory = s.moveEduprogcomp(eduprogcomps.Mandatory, fromIndex-1, int(putAfter))

	for i := range eduprogcomps.Mandatory {
		eduprogcomps.Mandatory[i].Code = strconv.Itoa(i + 1)
		_, err = s.Update(eduprogcomps.Mandatory[i])
		if err != nil {
			log.Printf("EduprogcompService: %s", err)
			return domain.Components{}, err
		}
	}

	return eduprogcomps, nil
}

func (s eduprogcompService) ReplaceVBBlock(firstCompId uint64, putAfter uint64) (domain.Components, error) {
	eduprogcompById, err := s.eduprogcompRepo.FindById(firstCompId)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Components{}, err
	}

	eduprogcomps, err := s.eduprogcompRepo.ShowListByEduprogIdWithType(eduprogcompById.EduprogId, "ВБ")
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Components{}, err
	}

	var components domain.Components
	var block domain.BlockInfo

	for _, eduprogcomp := range eduprogcomps {
		block.BlockName = eduprogcomp.BlockName
		block.BlockNum = eduprogcomp.BlockNum
		components.Selective = append(components.Selective, block)
	}
	components.Selective = s.uniqueBlocks(components.Selective)

	for i, info := range components.Selective {
		for _, eduprogcomp := range eduprogcomps {
			if eduprogcomp.BlockNum == info.BlockNum {
				components.Selective[i].CompsInBlock = append(components.Selective[i].CompsInBlock, eduprogcomp)
			}
		}
	}

	s.sortBlocks(components.Selective)

	fromIndex, err := strconv.Atoi(eduprogcompById.BlockNum)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Components{}, err
	}

	components.Selective = s.moveVBBlock(components.Selective, fromIndex-1, int(putAfter))

	for i, blockInfo := range components.Selective {
		components.Selective[i].BlockNum = strconv.Itoa(i + 1)
		for i2 := range blockInfo.CompsInBlock {
			components.Selective[i].CompsInBlock[i2].BlockNum = components.Selective[i].BlockNum
			_, err = s.Update(components.Selective[i].CompsInBlock[i2])
			if err != nil {
				log.Printf("EduprogcompService: %s", err)
				return domain.Components{}, err
			}
		}
	}

	return components, nil
}

func (s eduprogcompService) ReplaceVB(eduprogcompId, blockNum, putAfter uint64) (domain.Components, error) {
	eduprogcompById, err := s.eduprogcompRepo.FindById(eduprogcompId)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Components{}, err
	}

	eduprogcomps, err := s.eduprogcompRepo.ShowListByEduprogIdWithType(eduprogcompById.EduprogId, "ВБ")
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Components{}, err
	}

	var components domain.Components
	var block domain.BlockInfo

	for _, eduprogcomp := range eduprogcomps {
		block.BlockName = eduprogcomp.BlockName
		block.BlockNum = eduprogcomp.BlockNum
		components.Selective = append(components.Selective, block)
	}
	components.Selective = s.uniqueBlocks(components.Selective)

	for i, info := range components.Selective {
		for _, eduprogcomp := range eduprogcomps {
			if eduprogcomp.BlockNum == info.BlockNum {
				if eduprogcomp.Id != eduprogcompById.Id {
					components.Selective[i].CompsInBlock = append(components.Selective[i].CompsInBlock, eduprogcomp)
				}
			}
		}
	}

	s.sortBlocks(components.Selective)

	for i, blockInfo := range components.Selective {
		if blockInfo.BlockNum == strconv.FormatUint(blockNum, 10) {
			components.Selective[i].CompsInBlock = append(components.Selective[i].CompsInBlock, eduprogcompById)
			components.Selective[i].CompsInBlock = s.moveEduprogcomp(components.Selective[i].CompsInBlock, len(components.Selective[i].CompsInBlock)-1, int(putAfter))
		}
	}

	for i, blockInfo := range components.Selective {
		components.Selective[i].BlockNum = strconv.Itoa(i + 1)
		for i2 := range blockInfo.CompsInBlock {
			components.Selective[i].CompsInBlock[i2].Code = strconv.Itoa(i2 + 1)
			components.Selective[i].CompsInBlock[i2].BlockNum = components.Selective[i].BlockNum
			components.Selective[i].CompsInBlock[i2].BlockName = components.Selective[i].BlockName
			_, err = s.Update(components.Selective[i].CompsInBlock[i2])
			if err != nil {
				log.Printf("EduprogcompService: %s", err)
				return domain.Components{}, err
			}
		}
	}

	return components, nil
}

func (s eduprogcompService) moveEduprogcomp(slice []domain.Eduprogcomp, fromIndex, toIndex int) []domain.Eduprogcomp {
	element := slice[fromIndex]
	slice = append(slice[:fromIndex], slice[fromIndex+1:]...)
	slice = append(slice[:toIndex], append([]domain.Eduprogcomp{element}, slice[toIndex:]...)...)
	return slice
}

func (s eduprogcompService) moveVBBlock(slice []domain.BlockInfo, fromIndex, toIndex int) []domain.BlockInfo {
	element := slice[fromIndex]
	slice = append(slice[:fromIndex], slice[fromIndex+1:]...)
	slice = append(slice[:toIndex], append([]domain.BlockInfo{element}, slice[toIndex:]...)...)
	return slice
}

func (s eduprogcompService) uniqueBlocks(blocks []domain.BlockInfo) []domain.BlockInfo {
	var unique []domain.BlockInfo

loop:
	for _, l := range blocks {
		for i, u := range unique {
			if l.BlockName == u.BlockName {
				unique[i] = l
				continue loop
			}
		}
		unique = append(unique, l)
	}

	return unique
}

func (s eduprogcompService) sortBlocks(blocks []domain.BlockInfo) {
	sort.Slice(blocks, func(i, j int) bool {
		blockNumI, errI := strconv.Atoi(blocks[i].BlockNum)
		blockNumJ, errJ := strconv.Atoi(blocks[j].BlockNum)
		if errI != nil || errJ != nil {
			return false
		}
		return blockNumI < blockNumJ
	})
}
