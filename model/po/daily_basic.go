package po

type DailyBasic struct {
	ID int64 `gorm:"column:id;primaryKey"`

	TSCode        string  `gorm:"column:ts_code"`         //TS股票代码
	TradeDate     string  `gorm:"column:trade_date"`      //交易日期
	Close         float64 `gorm:"column:close"`           //当日收盘价
	TurnoverRate  float64 `gorm:"column:turnover_rate"`   //换手率（%）
	TurnoverRateF float64 `gorm:"column:turnover_rate_f"` //换手率（自由流通股）
	VolumeRatio   float64 `gorm:"column:volume_ratio"`    //量比
	PE            float64 `gorm:"column:pe"`              //市盈率（总市值/净利润， 亏损的PE为空）
	PETTM         float64 `gorm:"column:pe_ttm"`          //市盈率（TTM，亏损的PE为空）
	PB            float64 `gorm:"column:pb"`              //市净率（总市值/净资产）
	PS            float64 `gorm:"column:ps"`              //市销率
	PSTTM         float64 `gorm:"column:ps_ttm"`          //市销率（TTM）
	DVRatio       float64 `gorm:"column:dv_ratio"`        //股息率 （%）
	DVTTM         float64 `gorm:"column:dv_ttm"`          //股息率（TTM）（%）
	TotalShare    float64 `gorm:"column:total_share"`     //总股本 （万股）
	FloatShare    float64 `gorm:"column:float_share"`     //流通股本 （万）
	FreeShare     float64 `gorm:"column:free_share"`      //自由流通股本 （万）
	TotalMV       float64 `gorm:"column:total_mv"`        //总市值 （万元）
	CircMV        float64 `gorm:"column:circ_mv"`         //流通市值（万元）

	CreatedAt int64 `gorm:"column:created_at"`
	UpdatedAt int64 `gorm:"column:updated_at"`
	DeletedAt int64 `gorm:"column:deleted_at"`
}

func (d DailyBasic) TableName() string {
	return "daily_basic"
}

/*

名称	类型	描述
ts_code	str	TS股票代码
trade_date	str	交易日期
close	float	当日收盘价
turnover_rate	float	换手率（%）
turnover_rate_f	float	换手率（自由流通股）
volume_ratio	float	量比
pe	float	市盈率（总市值/净利润， 亏损的PE为空）
pe_ttm	float	市盈率（TTM，亏损的PE为空）
pb	float	市净率（总市值/净资产）
ps	float	市销率
ps_ttm	float	市销率（TTM）
dv_ratio	float	股息率 （%）
dv_ttm	float	股息率（TTM）（%）
total_share	float	总股本 （万股）
float_share	float	流通股本 （万股）
free_share	float	自由流通股本 （万）
total_mv	float	总市值 （万元）
circ_mv	float	流通市值（万元）

*/
