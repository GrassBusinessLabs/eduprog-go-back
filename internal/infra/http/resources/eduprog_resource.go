package resources

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"time"
)

type EduprogDto struct {
	Id             uint64        `json:"id"`
	Name           string        `json:"name"`
	EducationLevel string        `json:"education_level"`
	Stage          string        `json:"stage"`
	Speciality     string        `json:"speciality"`
	KnowledgeField string        `json:"knowledge_field"`
	UserId         uint64        `json:"user_id"`
	Components     ComponentsDto `json:"components"`
	UpdatedDate    time.Time     `json:"updated_date"`
}

type EduprogWithoutCompsDto struct {
	Id             uint64    `json:"id"`
	Name           string    `json:"name"`
	EducationLevel string    `json:"education_level"`
	Stage          string    `json:"stage"`
	Speciality     string    `json:"speciality"`
	KnowledgeField string    `json:"knowledge_field"`
	UserId         uint64    `json:"user_id"`
	UpdatedDate    time.Time `json:"updated_date"`
}

type ComponentsDto struct {
	Mandatory []EduprogcompDto `json:"mandatory"`
	Selective []EduprogcompDto `json:"selective"`
}

type CreditsDto struct {
	MandatoryCreditsForLevel uint64 `json:"credits_for_level"`
	SelectiveCreditsForLevel uint64 `json:"selective_credits_for_level"`
	TotalCredits             uint64 `json:"total_credits"`
	MandatoryCredits         uint64 `json:"mandatory_credits"`
	SelectiveCredits         uint64 `json:"selective_credits"`
	TotalFreeCredits         uint64 `json:"total_free_credits"`
	MandatoryFreeCredits     uint64 `json:"mandatory_free_credits"`
	SelectiveFreeCredits     uint64 `json:"selective_free_credits"`
}

type OPPLevelStructDto struct {
	Level            string `json:"level"`
	Stage            string `json:"stage"`
	MandatoryCredits uint64 `json:"mandatory_credits"`
	SelectiveCredits uint64 `json:"selective_credits"`
}

type EduprogsDto struct {
	Items []EduprogWithoutCompsDto `json:"items"`
	Total uint64                   `json:"total"`
	Pages uint                     `json:"pages"`
}

func (d EduprogDto) OPPLevelDomainToDto(level domain.OPPLevelStruct) OPPLevelStructDto {
	return OPPLevelStructDto{
		Level:            level.Level,
		Stage:            level.Stage,
		MandatoryCredits: level.MandatoryCredits,
		SelectiveCredits: level.SelectiveCredits,
	}
}

func (d EduprogDto) OPPLevelDomainToDtoCollection(levels []domain.OPPLevelStruct) []OPPLevelStructDto {
	result := make([]OPPLevelStructDto, len(levels))

	for i := range levels {
		result[i] = d.OPPLevelDomainToDto(levels[i])
	}

	return result
}

func (d EduprogDto) DomainToDto(eduprog domain.Eduprog) EduprogWithoutCompsDto {

	return EduprogWithoutCompsDto{
		Id:             eduprog.Id,
		Name:           eduprog.Name,
		EducationLevel: eduprog.EducationLevel,
		Stage:          eduprog.Stage,
		Speciality:     eduprog.Speciality,
		KnowledgeField: eduprog.KnowledgeField,
		UserId:         eduprog.UserId,
		UpdatedDate:    eduprog.UpdatedDate,
	}
}

func (d EduprogDto) DomainToDtoWithComps(eduprog domain.Eduprog, comp domain.Components) EduprogDto {
	var compDto EduprogcompDto
	return EduprogDto{
		Id:             eduprog.Id,
		Name:           eduprog.Name,
		EducationLevel: eduprog.EducationLevel,
		Stage:          eduprog.Stage,
		Speciality:     eduprog.Speciality,
		KnowledgeField: eduprog.KnowledgeField,
		UserId:         eduprog.UserId,
		Components:     compDto.DomainToDtoWCompCollection(comp),
		UpdatedDate:    eduprog.UpdatedDate,
	}
}

func (d EduprogcompDto) DomainToDtoWCompCollection(comps domain.Components) ComponentsDto {
	mandatory := make([]EduprogcompDto, len(comps.Mandatory))
	selective := make([]EduprogcompDto, len(comps.Selective))

	for i := range comps.Mandatory {
		mandatory[i] = d.DomainToDto(comps.Mandatory[i])
	}

	for i := range comps.Selective {
		selective[i] = d.DomainToDto(comps.Selective[i])
	}

	return ComponentsDto{
		Mandatory: mandatory,
		Selective: selective,
	}
}

func (d EduprogDto) DomainToDtoCollection(eduprogs domain.Eduprogs) EduprogsDto {
	result := make([]EduprogWithoutCompsDto, len(eduprogs.Items))

	for i := range eduprogs.Items {
		result[i] = d.DomainToDto(eduprogs.Items[i])
	}

	return EduprogsDto{Items: result, Pages: eduprogs.Pages, Total: eduprogs.Total}
}
