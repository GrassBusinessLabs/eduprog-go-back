package domain

import "time"

type Eduprog struct {
	Id             uint64
	Name           string
	EducationLevel string
	Stage          string
	Speciality     string
	KnowledgeField string
	UserId         uint64
	Components     Components
	CreatedDate    time.Time
	UpdatedDate    time.Time
	DeletedDate    *time.Time
}

type Components struct {
	Mandatory []Eduprogcomp
	Selective []Eduprogcomp
}

type Eduprogs struct {
	Items []Eduprog
	Total uint64
	Pages uint
}

func (e Eduprog) GetEduprogId() uint64 {
	return e.Id
}
