package resources

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type EduprogschemeDto struct {
	Id                 uint64 `json:"id"`
	SemesterNum        uint16 `json:"semester_num"`
	Discipline         string `json:"discipline"`
	EduprogId          uint64 `json:"eduprog_id"`
	EduprogcompId      uint64 `json:"eduprogcomp_id"`
	CreditsPerSemester uint16 `json:"credits_per_semester"`
}

func (d EduprogschemeDto) DomainToDto(eduprogscheme domain.Eduprogscheme) EduprogschemeDto {
	return EduprogschemeDto{
		Id:                 eduprogscheme.Id,
		SemesterNum:        eduprogscheme.SemesterNum,
		Discipline:         eduprogscheme.Discipline,
		EduprogId:          eduprogscheme.EduprogId,
		EduprogcompId:      eduprogscheme.EduprogcompId,
		CreditsPerSemester: eduprogscheme.CreditsPerSemester,
	}
}

func (d EduprogschemeDto) DomainToDtoCollection(eduprogscheme []domain.Eduprogscheme) []EduprogschemeDto {
	result := make([]EduprogschemeDto, len(eduprogscheme))

	for i := range eduprogscheme {
		result[i] = d.DomainToDto(eduprogscheme[i])
	}

	return result
}
