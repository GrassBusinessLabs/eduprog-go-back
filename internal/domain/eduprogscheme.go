package domain

import "time"

type Eduprogscheme struct { //eduprogschemeElement
	Id                 uint64
	SemesterNum        uint16
	EduprogId          uint64
	EduprogcompId      uint64
	CreditsPerSemester uint16
	CreatedDate        time.Time
	UpdatedDate        time.Time
}

func (e Eduprogscheme) GetEduprogschemeId() uint64 {
	return e.Id
}
