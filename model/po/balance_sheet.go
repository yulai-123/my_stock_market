package po

type BalanceSheet struct {
	ID int64 `gorm:"column:id;primaryKey"`

	TSCode     string `gorm:"column:ts_code"`
	AnnDate    string `gorm:"column:ann_date"`
	FAnnDate   string `gorm:"column:f_ann_date"`
	EndDate    string `gorm:"column:end_date"`
	ReportType string `gorm:"column:report_type"`
	CompType   string `gorm:"column:comp_type"`
	EndType    string `gorm:"column:end_type"`

	TotalAssets float64 `gorm:"column:total_assets"` // 总资产
	TotalLiab   float64 `gorm:"column:total_liab"`   // 总负债

	CreatedAt int64 `gorm:"column:created_at"`
	UpdatedAt int64 `gorm:"column:updated_at"`
	DeletedAt int64 `gorm:"column:deleted_at"`
}

func (b BalanceSheet) TableName() string {
	return "balance_sheet"
}
