package resources

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type EduprogschemeDto struct {
	Id                 uint64         `json:"id"`
	SemesterNum        uint16         `json:"semester_num"`
	DisciplineId       uint64         `json:"discipline_id"`
	EduprogId          uint64         `json:"eduprog_id"`
	EduprogcompId      uint64         `json:"eduprogcomp_id"`
	Eduprogcomp        EduprogcompDto `json:"eduprogcomp"`
	CreditsPerSemester uint16         `json:"credits_per_semester"`
}

func (d EduprogschemeDto) DomainToDto(eduprogscheme domain.Eduprogscheme, comp domain.Eduprogcomp) EduprogschemeDto {
	var compDto EduprogcompDto
	return EduprogschemeDto{
		Id:                 eduprogscheme.Id,
		SemesterNum:        eduprogscheme.SemesterNum,
		DisciplineId:       eduprogscheme.DisciplineId,
		EduprogId:          eduprogscheme.EduprogId,
		EduprogcompId:      eduprogscheme.EduprogcompId,
		Eduprogcomp:        compDto.DomainToDto(comp),
		CreditsPerSemester: eduprogscheme.CreditsPerSemester,
	}
}

func (d EduprogschemeDto) DomainToDtoCollection(eduprogscheme []domain.Eduprogscheme, educomp domain.Eduprogcomps) []EduprogschemeDto {
	result := make([]EduprogschemeDto, len(eduprogscheme))

	for i := range eduprogscheme {
		for i2 := range educomp.Items {
			if eduprogscheme[i].EduprogcompId == educomp.Items[i2].Id {
				result[i] = d.DomainToDto(eduprogscheme[i], educomp.Items[i2])
			}
		}

	}

	return result
}
