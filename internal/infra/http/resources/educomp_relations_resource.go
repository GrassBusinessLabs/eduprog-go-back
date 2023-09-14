package resources

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type EducompRelationsDto struct {
	EduprogId   uint64 `json:"eduprog_id"`
	BaseCompId  uint64 `json:"base_comp_id"`
	ChildCompId uint64 `json:"child_comp_id"`
}

type EducompWithPossibleRelationsDto struct {
	Id                uint64  `json:"id"`
	Code              string  `json:"code"`
	Name              string  `json:"name"`
	Credits           float64 `json:"credits"`
	ControlType       string  `json:"controlType"`
	Type              string  `json:"type"`
	BlockNum          string  `json:"blockNum"`
	BlockName         string  `json:"blockName"`
	EduprogId         uint64  `json:"eduprogId"`
	PossibleRelations []EduprogcompDto
}

func (d EducompWithPossibleRelationsDto) DomainToDto(relation domain.EducompWithPossibleRelations) EducompWithPossibleRelationsDto {
	var eduprogcompDto EduprogcompDto
	return EducompWithPossibleRelationsDto{
		Id:                relation.Id,
		Code:              relation.Code,
		Name:              relation.Name,
		Credits:           relation.Credits,
		ControlType:       relation.ControlType,
		Type:              relation.Type,
		BlockNum:          relation.BlockNum,
		BlockName:         relation.BlockName,
		EduprogId:         relation.EduprogId,
		PossibleRelations: eduprogcompDto.DomainToDtoCollection(relation.PossibleRelations),
	}
}

func (d EducompWithPossibleRelationsDto) DomainToDtoCollection(relations []domain.EducompWithPossibleRelations) []EducompWithPossibleRelationsDto {
	result := make([]EducompWithPossibleRelationsDto, len(relations))

	for i := range relations {
		result[i] = d.DomainToDto(relations[i])
	}

	return result
}

func (d EducompRelationsDto) DomainToDto(relation domain.EducompRelations) EducompRelationsDto {
	return EducompRelationsDto{
		EduprogId:   relation.EduprogId,
		BaseCompId:  relation.BaseCompId,
		ChildCompId: relation.ChildCompId,
	}
}

func (d EducompRelationsDto) DomainToDtoCollection(relation []domain.EducompRelations) []EducompRelationsDto {
	result := make([]EducompRelationsDto, len(relation))

	for i := range relation {
		result[i] = d.DomainToDto(relation[i])
	}

	return result
}
