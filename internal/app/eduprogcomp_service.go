package app

import (
	"errors"
	"fmt"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/resources"
	"log"
	"sort"
	"strconv"
)

type EduprogcompService interface {
	Save(eduprogcomp domain.Eduprogcomp) (domain.Eduprogcomp, error)
	Update(ref, req domain.Eduprogcomp) (domain.Eduprogcomp, error)
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
	GetCreditsInfo(eduprog domain.Eduprog) (resources.CreditsDto, error)
}

type eduprogcompService struct {
	eduprogcompRepo eduprog.EduprogcompRepository
	eduprogService  EduprogService
}

func NewEduprogcompService(er eduprog.EduprogcompRepository, es EduprogService) EduprogcompService {
	return eduprogcompService{
		eduprogcompRepo: er,
		eduprogService:  es,
	}
}

func (s eduprogcompService) Save(eduprogcomp domain.Eduprogcomp) (domain.Eduprogcomp, error) {
	exists, err := s.eduprogcompRepo.CheckName(eduprogcomp)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, err
	}
	if exists {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, fmt.Errorf("eduprogcomp with name '%s' already exists in VB block/OK list", eduprogcomp.Name)
	}

	eduprogById, err := s.eduprogService.FindById(eduprogcomp.EduprogId)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, err
	}

	creditsDto, err := s.GetCreditsInfo(eduprogById)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, err
	}

	if eduprogcomp.Type == "ОК" {
		if eduprogcomp.Credits > creditsDto.MandatoryFreeCredits {
			log.Printf("EduprogcompService: %s", err)
			return domain.Eduprogcomp{}, fmt.Errorf("too much credits")
		}
	} else if eduprogcomp.Type == "ВБ" {
		if eduprogcomp.Credits > creditsDto.SelectiveFreeCredits {
			log.Printf("EduprogcompService: %s", err)
			return domain.Eduprogcomp{}, fmt.Errorf("too much credits")
		}
	} else {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, fmt.Errorf("only 'ВБ' or 'ОК'")
	}

	maxCode := 0

	eduprogcomps, err := s.SortComponentsByMnS(eduprogcomp.EduprogId)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, err
	}

	if eduprogcomp.Type == "ОК" {
		for i := range eduprogcomps.Mandatory {
			code, err := strconv.Atoi(eduprogcomps.Mandatory[i].Code)
			if err != nil {
				log.Printf("EduprogcompService: %s", err)
				return domain.Eduprogcomp{}, err
			}
			if maxCode < code {
				maxCode = code
			}
		}
	} else if eduprogcomp.Type == "ВБ" {
		for i := range eduprogcomps.Selective {
			if eduprogcomps.Selective[i].BlockNum == eduprogcomp.BlockNum {
				for i2 := range eduprogcomps.Selective[i].CompsInBlock {
					code, err := strconv.Atoi(eduprogcomps.Selective[i].CompsInBlock[i2].Code)
					if err != nil {
						log.Printf("EduprogcompService: %s", err)
						return domain.Eduprogcomp{}, err
					}
					if maxCode < code {
						maxCode = code
					}
				}
			}
		}
	}

	eduprogcomp.Code = strconv.Itoa(maxCode + 1)

	e, err := s.eduprogcompRepo.Save(eduprogcomp)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, err
	}

	err = s.codesRedefine(eduprogcomp.EduprogId)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, err
	}

	return e, err
}

func (s eduprogcompService) Update(ref, req domain.Eduprogcomp) (domain.Eduprogcomp, error) {
	if req.Name != "" {
		exists, err := s.eduprogcompRepo.CheckName(req)
		if err != nil {
			log.Printf("EduprogcompService: %s", err)
			return domain.Eduprogcomp{}, err
		}
		if exists {
			return domain.Eduprogcomp{}, fmt.Errorf("eduprogcomp with name '%s' already exists in VB block/OK list", req.Name)
		}
		ref.Name = req.Name
	}
	if req.Credits != 0 {
		eduprogById, err := s.eduprogService.FindById(req.EduprogId)
		if err != nil {
			log.Printf("EduprogcompService: %s", err)
			return domain.Eduprogcomp{}, err
		}

		creditsDto, err := s.GetCreditsInfo(eduprogById)
		if err != nil {
			log.Printf("EduprogcompService: %s", err)
			return domain.Eduprogcomp{}, err
		}

		if req.Type == "ОК" {
			if req.Credits+(creditsDto.MandatoryCredits-ref.Credits) > creditsDto.MandatoryCreditsForLevel {
				log.Printf("EduprogcompService: %s", err)
				return domain.Eduprogcomp{}, fmt.Errorf("too much credits")
			}
		} else if req.Type == "ВБ" {
			if req.Credits+(creditsDto.SelectiveCredits-ref.Credits) > creditsDto.SelectiveCreditsForLevel {
				log.Printf("EduprogcompService: %s", err)
				return domain.Eduprogcomp{}, fmt.Errorf("too much credits")
			}
		} else {
			log.Printf("EduprogcompService: %s", err)
			return domain.Eduprogcomp{}, fmt.Errorf("wrong type, only 'ВБ' or 'ОК'")
		}
		ref.Credits = req.Credits
	}
	if req.ControlType != "" {
		ref.ControlType = req.ControlType
	}
	if req.Type != "" {
		ref.Type = req.Type
	}
	if req.BlockNum != "" {
		ref.BlockNum = req.BlockNum
	}
	if req.BlockName != "" {
		ref.BlockName = req.BlockName
	}
	if req.EduprogId != 0 {
		ref.EduprogId = req.EduprogId
	}

	e, err := s.eduprogcompRepo.Update(ref)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, err
	}

	err = s.codesRedefine(req.EduprogId)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return domain.Eduprogcomp{}, err
	}

	return e, err
}

