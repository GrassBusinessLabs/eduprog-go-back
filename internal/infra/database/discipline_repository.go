package database

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/upper/db/v4"
	"time"
)

const DisciplineTableName = "discipline"

type discipline struct {
	Id          uint64    `db:"id,omitempty"`
	Name        string    `db:"name"`
	EduprogId   uint64    `db:"eduprog_id"`
	CreatedDate time.Time `db:"created_date,omitempty"`
	UpdatedDate time.Time `db:"updated_date,omitempty"`
}

type DisciplineRepository interface {
	Save(discipline domain.Discipline) (domain.Discipline, error)
	Update(discipline domain.Discipline, id uint64) (domain.Discipline, error)
	ShowDisciplinesByEduprogId(eduprog_id uint64) ([]domain.Discipline, error)
	FindById(id uint64) (domain.Discipline, error)
	Delete(id uint64) error
}

type disciplineRepository struct {
	coll db.Collection
}

func NewDisciplineRepository(dbSession db.Session) DisciplineRepository {
	return disciplineRepository{
		coll: dbSession.Collection(DisciplineTableName),
	}
}

func (r disciplineRepository) Save(discipline domain.Discipline) (domain.Discipline, error) {
	e := r.mapDomainToModel(discipline)
	e.CreatedDate, e.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&e)
	if err != nil {
		return domain.Discipline{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r disciplineRepository) Update(discipline domain.Discipline, id uint64) (domain.Discipline, error) {
	e := r.mapDomainToModel(discipline)
	e.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": id}).Update(&e)
	if err != nil {
		return domain.Discipline{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r disciplineRepository) ShowDisciplinesByEduprogId(eduprog_id uint64) ([]domain.Discipline, error) {
	var d []discipline
	err := r.coll.Find(db.Cond{"eduprog_id": eduprog_id}).OrderBy("-semester_num").All(&d)
	if err != nil {
		return []domain.Discipline{}, err
	}

	return r.mapModelToDomainCollection(d), nil
}

func (r disciplineRepository) FindById(id uint64) (domain.Discipline, error) {
	var e discipline
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&e)
	if err != nil {
		return domain.Discipline{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r disciplineRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id}).Delete()
}

func (r disciplineRepository) mapDomainToModel(d domain.Discipline) discipline {
	return discipline{
		Id:          d.Id,
		Name:        d.Name,
		EduprogId:   d.EduprogId,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
	}
}

func (r disciplineRepository) mapModelToDomain(m discipline) domain.Discipline {
	return domain.Discipline{
		Id:          m.Id,
		Name:        m.Name,
		EduprogId:   m.EduprogId,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
	}
}

func (r disciplineRepository) mapModelToDomainCollection(m []discipline) []domain.Discipline {
	result := make([]domain.Discipline, len(m))

	for i := range m {
		result[i] = r.mapModelToDomain(m[i])
	}

	return result
}
