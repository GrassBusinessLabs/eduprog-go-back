package resources

import "github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"

type EduprogschemeDto struct {
	Id                 uint64         `json:"id"`
	SemesterNum        uint64         `json:"semester_num"`
	DisciplineId       uint64         `json:"discipline_id"`
	Row                uint64         `json:"row"`
	EduprogId          uint64         `json:"eduprog_id"`
	EduprogcompId      uint64         `json:"eduprogcomp_id"`
	Eduprogcomp        EduprogcompDto `json:"eduprogcomp"`
	CreditsPerSemester float64        `json:"credits_per_semester"`
}

func (d EduprogschemeDto) DomainToDto(eduprogscheme domain.Eduprogscheme, comp domain.Eduprogcomp) EduprogschemeDto {
	var compDto EduprogcompDto
	return EduprogschemeDto{
		Id:                 eduprogscheme.Id,
		SemesterNum:        eduprogscheme.SemesterNum,
		DisciplineId:       eduprogscheme.DisciplineId,
		Row:                eduprogscheme.Row,
		EduprogId:          eduprogscheme.EduprogId,
		EduprogcompId:      eduprogscheme.EduprogcompId,
		Eduprogcomp:        compDto.DomainToDto(comp),
		CreditsPerSemester: eduprogscheme.CreditsPerSemester,
	}
}

func (d EduprogschemeDto) DomainToDtoCollection(eduprogscheme []domain.Eduprogscheme, educomp []domain.Eduprogcomp) []EduprogschemeDto {
	result := make([]EduprogschemeDto, len(eduprogscheme))

	for i := range eduprogscheme {
		for i2 := range educomp {
			if eduprogscheme[i].EduprogcompId == educomp[i2].Id {
				result[i] = d.DomainToDto(eduprogscheme[i], educomp[i2])
			}
		}

	}

	return result
}
