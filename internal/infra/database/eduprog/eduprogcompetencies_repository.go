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
	Redefinition string `db:"redefinition"`
}

type EduprogcompetenciesRepository interface {
	AddCompetencyToEduprog(eduprogcompetency domain.Eduprogcompetencies) (domain.Eduprogcompetencies, error)
	UpdateCompetency(eduprogcompetency domain.Eduprogcompetencies, id uint64) (domain.Eduprogcompetencies, error)
	ShowCompetenciesByEduprogId(eduprogId uint64) ([]domain.Eduprogcompetencies, error)
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

	err := r.coll.InsertReturning(&ec)
	if err != nil {
		return domain.Eduprogcompetencies{}, err
	}

	return r.mapModelToDomain(ec), nil
}

func (r eduprogcompetenciesRepository) UpdateCompetency(eduprogcompetency domain.Eduprogcompetencies, id uint64) (domain.Eduprogcompetencies, error) {
	e := r.mapDomainToModel(eduprogcompetency)

	err := r.coll.Find(db.Cond{"id": id}).Update(&e)
	if err != nil {
		return domain.Eduprogcompetencies{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r eduprogcompetenciesRepository) ShowCompetenciesByEduprogId(eduprogId uint64) ([]domain.Eduprogcompetencies, error) {
	var ec []eduprogcompetencies
	err := r.coll.Find(db.Cond{"eduprog_id": eduprogId}).OrderBy("code").All(&ec)
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
		Redefinition: d.Redefinition,
	}
}

func (r eduprogcompetenciesRepository) mapModelToDomain(m eduprogcompetencies) domain.Eduprogcompetencies {
	return domain.Eduprogcompetencies{
		Id:           m.Id,
		CompetencyId: m.CompetencyId,
		EduprogId:    m.EduprogId,
		Type:         m.Type,
		Code:         m.Code,
		Redefinition: m.Redefinition,
	}
}

func (r eduprogcompetenciesRepository) mapModelToDomainCollection(m []eduprogcompetencies) []domain.Eduprogcompetencies {
	result := make([]domain.Eduprogcompetencies, len(m))

	for i := range m {
		result[i] = r.mapModelToDomain(m[i])
	}
	return result
}
