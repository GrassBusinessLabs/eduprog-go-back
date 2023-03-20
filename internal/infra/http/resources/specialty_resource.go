package resources

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type SpecialtyDto struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	KFCode         string `json:"kf_code"`
	KnowledgeField string `json:"knowledge_field"`
}

func (d SpecialtyDto) DomainToDto(specialty domain.Specialty) SpecialtyDto {
	return SpecialtyDto{
		Code:           specialty.Code,
		Name:           specialty.Name,
		KFCode:         specialty.KFCode,
		KnowledgeField: specialty.KnowledgeField,
	}
}

func (d SpecialtyDto) DomainToDtoCollection(specialty []domain.Specialty) []SpecialtyDto {
	result := make([]SpecialtyDto, len(specialty))

	for i := range specialty {
		result[i] = d.DomainToDto(specialty[i])
	}

	return result
}
