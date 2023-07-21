package eduprog

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/upper/db/v4"
	"time"
)

const EduprogschemeTableName = "eduprogscheme"

type eduprogscheme struct {
	Id                 uint64    `db:"id,omitempty"`
	SemesterNum        uint64    `db:"semester_num"`
	DisciplineId       uint64    `db:"discipline_id"`
	Row                uint64    `db:"row"`
	EduprogId          uint64    `db:"eduprog_id"`
	EduprogcompId      uint64    `db:"eduprogcomp_id"`
	CreditsPerSemester float64   `db:"credits_per_semester"`
	CreatedDate        time.Time `db:"created_date,omitempty"`
	UpdatedDate        time.Time `db:"updated_date,omitempty"`
}

type EduprogschemeRepository interface {
	SetComponentToEdprogscheme(eduprogscheme domain.Eduprogscheme) (domain.Eduprogscheme, error)
	UpdateComponentInEduprogscheme(eduprogscheme domain.Eduprogscheme, id uint64) (domain.Eduprogscheme, error)
	FindById(id uint64) (domain.Eduprogscheme, error)
	FindBySemesterNum(semester_num uint16, eduprog_id uint64) ([]domain.Eduprogscheme, error)
	ShowSchemeByEduprogId(eduprog_id uint64) ([]domain.Eduprogscheme, error)
	Delete(id uint64) error
}

type eduprogschemeRepository struct {
	coll db.Collection
}

func NewEduprogschemeRepository(dbSession db.Session) EduprogschemeRepository {
	return eduprogschemeRepository{
		coll: dbSession.Collection(EduprogschemeTableName),
	}
}

func (r eduprogschemeRepository) SetComponentToEdprogscheme(eduprogscheme domain.Eduprogscheme) (domain.Eduprogscheme, error) {
	es := r.mapDomainToModel(eduprogscheme)
	es.Id = 0
	es.CreatedDate, es.UpdatedDate = time.Now(), time.Now()

	err := r.coll.InsertReturning(&es)
	if err != nil {
		return domain.Eduprogscheme{}, err
	}

	return r.mapModelToDomain(es), nil
}

func (r eduprogschemeRepository) UpdateComponentInEduprogscheme(eduprogscheme domain.Eduprogscheme, id uint64) (domain.Eduprogscheme, error) {
	es := r.mapDomainToModel(eduprogscheme)
	es.UpdatedDate = time.Now()

	err := r.coll.Find(db.Cond{"id": id}).Update(&es)
	if err != nil {
		return domain.Eduprogscheme{}, err
	}

	return r.mapModelToDomain(es), nil
}

func (r eduprogschemeRepository) FindById(id uint64) (domain.Eduprogscheme, error) {
	var es eduprogscheme
	err := r.coll.Find(db.Cond{"id": id}).One(&es)
	if err != nil {
		return domain.Eduprogscheme{}, err
	}

	return r.mapModelToDomain(es), nil
}

func (r eduprogschemeRepository) FindBySemesterNum(semester_num uint16, eduprog_id uint64) ([]domain.Eduprogscheme, error) {
	var es []eduprogscheme
	err := r.coll.Find(db.Cond{"eduprog_id": eduprog_id, "semester_num": semester_num}).All(&es)
	if err != nil {
		return []domain.Eduprogscheme{}, err
	}

	return r.mapModelToDomainCollection(es), nil
}

func (r eduprogschemeRepository) ShowSchemeByEduprogId(eduprog_id uint64) ([]domain.Eduprogscheme, error) {
	var es []eduprogscheme
	err := r.coll.Find(db.Cond{"eduprog_id": eduprog_id, "discipline_id >": 0}).OrderBy("semester_num").All(&es)
	if err != nil {
		return []domain.Eduprogscheme{}, err
	}

	return r.mapModelToDomainCollection(es), nil
}

func (r eduprogschemeRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id}).Delete()
}

func (r eduprogschemeRepository) mapDomainToModel(d domain.Eduprogscheme) eduprogscheme {
	return eduprogscheme{
		Id:                 d.Id,
		SemesterNum:        d.SemesterNum,
		DisciplineId:       d.DisciplineId,
		Row:                d.Row,
		EduprogId:          d.EduprogId,
		EduprogcompId:      d.EduprogcompId,
		CreditsPerSemester: d.CreditsPerSemester,
		CreatedDate:        d.CreatedDate,
		UpdatedDate:        d.UpdatedDate,
	}
}

func (r eduprogschemeRepository) mapModelToDomain(m eduprogscheme) domain.Eduprogscheme {
	return domain.Eduprogscheme{
		Id:                 m.Id,
		SemesterNum:        m.SemesterNum,
		DisciplineId:       m.DisciplineId,
		Row:                m.Row,
		EduprogId:          m.EduprogId,
		EduprogcompId:      m.EduprogcompId,
		CreditsPerSemester: m.CreditsPerSemester,
		CreatedDate:        m.CreatedDate,
		UpdatedDate:        m.UpdatedDate,
	}
}

//func (r eduprogschemeRepository) mapDomainToModelCollection(d []domain.Eduprogscheme) []eduprogscheme {
//	result := make([]eduprogscheme, len(d))
//
//	for i := range d {
//		result[i] = r.mapDomainToModel(d[i])
//	}
//
//	return result
//}

func (r eduprogschemeRepository) mapModelToDomainCollection(m []eduprogscheme) []domain.Eduprogscheme {
	result := make([]domain.Eduprogscheme, len(m))

	for i := range m {
		result[i] = r.mapModelToDomain(m[i])
	}

	return result
}
