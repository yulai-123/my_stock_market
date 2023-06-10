package do

import (
	"context"
	"my_stock_market/model/po"
)

type Cashflow struct {
	ID        int64 `json:"id"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	DeletedAt int64 `json:"deleted_at"`

	TSCode     string `json:"ts_code"`     // TS代码
	AnnDate    string `json:"ann_date"`    // 公告日期
	FAnnDate   string `json:"f_ann_date"`  // 实际公告日期
	EndDate    string `json:"end_date"`    // 报告期
	ReportType string `json:"report_type"` // 报告类型
	CompType   string `json:"comp_type"`   // 公司类型
	EndType    string `json:"end_type"`    // 报告期

	NCashflowAct float64 `json:"n_cashflow_act"` // 经营活动产生的现金流量净额
}

func CashflowDO2PO(ctx context.Context, do *Cashflow) *po.Cashflow {
	return &po.Cashflow{
		ID:        do.ID,
		CreatedAt: do.CreatedAt,
		UpdatedAt: do.UpdatedAt,
		DeletedAt: do.DeletedAt,

		TSCode:     do.TSCode,
		AnnDate:    do.AnnDate,
		FAnnDate:   do.FAnnDate,
		EndDate:    do.EndDate,
		ReportType: do.ReportType,
		CompType:   do.CompType,
		EndType:    do.EndType,

		NCashflowAct: do.NCashflowAct,
	}
}

func CashflowPO2DO(ctx context.Context, po *po.Cashflow) *Cashflow {
	return &Cashflow{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
		DeletedAt: po.DeletedAt,

		TSCode:     po.TSCode,
		AnnDate:    po.AnnDate,
		FAnnDate:   po.FAnnDate,
		EndDate:    po.EndDate,
		ReportType: po.ReportType,
		CompType:   po.CompType,
		EndType:    po.EndType,

		NCashflowAct: po.NCashflowAct,
	}
}
