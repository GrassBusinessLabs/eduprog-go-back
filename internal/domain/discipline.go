package domain

import "time"

type Discipline struct {
	Id          uint64
	Name        string
	EduprogId   uint64
	CreatedDate time.Time
	UpdatedDate time.Time
}
