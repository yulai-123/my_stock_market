package po

type Cashflow struct {
	ID int64 `gorm:"column:id;primaryKey"`

	TSCode     string `gorm:"column:ts_code"`
	AnnDate    string `gorm:"column:ann_date"`
	FAnnDate   string `gorm:"column:f_ann_date"`
	EndDate    string `gorm:"column:end_date"`
	ReportType string `gorm:"column:report_type"`
	CompType   string `gorm:"column:comp_type"`
	EndType    string `gorm:"column:end_type"`

	NCashflowAct float64 `gorm:"column:n_cashflow_act"`

	CreatedAt int64 `gorm:"column:created_at"`
	UpdatedAt int64 `gorm:"column:updated_at"`
	DeletedAt int64 `gorm:"column:deleted_at"`
}

func (c Cashflow) TableName() string {
	return "cashflow"
}
