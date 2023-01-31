package resources

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
)

type EduprogcompDto struct {
	Id          uint64 `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Credits     uint64 `json:"credits"`
	ControlType string `json:"control_type"`
	Type        string `json:"type"`
	SubType     string `json:"sub_type"`
	Category    string `json:"category"`
	EduprogId   uint64 `json:"eduprog_id"`
}

type EduprogcompsDto struct {
	Items []EduprogcompDto `json:"items"`
	//Total uint64           `json:"total"`
	//Pages uint             `json:"pages"`
}

func (d EduprogcompDto) DomainToDto(eduprogcomp domain.Eduprogcomp) EduprogcompDto {
	return EduprogcompDto{
		Id:          eduprogcomp.Id,
		Code:        eduprogcomp.Code,
		Name:        eduprogcomp.Name,
		Credits:     eduprogcomp.Credits,
		ControlType: eduprogcomp.ControlType,
		Type:        eduprogcomp.Type,
		SubType:     eduprogcomp.SubType,
		Category:    eduprogcomp.Category,
		EduprogId:   eduprogcomp.EduprogId,
	}
}

func (d EduprogcompDto) DomainToDtoCollection(eduprogcomps domain.Eduprogcomps) EduprogcompsDto {
	result := make([]EduprogcompDto, len(eduprogcomps.Items))

	for i := range eduprogcomps.Items {
		result[i] = d.DomainToDto(eduprogcomps.Items[i])
	}

	return EduprogcompsDto{Items: result}
}
