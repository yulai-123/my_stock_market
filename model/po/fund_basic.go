package po

type FundBasic struct {
	ID            int64   `gorm:"column:id"`
	TSCode        string  `gorm:"column:ts_code"`
	Name          string  `gorm:"column:name"`
	Management    string  `gorm:"column:management"`
	Custodian     string  `gorm:"column:custodian"`
	FundType      string  `gorm:"column:fund_type"`
	FoundDate     string  `gorm:"column:found_date"`
	DueDate       string  `gorm:"column:due_date"`
	ListDate      string  `gorm:"column:list_date"`
	IssueDate     string  `gorm:"column:issue_date"`
	DelistDate    string  `gorm:"column:delist_date"`
	IssueAmount   float64 `gorm:"column:issue_amount"`
	MFee          float64 `gorm:"column:m_fee"`
	CFee          float64 `gorm:"column:c_fee"`
	DurationYear  float64 `gorm:"column:duration_year"`
	PValue        float64 `gorm:"column:p_value"`
	MinAmount     float64 `gorm:"column:min_amount"`
	ExpReturn     float64 `gorm:"column:exp_return"`
	Benchmark     string  `gorm:"column:benchmark"`
	Status        string  `gorm:"column:status"`
	InvestType    string  `gorm:"column:invest_type"`
	Type          string  `gorm:"column:type"`
	Trustee       string  `gorm:"column:trustee"`
	PurcStartdate string  `gorm:"column:purc_startdate"`
	RedmStartdate string  `gorm:"column:redm_startdate"`
	Market        string  `gorm:"column:market"`

	CreatedAt int64 `gorm:"column:created_at"`
	UpdatedAt int64 `gorm:"column:updated_at"`
	DeletedAt int64 `gorm:"column:deleted_at"`
}

func (f FundBasic) TableName() string {
	return "fund_basic"
}
