package eduprog

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const EduprogcompetenciesTableName = "eduprogcompetencies"

type eduprogcompetencies struct {
	Id           uint64 `db:"id,omitempty"`
	CompetencyId uint64 `db:"competency_id"`
	EduprogId    uint64 `db:"eduprog_id"`
	Type         string `db:"type"`
	Code         uint64 `db:"code"`
	Definition   string `db:"definition"`
}

type EduprogcompetenciesRepository interface {
	AddCompetencyToEduprog(eduprogcompetency domain.Eduprogcompetencies) (domain.Eduprogcompetencies, error)
	UpdateCompetency(eduprogcompetency domain.Eduprogcompetencies) (domain.Eduprogcompetencies, error)
	ShowCompetenciesByEduprogId(eduprogId uint64) ([]domain.Eduprogcompetencies, error)
	ShowCompetenciesByType(eduprogId uint64, ttype string) ([]domain.Eduprogcompetencies, error)
	FindById(competencyId uint64) (domain.Eduprogcompetencies, error)
	Delete(competencyId uint64) error
}

type eduprogcompetenciesRepository struct {
	coll db.Collection
}

func NewEduprogcompetenciesRepository(dbSession db.Session) EduprogcompetenciesRepository {
	return eduprogcompetenciesRepository{
		coll: dbSession.Collection(EduprogcompetenciesTableName),
	}
}

func (r eduprogcompetenciesRepository) AddCompetencyToEduprog(eduprogcompetency domain.Eduprogcompetencies) (domain.Eduprogcompetencies, error) {
	ec := r.mapDomainToModel(eduprogcompetency)
	ec.Id = 0
	err := r.coll.InsertReturning(&ec)
	if err != nil {
		return domain.Eduprogcompetencies{}, err
	}

	return r.mapModelToDomain(ec), nil
}

func (r eduprogcompetenciesRepository) UpdateCompetency(eduprogcompetency domain.Eduprogcompetencies) (domain.Eduprogcompetencies, error) {
	e := r.mapDomainToModel(eduprogcompetency)

	err := r.coll.Find(db.Cond{"id": eduprogcompetency.Id}).Update(&e)
	if err != nil {
		return domain.Eduprogcompetencies{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r eduprogcompetenciesRepository) ShowCompetenciesByEduprogId(eduprogId uint64) ([]domain.Eduprogcompetencies, error) {
	var ec []eduprogcompetencies
	err := r.coll.Find(db.Cond{"eduprog_id": eduprogId}).All(&ec)
	if err != nil {
		return []domain.Eduprogcompetencies{}, err
	}

	return r.mapModelToDomainCollection(ec), nil
}

func (r eduprogcompetenciesRepository) ShowCompetenciesByType(eduprogId uint64, ttype string) ([]domain.Eduprogcompetencies, error) {
	var ec []eduprogcompetencies
	if ttype == "ZK" {
		ttype = "ЗК"
	} else if ttype == "FK" {
		ttype = "ФК"
	} else if ttype == "PR" {
		ttype = "ПР"
	} else if ttype == "VFK" {
		ttype = "ВФК"
	} else if ttype == "VPR" {
		ttype = "ВПР"
	}
	err := r.coll.Find(db.Cond{"eduprog_id": eduprogId, "type": ttype}).OrderBy("code").All(&ec)
	if err != nil {
		return []domain.Eduprogcompetencies{}, err
	}

	return r.mapModelToDomainCollection(ec), nil
}

func (r eduprogcompetenciesRepository) FindById(competencyId uint64) (domain.Eduprogcompetencies, error) {
	var es eduprogcompetencies
	err := r.coll.Find(db.Cond{"id": competencyId}).One(&es)
	if err != nil {
		return domain.Eduprogcompetencies{}, err
	}

	return r.mapModelToDomain(es), nil
}

func (r eduprogcompetenciesRepository) Delete(competencyId uint64) error {
	return r.coll.Find(db.Cond{"id": competencyId}).Delete()
}

func (r eduprogcompetenciesRepository) mapDomainToModel(d domain.Eduprogcompetencies) eduprogcompetencies {
	return eduprogcompetencies{
		Id:           d.Id,
		CompetencyId: d.CompetencyId,
		EduprogId:    d.EduprogId,
		Type:         d.Type,
		Code:         d.Code,
		Definition:   d.Definition,
	}
}

func (r eduprogcompetenciesRepository) mapModelToDomain(m eduprogcompetencies) domain.Eduprogcompetencies {
	return domain.Eduprogcompetencies{
		Id:           m.Id,
		CompetencyId: m.CompetencyId,
		EduprogId:    m.EduprogId,
		Type:         m.Type,
		Code:         m.Code,
		Definition:   m.Definition,
	}
}

func (r eduprogcompetenciesRepository) mapModelToDomainCollection(m []eduprogcompetencies) []domain.Eduprogcompetencies {
	result := make([]domain.Eduprogcompetencies, len(m))

	for i := range m {
		result[i] = r.mapModelToDomain(m[i])
	}
	return result
}
