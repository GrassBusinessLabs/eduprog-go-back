package resources

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
)

type EduprogcompDto struct {
	Id          uint64  `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Credits     float64 `json:"credits"`
	ControlType string  `json:"control_type"`
	Type        string  `json:"type"`
	BlockNum    string  `json:"block_num"`
	BlockName   string  `json:"block_name"`
	Category    string  `json:"category"`
	EduprogId   uint64  `json:"eduprog_id"`
}

type BlockInfoDto struct {
	BlockNum     string           `json:"block_num"`
	BlockName    string           `json:"block_name"`
	CompsInBlock []EduprogcompDto `json:"comps_in_block"`
}

func (d EduprogcompDto) BlockInfoToDto(blockInfo domain.BlockInfo) BlockInfoDto {
	var compDto EduprogcompDto
	return BlockInfoDto{
		BlockNum:     blockInfo.BlockNum,
		BlockName:    blockInfo.BlockName,
		CompsInBlock: compDto.DomainToDtoCollection(blockInfo.CompsInBlock),
	}
}

func (d EduprogcompDto) BlockInfoToDtoCollection(blockinfo []domain.BlockInfo) []BlockInfoDto {
	result := make([]BlockInfoDto, len(blockinfo))

	for i := range blockinfo {
		result[i] = d.BlockInfoToDto(blockinfo[i])
	}

	return result
}

func (d EduprogcompDto) DomainToDto(eduprogcomp domain.Eduprogcomp) EduprogcompDto {
	return EduprogcompDto{
		Id:          eduprogcomp.Id,
		Code:        eduprogcomp.Code,
		Name:        eduprogcomp.Name,
		Credits:     eduprogcomp.Credits,
		ControlType: eduprogcomp.ControlType,
		Type:        eduprogcomp.Type,
		BlockNum:    eduprogcomp.BlockNum,
		BlockName:   eduprogcomp.BlockName,
		Category:    eduprogcomp.Category,
		EduprogId:   eduprogcomp.EduprogId,
	}
}

func (d EduprogcompDto) DomainToDtoCollection(eduprogcomps []domain.Eduprogcomp) []EduprogcompDto {
	result := make([]EduprogcompDto, len(eduprogcomps))

	for i := range eduprogcomps {
		result[i] = d.DomainToDto(eduprogcomps[i])
	}

	return result
}
