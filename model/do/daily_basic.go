package do

import (
	"context"
	"my_stock_market/model/po"
)

type DailyBasic struct {
	ID int64 `json:"id"`

	TSCode        string  `json:"ts_code"`         //TS股票代码
	TradeDate     string  `json:"trade_date"`      //交易日期
	Close         float64 `json:"close"`           //当日收盘价
	TurnoverRate  float64 `json:"turnover_rate"`   //换手率（%）
	TurnoverRateF float64 `json:"turnover_rate_f"` //换手率（自由流通股）
	VolumeRatio   float64 `json:"volume_ratio"`    //量比
	PE            float64 `json:"pe"`              //市盈率（总市值/净利润， 亏损的PE为空）
	PETTM         float64 `json:"pe_ttm"`          //市盈率（TTM，亏损的PE为空）
	PB            float64 `json:"pb"`              //市净率（总市值/净资产）
	PS            float64 `json:"ps"`              //市销率
	PSTTM         float64 `json:"ps_ttm"`          //市销率（TTM）
	DVRatio       float64 `json:"dv_ratio"`        //股息率 （%）
	DVTTM         float64 `json:"dv_ttm"`          //股息率（TTM）（%）
	TotalShare    float64 `json:"total_share"`     //总股本 （万股）
	FloatShare    float64 `json:"float_share"`     //流通股本 （万）
	FreeShare     float64 `json:"free_share"`      //自由流通股本 （万）
	TotalMV       float64 `json:"total_mv"`        //总市值 （万元）
	CircMV        float64 `json:"circ_mv"`         //流通市值（万元）

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	DeletedAt int64 `json:"deleted_at"`
}

func DailyBasicDo2PO(ctx context.Context, do *DailyBasic) *po.DailyBasic {
	return &po.DailyBasic{
		ID:            do.ID,
		TSCode:        do.TSCode,
		TradeDate:     do.TradeDate,
		Close:         do.Close,
		TurnoverRate:  do.TurnoverRate,
		TurnoverRateF: do.TurnoverRateF,
		VolumeRatio:   do.VolumeRatio,
		PE:            do.PE,
		PETTM:         do.PETTM,
		PB:            do.PB,
		PS:            do.PS,
		PSTTM:         do.PSTTM,
		DVRatio:       do.DVRatio,
		DVTTM:         do.DVTTM,
		TotalShare:    do.TotalShare,
		FloatShare:    do.FloatShare,
		FreeShare:     do.FreeShare,
		TotalMV:       do.TotalMV,
		CircMV:        do.CircMV,
		CreatedAt:     do.CreatedAt,
		UpdatedAt:     do.UpdatedAt,
		DeletedAt:     do.DeletedAt,
	}
}

func DailyBasicPO2DO(ctx context.Context, po *po.DailyBasic) *DailyBasic {
	return &DailyBasic{
		ID:            po.ID,
		TSCode:        po.TSCode,
		TradeDate:     po.TradeDate,
		Close:         po.Close,
		TurnoverRate:  po.TurnoverRate,
		TurnoverRateF: po.TurnoverRateF,
		VolumeRatio:   po.VolumeRatio,
		PE:            po.PE,
		PETTM:         po.PETTM,
		PB:            po.PB,
		PS:            po.PS,
		PSTTM:         po.PSTTM,
		DVRatio:       po.DVRatio,
		DVTTM:         po.DVTTM,
		TotalShare:    po.TotalShare,
		FloatShare:    po.FloatShare,
		FreeShare:     po.FreeShare,
		TotalMV:       po.TotalMV,
		CircMV:        po.CircMV,
		CreatedAt:     po.CreatedAt,
		UpdatedAt:     po.UpdatedAt,
		DeletedAt:     po.DeletedAt,
	}
}
