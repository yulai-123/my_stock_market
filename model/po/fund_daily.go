package po

type FundDaily struct {
	ID int64 `gorm:"column:id;primaryKey"`

	TSCode    string  `gorm:"column:ts_code"`
	TradeDate string  `gorm:"column:trade_date"`
	Open      float64 `gorm:"column:open"`
	High      float64 `gorm:"column:high"`
	Low       float64 `gorm:"column:low"`
	Close     float64 `gorm:"column:close"`
	PreClose  float64 `gorm:"column:pre_close"`
	Change    float64 `gorm:"column:change"`
	PctChg    float64 `gorm:"column:pct_chg"`
	Vol       float64 `gorm:"column:vol"`
	Amount    float64 `gorm:"column:amount"`

	CreatedAt int64 `gorm:"column:created_at"`
	UpdatedAt int64 `gorm:"column:updated_at"`
	DeletedAt int64 `gorm:"column:deleted_at"`
}

func (s FundDaily) TableName() string {
	return "fund_daily"
}
