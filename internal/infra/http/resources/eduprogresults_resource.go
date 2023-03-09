package resources

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type EduprogresultsDto struct {
	Id         uint64 `json:"id"`
	Type       string `json:"type"`
	Code       uint64 `json:"code"`
	Definition string `json:"definition"`
	EduprogId  uint64 `json:"eduprog_id"`
}

func (d EduprogresultsDto) DomainToDto(eduprogresult domain.Eduprogresult) EduprogresultsDto {
	return EduprogresultsDto{
		Id:         eduprogresult.Id,
		Type:       eduprogresult.Type,
		Code:       eduprogresult.Code,
		Definition: eduprogresult.Definition,
		EduprogId:  eduprogresult.EduprogId,
	}
}

func (d EduprogresultsDto) DomainToDtoCollection(eduprogresults []domain.Eduprogresult) []EduprogresultsDto {
	result := make([]EduprogresultsDto, len(eduprogresults))

	for i := range eduprogresults {
		result[i] = d.DomainToDto(eduprogresults[i])
	}

	return result
}
