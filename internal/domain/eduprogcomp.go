package domain

import "time"

type Eduprogcomp struct {
	Id          uint64
	Code        string
	Name        string
	Credits     float64
	ControlType string
	Type        string
	SubType     string
	Category    string
	EduprogId   uint64
	CreatedDate time.Time
	UpdatedDate time.Time
}

func (e Eduprogcomp) GetEduprogcompId() uint64 {
	return e.Id
}
