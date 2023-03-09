package eduprog

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const CompetenciesMatrixTableName = "competencies_matrix"

type competencies_martix struct {
	EduprogId    uint64 `db:"eduprog_id"`
	ComponentId  uint64 `db:"component_id"`
	CompetencyId uint64 `db:"competency_id"`
}

type CompetenciesMatrixRepository interface {
	CreateRelation(relation domain.CompetenciesMatrix) (domain.CompetenciesMatrix, error)
	ShowByEduprogId(eduprog_id uint64) ([]domain.CompetenciesMatrix, error)
	Delete(componentId uint64, competencyId uint64) error
}

type competenciesMatrixRepository struct {
	coll db.Collection
}

func NewCompetenciesMatrixRepository(dbSession db.Session) CompetenciesMatrixRepository {
	return competenciesMatrixRepository{
		coll: dbSession.Collection(CompetenciesMatrixTableName),
	}
}

func (r competenciesMatrixRepository) CreateRelation(relation domain.CompetenciesMatrix) (domain.CompetenciesMatrix, error) {
	cmr := r.mapDomainToModel(relation)
	err := r.coll.InsertReturning(&cmr)
	if err != nil {
		return domain.CompetenciesMatrix{}, err
	}

	return r.mapModelToDomain(cmr), nil
}

func (r competenciesMatrixRepository) ShowByEduprogId(eduprog_id uint64) ([]domain.CompetenciesMatrix, error) {
	var cmr []competencies_martix
	err := r.coll.Find(db.Cond{"eduprog_id": eduprog_id}).All(&cmr)
	if err != nil {
		return []domain.CompetenciesMatrix{}, err
	}
	return r.mapModelToDomainCollection(cmr), nil
}

func (r competenciesMatrixRepository) Delete(componentId uint64, competencyId uint64) error {
	return r.coll.Find(db.Cond{"component_id": componentId, "competency_id": competencyId}).Delete()
}

func (r competenciesMatrixRepository) mapDomainToModel(d domain.CompetenciesMatrix) competencies_martix {
	return competencies_martix{
		EduprogId:    d.EduprogId,
		ComponentId:  d.ComponentId,
		CompetencyId: d.CompetencyId,
	}
}

func (r competenciesMatrixRepository) mapModelToDomain(m competencies_martix) domain.CompetenciesMatrix {
	return domain.CompetenciesMatrix{
		EduprogId:    m.EduprogId,
		ComponentId:  m.ComponentId,
		CompetencyId: m.CompetencyId,
	}
}

func (r competenciesMatrixRepository) mapModelToDomainCollection(m []competencies_martix) []domain.CompetenciesMatrix {
	result := make([]domain.CompetenciesMatrix, len(m))

	for i := range m {
		result[i] = r.mapModelToDomain(m[i])
	}
	return result
}
