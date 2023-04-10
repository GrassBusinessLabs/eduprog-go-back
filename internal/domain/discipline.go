package domain

import "time"

type Discipline struct {
	Id          uint64
	Name        string
	Rows        uint64
	EduprogId   uint64
	CreatedDate time.Time
	UpdatedDate time.Time
}
