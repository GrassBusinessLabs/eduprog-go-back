package domain

import "time"

type Eduprogscheme struct { //eduprogschemeElement
	Id                 uint64
	SemesterNum        uint16
	Discipline         string
	EduprogId          uint64
	EduprogcompId      uint64
	Eduprogcomp        Eduprogcomp
	CreditsPerSemester uint16
	CreatedDate        time.Time
	UpdatedDate        time.Time
}

func (e Eduprogscheme) GetEduprogschemeId() uint64 {
	return e.Id
}
