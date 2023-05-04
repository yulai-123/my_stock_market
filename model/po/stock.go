package po

type Stock struct {
	ID         int64  `gorm:"column:id;primaryKey"`
	TSCode     string `gorm:"column:ts_code"`
	Symbol     string `gorm:"column:symbol"`
	Name       string `gorm:"column:name"`
	Area       string `gorm:"column:area"`
	Industry   string `gorm:"column:industry"`
	FullName   string `gorm:"column:fullname"`
	ENName     string `gorm:"column:enname"`
	CNSpell    string `gorm:"column:cnspell"`
	Market     string `gorm:"column:market"`
	Exchange   string `gorm:"column:exchange"`
	CurrType   string `gorm:"column:curr_type"`
	ListStatus string `gorm:"column:list_status"`
	ListDate   string `gorm:"column:list_date"`
	DeListDate string `gorm:"column:delist_date"`
	IsHs       string `gorm:"column:is_hs"`

	CreatedAt int64 `gorm:"column:created_at"`
	UpdatedAt int64 `gorm:"column:updated_at"`
	DeletedAt int64 `gorm:"column:deleted_at"`
}

func (s Stock) TableName() string {
	return "stock"
}
