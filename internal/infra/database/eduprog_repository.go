package database

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/upper/db/v4"
	"time"
)

const EduprogTableName = "eduprog"

type eduprog struct {
	Id             uint64     `db:"id,omitempty"`
	Name           string     `db:"name"`
	EducationLevel string     `db:"education_level"`
	Stage          string     `db:"stage"`
	Speciality     string     `db:"speciality"`
	KnowledgeField string     `db:"knowledge_field"`
	UserId         uint64     `db:"user_id"`
	CreatedDate    time.Time  `db:"created_date,omitempty"`
	UpdatedDate    time.Time  `db:"updated_date,omitempty"`
	DeletedDate    *time.Time `db:"deleted_date,omitempty"`
}

type EduprogRepository interface {
	Save(eduprog domain.Eduprog) (domain.Eduprog, error)
	Update(eduprog domain.Eduprog, id uint64) (domain.Eduprog, error)
	ShowList() (domain.Eduprogs, error)
	FindById(id uint64) (domain.Eduprog, error)
	Delete(id uint64) error
}

type eduprogRepository struct {
	coll db.Collection
}

func NewEduprogRepository(dbSession db.Session) EduprogRepository {
	return eduprogRepository{
		coll: dbSession.Collection(EduprogTableName),
	}
}

func (r eduprogRepository) Save(eduprog domain.Eduprog) (domain.Eduprog, error) {
	e := r.mapDomainToModel(eduprog)
	e.CreatedDate, e.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&e)
	if err != nil {
		return domain.Eduprog{}, err
	}
	return r.mapModelToDomain(e), nil
}

func (r eduprogRepository) Update(eduprog domain.Eduprog, id uint64) (domain.Eduprog, error) {
	e := r.mapDomainToModel(eduprog)
	e.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": id}).Update(&e)
	if err != nil {
		return domain.Eduprog{}, err
	}
	return r.mapModelToDomain(e), nil
}

func (r eduprogRepository) ShowList() (domain.Eduprogs, error) {
	var eduprog_slice []eduprog
	var eduprogs domain.Eduprogs
	res := r.coll.Find(db.Cond{"deleted_date": nil})
	err := res.All(&eduprog_slice)
	if err != nil {
		return domain.Eduprogs{}, err
	}

	for i := range eduprog_slice {
		eduprogs.Items = append(eduprogs.Items, r.mapModelToDomain(eduprog_slice[i]))
	}
	eduprogs.Total = uint64(len(eduprog_slice))
	return eduprogs, nil
}

func (r eduprogRepository) FindById(id uint64) (domain.Eduprog, error) {
	var e eduprog
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&e)
	if err != nil {
		return domain.Eduprog{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r eduprogRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r eduprogRepository) mapDomainToModel(d domain.Eduprog) eduprog {
	return eduprog{
		Id:             d.Id,
		Name:           d.Name,
		EducationLevel: d.EducationLevel,
		Stage:          d.Stage,
		Speciality:     d.Speciality,
		KnowledgeField: d.KnowledgeField,
		UserId:         d.UserId,
		CreatedDate:    d.CreatedDate,
		UpdatedDate:    d.UpdatedDate,
		DeletedDate:    d.DeletedDate,
	}
}

func (r eduprogRepository) mapModelToDomain(m eduprog) domain.Eduprog {
	return domain.Eduprog{
		Id:             m.Id,
		Name:           m.Name,
		EducationLevel: m.EducationLevel,
		Stage:          m.Stage,
		Speciality:     m.Speciality,
		KnowledgeField: m.KnowledgeField,
		UserId:         m.UserId,
		CreatedDate:    m.CreatedDate,
		UpdatedDate:    m.UpdatedDate,
		DeletedDate:    m.DeletedDate,
	}
}
