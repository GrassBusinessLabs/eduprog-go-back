package requests

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
)

type SetComponentToEdprogschemeRequest struct {
	SemesterNum        uint64  `json:"semester_num" validate:"required,number"`
	DisciplineId       uint64  `json:"discipline_id" validate:"required,number"`
	Row                uint64  `json:"row" validate:"number"`
	EduprogId          uint64  `json:"eduprog_id" validate:"required,number"`
	EduprogcompId      uint64  `json:"eduprogcomp_id" validate:"required,number"`
	CreditsPerSemester float64 `json:"credits_per_semester" validate:"required,number"`
}

type UpdateComponentInEduprogschemeRequest struct {
	SemesterNum        uint64  `json:"semester_num" validate:"number"`
	DisciplineId       uint64  `json:"discipline_id" validate:"number"`
	Row                uint64  `json:"row" validate:"number"`
	EduprogId          uint64  `json:"eduprog_id" validate:"number"`
	EduprogcompId      uint64  `json:"eduprogcomp_id" validate:"number"`
	CreditsPerSemester float64 `json:"credits_per_semester" validate:"number"`
}

func (r SetComponentToEdprogschemeRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogscheme{
		SemesterNum:        r.SemesterNum,
		DisciplineId:       r.DisciplineId,
		Row:                r.Row,
		EduprogId:          r.EduprogId,
		EduprogcompId:      r.EduprogcompId,
		CreditsPerSemester: r.CreditsPerSemester,
	}, nil
}

func (r UpdateComponentInEduprogschemeRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogscheme{
		SemesterNum:        r.SemesterNum,
		DisciplineId:       r.DisciplineId,
		Row:                r.Row,
		EduprogId:          r.EduprogId,
		EduprogcompId:      r.EduprogcompId,
		CreditsPerSemester: r.CreditsPerSemester,
	}, nil
}
