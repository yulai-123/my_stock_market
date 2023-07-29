package do

import (
	"context"
	"my_stock_market/model/po"
)

type Index struct {
	ID       int64  `json:"id"`
	TSCode   string `json:"ts_code"`
	Name     string `json:"name"`
	FullName string `json:"fullname"`
	Market   string `json:"market"`
	ListDate string `json:"list_date"`

	Publisher  string  `json:"publisher"`
	IndexType  string  `json:"index_type"`
	Category   string  `json:"category"`
	BaseDate   string  `json:"base_date"`
	BasePoint  float64 `json:"base_point"`
	WeightRule string  `json:"weight_rule"`
	Desc       string  `json:"desc"`
	ExpDate    string  `json:"exp_date"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	DeletedAt int64 `json:"deleted_at"`
}

func IndexDO2PO(ctx context.Context, do *Index) (*po.Index, error) {
	if do == nil {
		return nil, nil
	}

	return &po.Index{
		ID:         do.ID,
		TSCode:     do.TSCode,
		Name:       do.Name,
		Market:     do.Market,
		ListDate:   do.ListDate,
		Publisher:  do.Publisher,
		IndexType:  do.IndexType,
		Category:   do.Category,
		BaseDate:   do.BaseDate,
		BasePoint:  do.BasePoint,
		WeightRule: do.WeightRule,
		Desc:       do.Desc,
		ExpDate:    do.ExpDate,
		CreatedAt:  do.CreatedAt,
		UpdatedAt:  do.UpdatedAt,
		DeletedAt:  do.DeletedAt,
	}, nil
}

func IndexPO2DO(ctx context.Context, po *po.Index) *Index {
	if po == nil {
		return nil
	}

	return &Index{
		ID:         po.ID,
		TSCode:     po.TSCode,
		Name:       po.Name,
		Market:     po.Market,
		ListDate:   po.ListDate,
		Publisher:  po.Publisher,
		IndexType:  po.IndexType,
		Category:   po.Category,
		BaseDate:   po.BaseDate,
		BasePoint:  po.BasePoint,
		WeightRule: po.WeightRule,
		Desc:       po.Desc,
		ExpDate:    po.ExpDate,
		CreatedAt:  po.CreatedAt,
		UpdatedAt:  po.UpdatedAt,
		DeletedAt:  po.DeletedAt,
	}
}
