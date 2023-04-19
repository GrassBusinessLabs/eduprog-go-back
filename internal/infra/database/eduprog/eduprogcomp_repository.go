package eduprog

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/upper/db/v4"
	"reflect"
	"sort"
	"strconv"
	"time"
)

const EduprogcompTableName = "eduprogcomp"
const (
	MandCompType   = "ОК"
	SelectCompType = "ВБ"
)

type eduprogcomp struct {
	Id          uint64    `db:"id,omitempty"`
	Code        string    `db:"code"`
	Name        string    `db:"name"`
	Credits     float64   `db:"credits"`
	ControlType string    `db:"control_type"`
	Type        string    `db:"type"`
	BlockNum    string    `db:"block_num"`
	BlockName   string    `db:"block_name"`
	Category    string    `db:"category"`
	EduprogId   uint64    `db:"eduprog_id"`
	CreatedDate time.Time `db:"created_date,omitempty"`
	UpdatedDate time.Time `db:"updated_date,omitempty"`
}

type EduprogcompRepository interface {
	Save(eduprogcomp domain.Eduprogcomp) (domain.Eduprogcomp, error)
	Update(eduprogcomp domain.Eduprogcomp, id uint64) (domain.Eduprogcomp, error)
	ShowList() ([]domain.Eduprogcomp, error)
	FindById(id uint64) (domain.Eduprogcomp, error)
	FindByWODeleteDate(id uint64) (domain.Eduprogcomp, error)
	FindByBlockNum(id uint64, blockNum string) ([]domain.Eduprogcomp, error)
	ShowListByEduprogId(eduprog_id uint64) ([]domain.Eduprogcomp, error)
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
	e.Id = 0
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

func (r eduprogcompRepository) ShowList() ([]domain.Eduprogcomp, error) {
	var eduprogcomps []eduprogcomp

	err := r.coll.Find().All(&eduprogcomps)
	if err != nil {
		return []domain.Eduprogcomp{}, err
	}

	return r.mapModelToDomainCollection(eduprogcomps), nil
}

func (r eduprogcompRepository) ShowListByEduprogId(eduprog_id uint64) ([]domain.Eduprogcomp, error) {
	var eduprogcomps []eduprogcomp

	err := r.coll.Find(db.Cond{"eduprog_id": eduprog_id}).OrderBy("code").All(&eduprogcomps)
	if err != nil {
		return []domain.Eduprogcomp{}, err
	}

	return r.mapModelToDomainCollection(eduprogcomps), nil
}

func (r eduprogcompRepository) FindById(id uint64) (domain.Eduprogcomp, error) {
	var e eduprogcomp
	err := r.coll.Find(db.Cond{"id": id}).One(&e)
	if err != nil {
		return domain.Eduprogcomp{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r eduprogcompRepository) FindByWODeleteDate(id uint64) (domain.Eduprogcomp, error) {
	var e eduprogcomp
	err := r.coll.Find(db.Cond{"id": id}).One(&e)
	if err != nil {
		return domain.Eduprogcomp{}, err
	}

	return r.mapModelToDomain(e), nil
}

func (r eduprogcompRepository) FindByBlockNum(id uint64, blockNum string) ([]domain.Eduprogcomp, error) {
	var e []eduprogcomp
	err := r.coll.Find(db.Cond{"eduprog_id": id, "block_num": blockNum}).All(&e)
	if err != nil {
		return []domain.Eduprogcomp{}, err
	}

	return r.mapModelToDomainCollection(e), nil
}

func (r eduprogcompRepository) SortComponentsByMnS(eduprog_id uint64) (domain.Components, error) {
	var mandeduprogcomp_slice []eduprogcomp
	var seleduprogcomp_slice []eduprogcomp
	var components domain.Components

	err := r.coll.Find(db.Cond{"eduprog_id": eduprog_id, "type": MandCompType}).All(&mandeduprogcomp_slice)
	if err != nil {
		return domain.Components{}, err
	}
	err = r.coll.Find(db.Cond{"eduprog_id": eduprog_id, "type": SelectCompType}).All(&seleduprogcomp_slice)
	if err != nil {
		return domain.Components{}, err
	}

	components.Mandatory = r.mapModelToDomainCollection(mandeduprogcomp_slice)
	components.Selective = r.GetVBBlocksDomain(r.mapModelToDomainCollection(seleduprogcomp_slice))

	return components, err
}

func (r eduprogcompRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id}).Delete()
}

func (r eduprogcompRepository) mapDomainToModel(d domain.Eduprogcomp) eduprogcomp {
	return eduprogcomp{
		Id:          d.Id,
		Code:        d.Code,
		Name:        d.Name,
		Credits:     d.Credits,
		ControlType: d.ControlType,
		Type:        d.Type,
		BlockNum:    d.BlockNum,
		BlockName:   d.BlockName,
		Category:    d.Category,
		EduprogId:   d.EduprogId,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
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
		BlockNum:    m.BlockNum,
		BlockName:   m.BlockName,
		Category:    m.Category,
		EduprogId:   m.EduprogId,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
	}
}

func (r eduprogcompRepository) mapModelToDomainCollection(m []eduprogcomp) []domain.Eduprogcomp {
	result := make([]domain.Eduprogcomp, len(m))

	for i := range m {
		result[i] = r.mapModelToDomain(m[i])
	}

	return result
}

func (r eduprogcompRepository) GetVBBlocksDomain(eduprogcomps []domain.Eduprogcomp) []domain.BlockInfo {
	var blockInfo []domain.BlockInfo
	for i := range eduprogcomps {
		var temp domain.BlockInfo
		temp.BlockNum = eduprogcomps[i].BlockNum
		temp.BlockName = eduprogcomps[i].BlockName
		blockInfo = append(blockInfo, temp)
	}
	blockInfo = RemoveDuplicatesByField(blockInfo, "BlockNum")
	for i := range blockInfo {
		for i2 := range eduprogcomps {
			if blockInfo[i].BlockNum == eduprogcomps[i2].BlockNum {
				blockInfo[i].CompsInBlock = append(blockInfo[i].CompsInBlock, eduprogcomps[i2])
			}
		}
		blockInfo[i].CompsInBlock = sortByCode(blockInfo[i].CompsInBlock)
	}

	sortBlocks(blockInfo)
	for i := range blockInfo {
		for i2, elem := range blockInfo[i].CompsInBlock {
			elem.Code = strconv.Itoa(i2 + 1)
			blockInfo[i].CompsInBlock[i2] = elem
		}

	}
	return blockInfo
}

func sortBlocks(blocks []domain.BlockInfo) {
	sort.Slice(blocks, func(i, j int) bool {
		blockNumI, errI := strconv.Atoi(blocks[i].BlockNum)
		blockNumJ, errJ := strconv.Atoi(blocks[j].BlockNum)
		if errI != nil || errJ != nil {
			// handle error cases where blockNum is not an integer
			return false
		}
		return blockNumI < blockNumJ
	})
}

func sortByCode(eduprogcomps []domain.Eduprogcomp) []domain.Eduprogcomp {
	sort.Slice(eduprogcomps, func(i, j int) bool {
		// Parse the Code field as integers and compare them
		codeI, errI := strconv.ParseUint(eduprogcomps[i].Code, 10, 64)
		codeJ, errJ := strconv.ParseUint(eduprogcomps[j].Code, 10, 64)
		if errI != nil || errJ != nil {
			return eduprogcomps[i].Code < eduprogcomps[j].Code
		}
		return codeI < codeJ
	})
	return eduprogcomps
}

func RemoveDuplicatesByField(mySlice []domain.BlockInfo, fieldName string) []domain.BlockInfo {
	unique := make(map[string]bool)
	result := make([]domain.BlockInfo, 0)
	for _, v := range mySlice {
		fieldValue := reflect.ValueOf(v).FieldByName(fieldName).String()
		if !unique[fieldValue] {
			unique[fieldValue] = true
			result = append(result, v)
		}
	}
	return result
}
