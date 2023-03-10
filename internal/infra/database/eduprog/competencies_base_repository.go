package eduprog

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const CompetenciesBaseTableName = "competencies_base"

type competencies_base struct {
	Id         uint64 `db:"id"`
	Type       string `db:"type"`
	Definition string `db:"definition"`
	Specialty  string `db:"specialty"`
}

type CompetenciesBaseRepository interface {
	ShowAllCompetencies() ([]domain.CompetenciesBase, error)
	ShowZK() ([]domain.CompetenciesBase, error)
	ShowFK() ([]domain.CompetenciesBase, error)
	FindById(id uint64) (domain.CompetenciesBase, error)
}

type competenciesBaseRepository struct {
	coll db.Collection
}

func NewCompetenciesBaseRepository(dbSession db.Session) CompetenciesBaseRepository {
	return competenciesBaseRepository{
		coll: dbSession.Collection(CompetenciesBaseTableName),
	}
}

func (r competenciesBaseRepository) ShowAllCompetencies() ([]domain.CompetenciesBase, error) {
	var c []competencies_base
	err := r.coll.Find().All(&c)
	if err != nil {
		return []domain.CompetenciesBase{}, err
	}

	return r.mapModelToDomainCollection(c), nil
}

func (r competenciesBaseRepository) ShowZK() ([]domain.CompetenciesBase, error) {
	var c []competencies_base
	err := r.coll.Find(db.Cond{"type": "ЗК"}).All(&c)
	if err != nil {
		return []domain.CompetenciesBase{}, err
	}

	return r.mapModelToDomainCollection(c), nil
}
func (r competenciesBaseRepository) ShowFK() ([]domain.CompetenciesBase, error) {
	var c []competencies_base
	err := r.coll.Find(db.Cond{"type": "ФК"}).All(&c)
	if err != nil {
		return []domain.CompetenciesBase{}, err
	}

	return r.mapModelToDomainCollection(c), nil
}

func (r competenciesBaseRepository) FindById(id uint64) (domain.CompetenciesBase, error) {
	var e competencies_base
	err := r.coll.Find(db.Cond{"id": id}).One(&e)
	if err != nil {
		return domain.CompetenciesBase{}, err
	}

	return r.mapModelToDomain(e), nil
}

// nolint
func (r competenciesBaseRepository) mapDomainToModel(d domain.CompetenciesBase) competencies_base {
	return competencies_base{
		Id:         d.Id,
		Type:       d.Type,
		Definition: d.Definition,
		Specialty:  d.Specialty,
	}
}

func (r competenciesBaseRepository) mapModelToDomain(m competencies_base) domain.CompetenciesBase {
	return domain.CompetenciesBase{
		Id:         m.Id,
		Type:       m.Type,
		Definition: m.Definition,
		Specialty:  m.Specialty,
	}
}

func (r competenciesBaseRepository) mapModelToDomainCollection(m []competencies_base) []domain.CompetenciesBase {
	result := make([]domain.CompetenciesBase, len(m))

	for i := range m {
		result[i] = r.mapModelToDomain(m[i])
	}
	return result
}
