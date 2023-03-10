package eduprog

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const ResultsMatrixTableName = "results_matrix"

type results_martix struct {
	EduprogId       uint64 `db:"eduprog_id"`
	ComponentId     uint64 `db:"component_id"`
	EduprogresultId uint64 `db:"eduprogresult_id"`
}

type ResultsMatrixRepository interface {
	CreateRelation(relation domain.ResultsMatrix) (domain.ResultsMatrix, error)
	ShowByEduprogId(eduprog_id uint64) ([]domain.ResultsMatrix, error)
	Delete(componentId uint64, eduprogresultId uint64) error
}

type resultsMatrixRepository struct {
	coll db.Collection
}

func NewResultsMatrixRepository(dbSession db.Session) ResultsMatrixRepository {
	return resultsMatrixRepository{
		coll: dbSession.Collection(ResultsMatrixTableName),
	}
}

func (r resultsMatrixRepository) CreateRelation(relation domain.ResultsMatrix) (domain.ResultsMatrix, error) {
	cmr := r.mapDomainToModel(relation)
	err := r.coll.InsertReturning(&cmr)
	if err != nil {
		return domain.ResultsMatrix{}, err
	}

	return r.mapModelToDomain(cmr), nil
}

func (r resultsMatrixRepository) ShowByEduprogId(eduprog_id uint64) ([]domain.ResultsMatrix, error) {
	var cmr []results_martix
	err := r.coll.Find(db.Cond{"eduprog_id": eduprog_id}).All(&cmr)
	if err != nil {
		return []domain.ResultsMatrix{}, err
	}
	return r.mapModelToDomainCollection(cmr), nil
}

func (r resultsMatrixRepository) Delete(componentId uint64, eduprogresultId uint64) error {
	return r.coll.Find(db.Cond{"component_id": componentId, "eduprogresult_id": eduprogresultId}).Delete()
}

func (r resultsMatrixRepository) mapDomainToModel(d domain.ResultsMatrix) results_martix {
	return results_martix{
		EduprogId:       d.EduprogId,
		ComponentId:     d.ComponentId,
		EduprogresultId: d.EduprogresultId,
	}
}

func (r resultsMatrixRepository) mapModelToDomain(m results_martix) domain.ResultsMatrix {
	return domain.ResultsMatrix{
		EduprogId:       m.EduprogId,
		ComponentId:     m.ComponentId,
		EduprogresultId: m.EduprogresultId,
	}
}

func (r resultsMatrixRepository) mapModelToDomainCollection(m []results_martix) []domain.ResultsMatrix {
	result := make([]domain.ResultsMatrix, len(m))

	for i := range m {
		result[i] = r.mapModelToDomain(m[i])
	}
	return result
}
