package eduprog

import (
	"errors"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/upper/db/v4"
	"time"
)

type OPPLevel string

const (
	EduprogTableName          = "eduprog"
	EntryLevel       OPPLevel = "Початковий рівень (короткий цикл)"
	FirstLevel       OPPLevel = "Перший (бакалаврський) рівень"
	SecondLevel      OPPLevel = "Другий (магістерський) рівень"
	ThirdLevel       OPPLevel = "Третій (освітньо-науковий/освітньо-творчий) рівень"
)

type eduprog struct {
	Id             uint64     `db:"id,omitempty"`
	Name           string     `db:"name"`
	EducationLevel string     `db:"education_level"`
	Stage          string     `db:"stage"`
	SpecialtyCode  string     `db:"speciality_code"`
	Speciality     string     `db:"speciality"`
	KFCode         string     `db:"kf_code"`
	KnowledgeField string     `db:"knowledge_field"`
	UserId         uint64     `db:"user_id"`
	CreatedDate    time.Time  `db:"created_date,omitempty"`
	UpdatedDate    time.Time  `db:"updated_date,omitempty"`
	DeletedDate    *time.Time `db:"deleted_date,omitempty"`
}

type EduprogRepository interface {
	Save(eduprog domain.Eduprog) (domain.Eduprog, error)
	Update(eduprog domain.Eduprog, id uint64) (domain.Eduprog, error)
	ShowList() (domain.Eduprogs, error)
	FindById(id uint64) (domain.Eduprog, error)
	GetOPPLevelsList() ([]domain.OPPLevelStruct, error)
	GetOPPLevelData(level string) (domain.OPPLevelStruct, error)
	Delete(id uint64) error
}

type eduprogRepository struct {
	coll db.Collection
}

func NewEduprogRepository(dbSession db.Session) EduprogRepository {
	return eduprogRepository{
		coll: dbSession.Collection(EduprogTableName),
	}
}

func (r eduprogRepository) Save(eduprog domain.Eduprog) (domain.Eduprog, error) {
	e := r.mapDomainToModel(eduprog)
	e.CreatedDate, e.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&e)
	if err != nil {
		return domain.Eduprog{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r eduprogRepository) Update(eduprog domain.Eduprog, id uint64) (domain.Eduprog, error) {
	e := r.mapDomainToModel(eduprog)
	e.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": id}).Update(&e)
	if err != nil {
		return domain.Eduprog{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r eduprogRepository) ShowList() (domain.Eduprogs, error) {
	var eduprog_slice []eduprog
	var eduprogs domain.Eduprogs

	err := r.coll.Find(db.Cond{"deleted_date": nil}).All(&eduprog_slice)
	if err != nil {
		return domain.Eduprogs{}, err
	}

	for i := range eduprog_slice {
		eduprogs.Items = append(eduprogs.Items, r.mapModelToDomain(eduprog_slice[i]))
	}
	eduprogs.Total = uint64(len(eduprog_slice))

	return eduprogs, nil
}

func (r eduprogRepository) FindById(id uint64) (domain.Eduprog, error) {
	var e eduprog
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&e)
	if err != nil {
		return domain.Eduprog{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r eduprogRepository) GetOPPLevelsList() ([]domain.OPPLevelStruct, error) {
	var levels []domain.OPPLevelStruct
	levelData, err := r.GetOPPLevelData(string(EntryLevel))
	if err != nil {
		return []domain.OPPLevelStruct{}, err
	}
	levels = append(levels, levelData)
	levelData, err = r.GetOPPLevelData(string(FirstLevel))
	if err != nil {
		return []domain.OPPLevelStruct{}, err
	}
	levels = append(levels, levelData)
	levelData, err = r.GetOPPLevelData(string(SecondLevel))
	if err != nil {
		return []domain.OPPLevelStruct{}, err
	}
	levels = append(levels, levelData)
	levelData, err = r.GetOPPLevelData(string(ThirdLevel))
	if err != nil {
		return []domain.OPPLevelStruct{}, err
	}
	levels = append(levels, levelData)

	return levels, nil
}

func (r eduprogRepository) GetOPPLevelData(level string) (domain.OPPLevelStruct, error) {
	var edulevel domain.OPPLevelStruct

	switch level {
	case string(EntryLevel):
		edulevel.Level = string(EntryLevel)
		edulevel.Stage = "Молодший бакалавр"
		edulevel.MandatoryCredits = 90
		edulevel.SelectiveCredits = 30
		return edulevel, nil
	case string(FirstLevel):
		edulevel.Level = string(FirstLevel)
		edulevel.Stage = "Бакалавр"
		edulevel.MandatoryCredits = 180
		edulevel.SelectiveCredits = 60
		return edulevel, nil
	case string(SecondLevel):
		edulevel.Level = string(SecondLevel)
		edulevel.Stage = "Магістр"
		edulevel.MandatoryCredits = 90
		edulevel.SelectiveCredits = 30
		return edulevel, nil
	case string(ThirdLevel):
		edulevel.Level = string(ThirdLevel)
		edulevel.Stage = "Освітньо-науковий/освітньо-творчий рівень"
		edulevel.MandatoryCredits = 30
		edulevel.SelectiveCredits = 30
		return edulevel, nil
	default:
		return edulevel, errors.New("eduprog level error: no such level in enumeration, can use only `Початковий рівень (короткий цикл)`, `Перший (бакалаврський) рівень`, `Другий (магістерський) рівень`, `Третій (освітньо-науковий/освітньо-творчий) рівень`; use method LevelsList")
	}

}

func (r eduprogRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r eduprogRepository) mapDomainToModel(d domain.Eduprog) eduprog {
	return eduprog{
		Id:             d.Id,
		Name:           d.Name,
		EducationLevel: d.EducationLevel,
		Stage:          d.Stage,
		SpecialtyCode:  d.SpecialtyCode,
		Speciality:     d.Speciality,
		KFCode:         d.KFCode,
		KnowledgeField: d.KnowledgeField,
		UserId:         d.UserId,
		CreatedDate:    d.CreatedDate,
		UpdatedDate:    d.UpdatedDate,
		DeletedDate:    d.DeletedDate,
	}
}

func (r eduprogRepository) mapModelToDomain(m eduprog) domain.Eduprog {
	return domain.Eduprog{
		Id:             m.Id,
		Name:           m.Name,
		EducationLevel: m.EducationLevel,
		Stage:          m.Stage,
		SpecialtyCode:  m.SpecialtyCode,
		Speciality:     m.Speciality,
		KFCode:         m.KFCode,
		KnowledgeField: m.KnowledgeField,
		UserId:         m.UserId,
		CreatedDate:    m.CreatedDate,
		UpdatedDate:    m.UpdatedDate,
		DeletedDate:    m.DeletedDate,
	}
}
