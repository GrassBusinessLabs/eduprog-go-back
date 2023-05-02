package domain

import "time"

type Eduprogcomp struct {
	Id          uint64
	Code        string
	Name        string
	Credits     float64
	FreeCredits float64
	ControlType string
	Type        string
	BlockNum    string
	BlockName   string
	Category    string
	EduprogId   uint64
	CreatedDate time.Time
	UpdatedDate time.Time
}

type BlockInfo struct {
	BlockNum     string
	BlockName    string
	CompsInBlock []Eduprogcomp
}

func (e Eduprogcomp) GetEduprogcompId() uint64 {
	return e.Id
}
