package eduprog

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const EducompRelationsTableName = "educomp_relations"

type educomp_relations struct {
	EduprogId   uint64 `db:"eduprog_id"`
	BaseCompId  uint64 `db:"base_comp_id"`
	ChildCompId uint64 `db:"child_comp_id"`
}

type EducompRelationsRepository interface {
	CreateRelation(relation domain.Educomp_relations) (domain.Educomp_relations, error)
	ShowByEduprogId(eduprog_id uint64) ([]domain.Educomp_relations, error)
	Delete(base_comp_id uint64, child_comp_id uint64) error
}

type educompRelationsRepository struct {
	coll db.Collection
}

func NewEducompRelationsRepository(dbSession db.Session) EducompRelationsRepository {
	return educompRelationsRepository{
		coll: dbSession.Collection(EducompRelationsTableName),
	}
}

func (r educompRelationsRepository) CreateRelation(relation domain.Educomp_relations) (domain.Educomp_relations, error) {
	er := r.mapDomainToModel(relation)
	err := r.coll.InsertReturning(&er)
	if err != nil {
		return domain.Educomp_relations{}, err
	}

	return r.mapModelToDomain(er), nil
}

func (r educompRelationsRepository) ShowByEduprogId(eduprog_id uint64) ([]domain.Educomp_relations, error) {
	var er []educomp_relations
	err := r.coll.Find(db.Cond{"eduprog_id": eduprog_id}).All(&er)
	if err != nil {
		return []domain.Educomp_relations{}, err
	}
	return r.mapModelToDomainCollection(er), nil
}

func (r educompRelationsRepository) Delete(base_comp_id uint64, child_comp_id uint64) error {
	return r.coll.Find(db.Cond{"base_comp_id": base_comp_id, "child_comp_id": child_comp_id}).Delete()
}

func (r educompRelationsRepository) mapDomainToModel(d domain.Educomp_relations) educomp_relations {
	return educomp_relations{
		EduprogId:   d.EduprogId,
		BaseCompId:  d.BaseCompId,
		ChildCompId: d.ChildCompId,
	}
}

func (r educompRelationsRepository) mapModelToDomain(m educomp_relations) domain.Educomp_relations {
	return domain.Educomp_relations{
		EduprogId:   m.EduprogId,
		BaseCompId:  m.BaseCompId,
		ChildCompId: m.ChildCompId,
	}
}

func (r educompRelationsRepository) mapModelToDomainCollection(m []educomp_relations) []domain.Educomp_relations {
	result := make([]domain.Educomp_relations, len(m))

	for i := range m {
		result[i] = r.mapModelToDomain(m[i])
	}
	return result
}
