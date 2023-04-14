package domain

import "time"

type Eduprog struct {
	Id             uint64
	Name           string
	EducationLevel string
	Stage          string
	SpecialtyCode  string
	Speciality     string
	KFCode         string
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

type OPPLevelStruct struct {
	Level            string
	Stage            string
	MandatoryCredits float64
	SelectiveCredits float64
}

type Components struct {
	Mandatory []Eduprogcomp
	Selective []BlockInfo
}

func (e Eduprog) GetEduprogId() uint64 {
	return e.Id
}
