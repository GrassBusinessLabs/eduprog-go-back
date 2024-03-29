package domain

import "time"

type Eduprogscheme struct { //eduprogschemeElement
	Id                 uint64
	SemesterNum        uint64
	DisciplineId       uint64
	Row                uint64
	EduprogId          uint64
	EduprogcompId      uint64
	Eduprogcomp        Eduprogcomp
	CreditsPerSemester float64
	CreatedDate        time.Time
	UpdatedDate        time.Time
}

type ExpandEduprogScheme struct {
	ExpandTo           string
	CreditsPerSemester float64
}

func (e Eduprogscheme) GetEduprogschemeId() uint64 {
	return e.Id
}
