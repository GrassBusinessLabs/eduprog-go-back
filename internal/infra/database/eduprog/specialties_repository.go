package eduprog

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const SpecialtiesTableName = "specialties"

type specialty struct {
	Code           string `db:"code"`
	Name           string `db:"name"`
	KFCode         string `db:"kf_code"`
	KnowledgeField string `db:"knowledge_field"`
}

type SpecialtiesRepository interface {
	CreateSpecialty(specialty domain.Specialty) (domain.Specialty, error)
	UpdateSpecialty(specialty domain.Specialty, code string) (domain.Specialty, error)
	ShowAllSpecialties() ([]domain.Specialty, error)
	ShowByKFCode(kfCode string) ([]domain.Specialty, error)
	FindByCode(code string) (domain.Specialty, error)
	Delete(code string) error
}

type specialtiesRepository struct {
	coll db.Collection
}

func NewSpecialtiesRepository(dbSession db.Session) SpecialtiesRepository {
	return specialtiesRepository{
		coll: dbSession.Collection(SpecialtiesTableName),
	}
}

func (r specialtiesRepository) CreateSpecialty(specialty domain.Specialty) (domain.Specialty, error) {
	s := r.mapDomainToModel(specialty)
	err := r.coll.InsertReturning(&s)
	if err != nil {
		return domain.Specialty{}, err
	}

	return r.mapModelToDomain(s), nil
}

func (r specialtiesRepository) UpdateSpecialty(specialty domain.Specialty, code string) (domain.Specialty, error) {
	s := r.mapDomainToModel(specialty)

	err := r.coll.Find(db.Cond{"code": code}).Update(&s)
	if err != nil {
		return domain.Specialty{}, err
	}

	return r.mapModelToDomain(s), nil
}

func (r specialtiesRepository) ShowAllSpecialties() ([]domain.Specialty, error) {
	var s []specialty
	err := r.coll.Find().OrderBy("code").All(&s)
	if err != nil {
		return []domain.Specialty{}, err
	}

	return r.mapModelToDomainCollection(s), nil
}

func (r specialtiesRepository) ShowByKFCode(kfCode string) ([]domain.Specialty, error) {
	var s []specialty
	err := r.coll.Find(db.Cond{"kf_code": kfCode}).All(&s)
	if err != nil {
		return []domain.Specialty{}, err
	}
	return r.mapModelToDomainCollection(s), nil
}

func (r specialtiesRepository) FindByCode(code string) (domain.Specialty, error) {
	var s specialty
	err := r.coll.Find(db.Cond{"code": code}).One(&s)
	if err != nil {
		return domain.Specialty{}, err
	}

	return r.mapModelToDomain(s), nil
}

func (r specialtiesRepository) Delete(code string) error {
	return r.coll.Find(db.Cond{"code": code}).Delete()
}

func (r specialtiesRepository) mapDomainToModel(d domain.Specialty) specialty {
	return specialty{
		Code:           d.Code,
		Name:           d.Name,
		KFCode:         d.KFCode,
		KnowledgeField: d.KnowledgeField,
	}
}

func (r specialtiesRepository) mapModelToDomain(m specialty) domain.Specialty {
	return domain.Specialty{
		Code:           m.Code,
		Name:           m.Name,
		KFCode:         m.KFCode,
		KnowledgeField: m.KnowledgeField,
	}
}

func (r specialtiesRepository) mapModelToDomainCollection(m []specialty) []domain.Specialty {
	result := make([]domain.Specialty, len(m))

	for i := range m {
		result[i] = r.mapModelToDomain(m[i])
	}
	return result
}
