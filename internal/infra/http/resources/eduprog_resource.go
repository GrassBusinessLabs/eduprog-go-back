package resources

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type EduprogDto struct {
	Id             uint64 `json:"id"`
	Name           string `json:"name"`
	EducationLevel string `json:"education_level"`
	Stage          string `json:"stage"`
	Speciality     string `json:"speciality"`
	KnowledgeField string `json:"knowledge_field"`
	UserId         uint64 `json:"user_id"`
}

func (d EduprogDto) DomainToDto(eduprog domain.Eduprog) EduprogDto {
	return EduprogDto{
		Id:             eduprog.Id,
		Name:           eduprog.Name,
		EducationLevel: eduprog.EducationLevel,
		Stage:          eduprog.Stage,
		Speciality:     eduprog.Speciality,
		KnowledgeField: eduprog.KnowledgeField,
		UserId:         eduprog.UserId,
	}
}
