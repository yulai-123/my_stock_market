package do

import (
	"context"
	"my_stock_market/model/po"
)

type Stock struct {
	ID         int64  `json:"id"`
	TSCode     string `json:"ts_code"`
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
	Area       string `json:"area"`
	Industry   string `json:"industry"`
	FullName   string `json:"fullname"`
	ENName     string `json:"enname"`
	CNSpell    string `json:"cnspell"`
	Market     string `json:"market"`
	Exchange   string `json:"exchange"`
	CurrType   string `json:"curr_type"`
	ListStatus string `json:"list_status"`
	ListDate   string `json:"list_date"`
	DeListDate string `json:"de_list_date"`
	IsHs       string `json:"is_hs"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	DeletedAt int64 `json:"deleted_at"`
}

func StockDO2PO(ctx context.Context, do *Stock) (*po.Stock, error) {
	if do == nil {
		return nil, nil
	}

	return &po.Stock{
		ID:         do.ID,
		TSCode:     do.TSCode,
		Symbol:     do.Symbol,
		Name:       do.Name,
		Area:       do.Area,
		Industry:   do.Industry,
		FullName:   do.FullName,
		ENName:     do.ENName,
		CNSpell:    do.CNSpell,
		Market:     do.Market,
		Exchange:   do.Exchange,
		CurrType:   do.CurrType,
		ListStatus: do.ListStatus,
		ListDate:   do.ListDate,
		DeListDate: do.DeListDate,
		IsHs:       do.IsHs,
	}, nil
}

func StockPO2DO(ctx context.Context, po *po.Stock) *Stock {
	if po == nil {
		return nil
	}

	return &Stock{
		ID:         po.ID,
		TSCode:     po.TSCode,
		Symbol:     po.Symbol,
		Name:       po.Name,
		Area:       po.Area,
		Industry:   po.Industry,
		FullName:   po.FullName,
		ENName:     po.ENName,
		CNSpell:    po.CNSpell,
		Market:     po.Market,
		Exchange:   po.Exchange,
		CurrType:   po.CurrType,
		ListStatus: po.ListStatus,
		ListDate:   po.ListDate,
		DeListDate: po.DeListDate,
		IsHs:       po.IsHs,
		CreatedAt:  po.CreatedAt,
		UpdatedAt:  po.UpdatedAt,
		DeletedAt:  po.DeletedAt,
	}
}