func (s eduprogcompService) GetCreditsInfo(eduprog domain.Eduprog) (resources.CreditsDto, error) {
	var creditsDto resources.CreditsDto

	comps, err := s.SortComponentsByMnS(eduprog.Id)
	if err != nil {
		log.Printf("EduprogController: %s", err)
		return creditsDto, err
	}

	levelData, err := s.eduprogService.GetOPPLevelData(eduprog.EducationLevel)
	if err != nil {
		log.Printf("EduprogController: %s", err)
		return creditsDto, err
	}

	for i := range comps.Selective {
		var minFromBlock domain.Eduprogcomp
		var maxFromBlock domain.Eduprogcomp
		minFromBlock.Credits = 500
		maxFromBlock.Credits = 0
		for _, comp := range comps.Selective[i].CompsInBlock {
			creditsDto.SelectiveCredits += comp.Credits
			if comp.Credits < minFromBlock.Credits {
				minFromBlock = comp
			}
			if comp.Credits > minFromBlock.Credits {
				maxFromBlock = comp
			}
		}
		creditsDto.MinCreditsForVB += minFromBlock.Credits
		creditsDto.MaxCreditsForVB += maxFromBlock.Credits
	}

	for _, comp := range comps.Mandatory {
		creditsDto.MandatoryCredits += comp.Credits
	}

	creditsDto.MandatoryCreditsForLevel = levelData.MandatoryCredits
	creditsDto.SelectiveCreditsForLevel = levelData.SelectiveCredits
	creditsDto.TotalCredits = creditsDto.SelectiveCredits + creditsDto.MandatoryCredits
	creditsDto.TotalFreeCredits = (creditsDto.MandatoryCreditsForLevel + creditsDto.SelectiveCreditsForLevel) - creditsDto.TotalCredits
	creditsDto.MandatoryFreeCredits = creditsDto.MandatoryCreditsForLevel - creditsDto.MandatoryCredits
	creditsDto.SelectiveFreeCredits = creditsDto.SelectiveCreditsForLevel - creditsDto.SelectiveCredits

	return creditsDto, nil
}

func (s eduprogcompService) codesRedefine(eduprogId uint64) error {
	eduprogcomps, err := s.SortComponentsByMnS(eduprogId)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return err
	}

	for i := range eduprogcomps.Mandatory {
		eduprogcomps.Mandatory[i].Code = strconv.Itoa(i + 1)
		_, err = s.eduprogcompRepo.Update(eduprogcomps.Mandatory[i])
		if err != nil {
			log.Printf("EduprogcompService: %s", err)
			return err
		}
	}

	for i := range eduprogcomps.Selective {
		eduprogcomps.Selective[i].BlockNum = strconv.Itoa(i + 1)
		for i2 := range eduprogcomps.Selective[i].CompsInBlock {
			eduprogcomps.Selective[i].CompsInBlock[i2].Code = strconv.Itoa(i2 + 1)
			_, err = s.eduprogcompRepo.Update(eduprogcomps.Selective[i].CompsInBlock[i2])
			if err != nil {
				log.Printf("EduprogcompService: %s", err)
				return err
			}
		}
	}

	return nil
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
		updEduprogcomp, err := s.eduprogcompRepo.Update(edcompById)
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

	err = s.codesRedefine(eduprogcomp.EduprogId)
	if err != nil {
		log.Printf("EduprogcompService: %s", err)
		return err
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
		_, err = s.eduprogcompRepo.Update(eduprogcomps.Mandatory[i])
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
			_, err = s.eduprogcompRepo.Update(components.Selective[i].CompsInBlock[i2])
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
			_, err = s.eduprogcompRepo.Update(components.Selective[i].CompsInBlock[i2])
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
