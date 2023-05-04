package do

import (
	"context"
	"my_stock_market/model/po"
)

type StockMonthly struct {
	ID int64 `json:"id"`

	TSCode    string  `json:"ts_code"`
	TradeDate string  `json:"trade_date"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	PreClose  float64 `json:"pre_close"`
	Change    float64 `json:"change"`
	PctChg    float64 `json:"pct_chg"`
	Vol       float64 `json:"vol"`
	Amount    float64 `json:"amount"`

	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
}

func StockMonthlyDO2PO(ctx context.Context, do *StockMonthly) *po.StockMonthly {
	if do == nil {
		return nil
	}
	return &po.StockMonthly{
		ID:        do.ID,
		TSCode:    do.TSCode,
		TradeDate: do.TradeDate,
		Open:      do.Open,
		High:      do.High,
		Low:       do.Low,
		Close:     do.Close,
		PreClose:  do.PreClose,
		Change:    do.Change,
		PctChg:    do.PctChg,
		Vol:       do.Vol,
		Amount:    do.Amount,
	}
}

func StockMonthlyPO2DO(ctx context.Context, po *StockMonthly) *StockMonthly {
	if po == nil {
		return nil
	}
	return &StockMonthly{
		ID:        po.ID,
		TSCode:    po.TSCode,
		TradeDate: po.TradeDate,
		Open:      po.Open,
		High:      po.High,
		Low:       po.Low,
		Close:     po.Close,
		PreClose:  po.PreClose,
		Change:    po.Change,
		PctChg:    po.PctChg,
		Vol:       po.Vol,
		Amount:    po.Amount,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
		DeletedAt: po.DeletedAt,
	}
}
