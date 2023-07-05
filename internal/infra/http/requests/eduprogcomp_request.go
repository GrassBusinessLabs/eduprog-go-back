package requests

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
)

type CreateEduprogcompRequest struct {
	Name        string  `json:"name" validate:"required,gte=1,max=50"`
	Credits     float64 `json:"credits" validate:"required,number,gt=0"`
	ControlType string  `json:"control_type" validate:"required,gte=1,max=50"`
	Type        string  `json:"type" validate:"required,gte=1,max=50"`
	BlockNum    string  `json:"block_num"`
	BlockName   string  `json:"block_name"`
	Category    string  `json:"category"`
	EduprogId   uint64  `json:"eduprog_id" validate:"required,number"`
}

type UpdateEduprogcompRequest struct {
	Name        string  `json:"name" validate:"gte=1,max=50"`
	Credits     float64 `json:"credits" validate:"number,gt=0"`
	ControlType string  `json:"control_type" validate:"gte=1,max=50"`
	Type        string  `json:"type" validate:"gte=1,max=50"`
	BlockNum    string  `json:"block_num"`
	BlockName   string  `json:"block_name"`
	Category    string  `json:"category"`
	EduprogId   uint64  `json:"eduprog_id" validate:"number"`
}

type UpdateBlockName struct {
	BlockNum  string `json:"block_num"`
	BlockName string `json:"block_name"`
}

func (r CreateEduprogcompRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogcomp{

		Name:        r.Name,
		Credits:     r.Credits,
		ControlType: r.ControlType,
		Type:        r.Type,
		BlockNum:    r.BlockNum,
		BlockName:   r.BlockName,
		Category:    r.Category,
		EduprogId:   r.EduprogId,
	}, nil
}

func (r UpdateEduprogcompRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogcomp{
		Name:        r.Name,
		Credits:     r.Credits,
		ControlType: r.ControlType,
		Type:        r.Type,
		BlockNum:    r.BlockNum,
		BlockName:   r.BlockName,
		Category:    r.Category,
		EduprogId:   r.EduprogId,
	}, nil
}

func (r UpdateBlockName) ToDomainModel() (interface{}, error) {
	return domain.Eduprogcomp{
		BlockNum:  r.BlockNum,
		BlockName: r.BlockName,
	}, nil
}

type SendEduprogcompRequest struct {
	Id uint64 `json:"id"`
}

type SendEduprogcompSliceRequest struct {
	Eduprogcomps []SendEduprogcompRequest `json:"eduprogcomps" validate:"required,dive"`
}

func (r SendEduprogcompRequest) ToDomainModel() (interface{}, error) {
	return domain.Eduprogcomp{
		Id: r.Id,
	}, nil
}

func (r SendEduprogcompSliceRequest) ToDomainModel() (interface{}, error) {
	var eduprogcomps []domain.Eduprogcomp
	for _, eduprogcompReq := range r.Eduprogcomps {
		eduprogcomp, err := eduprogcompReq.ToDomainModel()
		if err != nil {
			return nil, err
		}
		eduprogcomps = append(eduprogcomps, eduprogcomp.(domain.Eduprogcomp))
	}
	return eduprogcomps, nil
}
