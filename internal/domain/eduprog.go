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
	CreatedDate    time.Time
	UpdatedDate    time.Time
	DeletedDate    *time.Time
}
