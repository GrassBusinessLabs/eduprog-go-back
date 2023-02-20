package requests

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
)

type SetComponentToEdprogschemeRequest struct {
	SemesterNum        uint16 `json:"semester_num" validate:"required"`
	Discipline         string `json:"discipline" validate:"required"`
	DisciplineId       uint64 `json:"discipline_id"`
	EduprogId          uint64 `json:"eduprog_id" validate:"required"`
	EduprogcompId      uint64 `json:"eduprogcomp_id" validate:"required"`
	CreditsPerSemester uint16 `json:"credits_per_semester" validate:"required"`
}

type UpdateComponentInEduprogschemeRequest struct {
	SemesterNum        uint16 `json:"semester_num" validate:"required"`
	Discipline         string `json:"discipline" validate:"required"`
	DisciplineId       uint64 `json:"discipline_id"`
	EduprogId          uint64 `json:"eduprog_id" validate:"required"`
	EduprogcompId      uint64 `json:"eduprogcomp_id" validate:"required"`
	CreditsPerSemester uint16 `json:"credits_per_semester" validate:"required"`
}

func (r SetComponentToEdprogschemeRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogscheme{
		SemesterNum:        r.SemesterNum,
		Discipline:         r.Discipline,
		EduprogId:          r.EduprogId,
		EduprogcompId:      r.EduprogcompId,
		CreditsPerSemester: r.CreditsPerSemester,
	}, nil
}

func (r UpdateComponentInEduprogschemeRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogscheme{
		SemesterNum:        r.SemesterNum,
		Discipline:         r.Discipline,
		EduprogId:          r.EduprogId,
		EduprogcompId:      r.EduprogcompId,
		CreditsPerSemester: r.CreditsPerSemester,
	}, nil
}
