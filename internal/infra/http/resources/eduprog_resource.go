package resources

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"time"
)

type EduprogDto struct {
	Id             uint64            `json:"id"`
	Name           string            `json:"name"`
	EducationLevel string            `json:"education_level"`
	Stage          string            `json:"stage"`
	Speciality     string            `json:"speciality"`
	KnowledgeField string            `json:"knowledge_field"`
	UserId         uint64            `json:"user_id"`
	Components     domain.Components `json:"components"`
	UpdatedDate    time.Time         `json:"updated_date"`
}

type ComponentsDto struct {
	Mandatory []domain.Eduprogcomp `json:"mandatory"`
	Selective []domain.Eduprogcomp `json:"selective"`
}

type EduprogsDto struct {
	Items []EduprogDto `json:"items"`
	Total uint64       `json:"total"`
	Pages uint         `json:"pages"`
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
		Components:     eduprog.Components,
		UpdatedDate:    eduprog.UpdatedDate,
	}
}

//func (d EduprogDto) DomainToDtoWCompCollection(eduprog domain.Eduprog, comps domain.Components) EduprogDto {
//
//}

func (d EduprogDto) DomainToDtoCollection(eduprogs domain.Eduprogs) EduprogsDto {
	result := make([]EduprogDto, len(eduprogs.Items))

	for i := range eduprogs.Items {
		result[i] = d.DomainToDto(eduprogs.Items[i])
	}

	return EduprogsDto{Items: result, Pages: eduprogs.Pages, Total: eduprogs.Total}
}
