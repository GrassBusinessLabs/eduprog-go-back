package eduprog

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const EduprogresultsTableName = "eduprogresults"

type eduprogresult struct {
	Id         uint64 `db:"id,omitempty"`
	Type       string `db:"type"`
	Code       uint64 `db:"code"`
	Definition string `db:"definition"`
	EduprogId  uint64 `db:"eduprog_id"`
}

type EduprogresultsRepository interface {
	AddEduprogresultToEduprog(eduprogresult domain.Eduprogresult) (domain.Eduprogresult, error)
	UpdateEduprogresult(eduprogresult domain.Eduprogresult, id uint64) (domain.Eduprogresult, error)
	ShowEduprogResultsByEduprogId(eduprogId uint64) ([]domain.Eduprogresult, error)
	FindById(eduprogresultId uint64) (domain.Eduprogresult, error)
	Delete(eduprogresultId uint64) error
}

type eduprogresultsRepository struct {
	coll db.Collection
}

func NewEduprogresultsRepository(dbSession db.Session) EduprogresultsRepository {
	return eduprogresultsRepository{
		coll: dbSession.Collection(EduprogresultsTableName),
	}
}

func (r eduprogresultsRepository) AddEduprogresultToEduprog(eduprogresult domain.Eduprogresult) (domain.Eduprogresult, error) {
	er := r.mapDomainToModel(eduprogresult)

	err := r.coll.InsertReturning(&er)
	if err != nil {
		return domain.Eduprogresult{}, err
	}

	return r.mapModelToDomain(er), nil
}

func (r eduprogresultsRepository) UpdateEduprogresult(eduprogresult domain.Eduprogresult, id uint64) (domain.Eduprogresult, error) {
	e := r.mapDomainToModel(eduprogresult)

	err := r.coll.Find(db.Cond{"id": id}).Update(&e)
	if err != nil {
		return domain.Eduprogresult{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r eduprogresultsRepository) ShowEduprogResultsByEduprogId(eduprogId uint64) ([]domain.Eduprogresult, error) {
	var er []eduprogresult
	err := r.coll.Find(db.Cond{"eduprog_id": eduprogId}).All(&er)
	if err != nil {
		return []domain.Eduprogresult{}, err
	}

	return r.mapModelToDomainCollection(er), nil
}

func (r eduprogresultsRepository) FindById(eduprogresultId uint64) (domain.Eduprogresult, error) {
	var er eduprogresult
	err := r.coll.Find(db.Cond{"id": eduprogresultId}).One(&er)
	if err != nil {
		return domain.Eduprogresult{}, err
	}

	return r.mapModelToDomain(er), nil
}

func (r eduprogresultsRepository) Delete(eduprogresultId uint64) error {
	return r.coll.Find(db.Cond{"id": eduprogresultId}).Delete()
}

func (r eduprogresultsRepository) mapDomainToModel(d domain.Eduprogresult) eduprogresult {
	return eduprogresult{
		Id:         d.Id,
		EduprogId:  d.EduprogId,
		Type:       d.Type,
		Code:       d.Code,
		Definition: d.Definition,
	}
}

func (r eduprogresultsRepository) mapModelToDomain(m eduprogresult) domain.Eduprogresult {
	return domain.Eduprogresult{
		Id:         m.Id,
		EduprogId:  m.EduprogId,
		Type:       m.Type,
		Code:       m.Code,
		Definition: m.Definition,
	}
}

func (r eduprogresultsRepository) mapModelToDomainCollection(m []eduprogresult) []domain.Eduprogresult {
	result := make([]domain.Eduprogresult, len(m))

	for i := range m {
		result[i] = r.mapModelToDomain(m[i])
	}
	return result
}
