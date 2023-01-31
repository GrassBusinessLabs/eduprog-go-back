package database

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/upper/db/v4"
	"time"
)

const EduprogcompTableName = "eduprogcomp"
const (
	MandCompType   = "ОК"
	SelectCompType = "ВБ"
)

type eduprogcomp struct {
	Id          uint64     `db:"id,omitempty"`
	Code        string     `db:"code"`
	Name        string     `db:"name"`
	Credits     uint64     `db:"credits"`
	ControlType string     `db:"control_type"`
	Type        string     `db:"type"`
	SubType     string     `db:"sub_type"`
	Category    string     `db:"category"`
	EduprogId   uint64     `db:"eduprog_id"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type EduprogcompRepository interface {
	Save(eduprogcomp domain.Eduprogcomp) (domain.Eduprogcomp, error)
	Update(eduprogcomp domain.Eduprogcomp, id uint64) (domain.Eduprogcomp, error)
	ShowList() (domain.Eduprogcomps, error)
	FindById(id uint64) (domain.Eduprogcomp, error)
	SortComponentsByMnS(eduprog_id uint64) (domain.Components, error)
	Delete(id uint64) error
}

type eduprogcompRepository struct {
	coll db.Collection
}

func NewEduprogcompRepository(dbSession db.Session) EduprogcompRepository {
	return eduprogcompRepository{
		coll: dbSession.Collection(EduprogcompTableName),
	}
}

func (r eduprogcompRepository) Save(eduprogcomp domain.Eduprogcomp) (domain.Eduprogcomp, error) {
	e := r.mapDomainToModel(eduprogcomp)
	e.CreatedDate, e.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&e)
	if err != nil {
		return domain.Eduprogcomp{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r eduprogcompRepository) Update(eduprogcomp domain.Eduprogcomp, id uint64) (domain.Eduprogcomp, error) {
	e := r.mapDomainToModel(eduprogcomp)
	e.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": id}).Update(&e)
	if err != nil {
		return domain.Eduprogcomp{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r eduprogcompRepository) ShowList() (domain.Eduprogcomps, error) {
	var eduprogcomp_slice []eduprogcomp
	var eduprogcomps domain.Eduprogcomps

	err := r.coll.Find(db.Cond{"deleted_date": nil}).All(&eduprogcomp_slice)
	if err != nil {
		return domain.Eduprogcomps{}, err
	}

	for i := range eduprogcomp_slice {
		eduprogcomps.Items = append(eduprogcomps.Items, r.mapModelToDomain(eduprogcomp_slice[i]))
	}

	//eduprogcomps.Total = uint64(len(eduprogcomp_slice))

	return eduprogcomps, nil
}

func (r eduprogcompRepository) FindById(id uint64) (domain.Eduprogcomp, error) {
	var e eduprogcomp
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&e)
	if err != nil {
		return domain.Eduprogcomp{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r eduprogcompRepository) SortComponentsByMnS(eduprog_id uint64) (domain.Components, error) {
	var eduprogcomp_slice []eduprogcomp
	var components domain.Components

	err := r.coll.Find(db.Cond{"eduprog_id": eduprog_id, "deleted_date": nil}).All(&eduprogcomp_slice)
	if err != nil {
		return domain.Components{}, err
	}

	for i := range eduprogcomp_slice {
		if eduprogcomp_slice[i].Type == MandCompType {
			components.Mandatory = append(components.Mandatory, r.mapModelToDomain(eduprogcomp_slice[i]))
		} else if eduprogcomp_slice[i].Type == SelectCompType {
			components.Selective = append(components.Selective, r.mapModelToDomain(eduprogcomp_slice[i]))
		}
	}

	return components, err
}

func (r eduprogcompRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r eduprogcompRepository) mapDomainToModel(d domain.Eduprogcomp) eduprogcomp {
	return eduprogcomp{
		Id:          d.Id,
		Code:        d.Code,
		Name:        d.Name,
		Credits:     d.Credits,
		ControlType: d.ControlType,
		Type:        d.Type,
		SubType:     d.SubType,
		Category:    d.Category,
		EduprogId:   d.EduprogId,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r eduprogcompRepository) mapModelToDomain(m eduprogcomp) domain.Eduprogcomp {
	return domain.Eduprogcomp{
		Id:          m.Id,
		Code:        m.Code,
		Name:        m.Name,
		Credits:     m.Credits,
		ControlType: m.ControlType,
		Type:        m.Type,
		SubType:     m.SubType,
		Category:    m.SubType,
		EduprogId:   m.EduprogId,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}
