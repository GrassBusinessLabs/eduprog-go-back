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

type Eduprogs struct {
	Items []Eduprog
	Total uint64
	Pages uint
}

type Components struct {
	Mandatory []Eduprogcomp
	Selective []Eduprogcomp
}

func (e Eduprog) GetEduprogId() uint64 {
	return e.Id
}
