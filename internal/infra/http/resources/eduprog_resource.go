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
	SpecialityCode string        `json:"speciality_code"`
	Speciality     string        `json:"speciality"`
	KFCode         string        `json:"kf_code"`
	KnowledgeField string        `json:"knowledge_field"`
	UserId         uint64        `json:"user_id"`
	Components     ComponentsDto `json:"components"`
	UpdatedDate    time.Time     `json:"updated_date"`
	ApprovalYear   int           `json:"approval_year"`
	ChildOf        uint64        `json:"child_of"`
}

type EduprogWithoutCompsDto struct {
	Id             uint64    `json:"id"`
	Name           string    `json:"name"`
	EducationLevel string    `json:"education_level"`
	Stage          string    `json:"stage"`
	SpecialityCode string    `json:"speciality_code"`
	Speciality     string    `json:"speciality"`
	KFCode         string    `json:"kf_code"`
	KnowledgeField string    `json:"knowledge_field"`
	UserId         uint64    `json:"user_id"`
	UpdatedDate    time.Time `json:"updated_date"`
	ApprovalYear   int       `json:"approval_year"`
	ChildOf        uint64    `json:"child_of"`
}

type ComponentsDto struct {
	Mandatory []EduprogcompDto `json:"mandatory"`
	Selective []BlockInfoDto   `json:"selective"`
}

type CreditsDto struct {
	MandatoryCreditsForLevel float64 `json:"credits_for_level"`
	SelectiveCreditsForLevel float64 `json:"selective_credits_for_level"`
	TotalCredits             float64 `json:"total_credits"`
	MandatoryCredits         float64 `json:"mandatory_credits"`
	SelectiveCredits         float64 `json:"selective_credits"`
	TotalFreeCredits         float64 `json:"total_free_credits"`
	MandatoryFreeCredits     float64 `json:"mandatory_free_credits"`
	SelectiveFreeCredits     float64 `json:"selective_free_credits"`
	MinCreditsForVB          float64 `json:"min_credits_for_vb"`
	MaxCreditsForVB          float64 `json:"max_credits_for_vb"`
}

type OPPLevelStructDto struct {
	Level            string  `json:"level"`
	Stage            string  `json:"stage"`
	MandatoryCredits float64 `json:"mandatory_credits"`
	SelectiveCredits float64 `json:"selective_credits"`
}

type EduprogsDto struct {
	Items []EduprogWithoutCompsDto `json:"items"`
	Total uint64                   `json:"total"`
	Pages uint                     `json:"pages"`
}

func (d EduprogDto) DomainToDto(eduprog domain.Eduprog) EduprogWithoutCompsDto {

	return EduprogWithoutCompsDto{
		Id:             eduprog.Id,
		Name:           eduprog.Name,
		EducationLevel: eduprog.EducationLevel,
		Stage:          eduprog.Stage,
		SpecialityCode: eduprog.SpecialtyCode,
		Speciality:     eduprog.Speciality,
		KFCode:         eduprog.KFCode,
		KnowledgeField: eduprog.KnowledgeField,
		UserId:         eduprog.UserId,
		UpdatedDate:    eduprog.UpdatedDate,
		ApprovalYear:   eduprog.ApprovalYear,
		ChildOf:        eduprog.ChildOf,
	}
}

func (d EduprogDto) DomainToDtoWithComps(eduprog domain.Eduprog, comp domain.Components, selBlocks []domain.BlockInfo) EduprogDto {
	var compDto EduprogcompDto
	return EduprogDto{
		Id:             eduprog.Id,
		Name:           eduprog.Name,
		EducationLevel: eduprog.EducationLevel,
		Stage:          eduprog.Stage,
		SpecialityCode: eduprog.SpecialtyCode,
		Speciality:     eduprog.Speciality,
		KFCode:         eduprog.KFCode,
		KnowledgeField: eduprog.KnowledgeField,
		UserId:         eduprog.UserId,
		Components:     compDto.DomainToDtoWCompCollection(comp, selBlocks),
		UpdatedDate:    eduprog.UpdatedDate,
		ApprovalYear:   eduprog.ApprovalYear,
		ChildOf:        eduprog.ChildOf,
	}
}

func (d EduprogcompDto) DomainToDtoWCompCollection(comps domain.Components, selBlocks []domain.BlockInfo) ComponentsDto {
	mandatory := make([]EduprogcompDto, len(comps.Mandatory))
	selective := make([]BlockInfoDto, len(selBlocks))

	for i := range comps.Mandatory {
		mandatory[i] = d.DomainToDto(comps.Mandatory[i])
	}

	for i := range selBlocks {
		selective[i] = d.BlockInfoToDto(selBlocks[i])
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
