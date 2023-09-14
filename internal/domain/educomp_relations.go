package domain

type EducompRelations struct {
	EduprogId   uint64
	BaseCompId  uint64
	ChildCompId uint64
}

type EducompWithPossibleRelations struct {
	Id                uint64
	Code              string
	Name              string
	Credits           float64
	ControlType       string
	Type              string
	BlockNum          string
	BlockName         string
	EduprogId         uint64
	PossibleRelations []Eduprogcomp
}
