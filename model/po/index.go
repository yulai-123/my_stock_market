package po

type Index struct {
	ID       int64  `gorm:"column:id;primaryKey"`
	TSCode   string `gorm:"column:ts_code"`
	Name     string `gorm:"column:name"`
	Market   string `gorm:"column:market"`
	ListDate string `gorm:"column:list_date"`

	Publisher  string  `gorm:"column:publisher"`
	IndexType  string  `gorm:"column:index_type"`
	Category   string  `gorm:"column:category"`
	BaseDate   string  `gorm:"column:base_date"`
	BasePoint  float64 `gorm:"column:base_point"`
	WeightRule string  `gorm:"column:weight_rule"`
	Desc       string  `gorm:"column:desc"`
	ExpDate    string  `gorm:"column:exp_date"`

	CreatedAt int64 `gorm:"column:created_at"`
	UpdatedAt int64 `gorm:"column:updated_at"`
	DeletedAt int64 `gorm:"column:deleted_at"`
}

func (s Index) TableName() string {
	return "index"
}
