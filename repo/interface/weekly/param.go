package weekly

import "my_stock_market/model/do"

type BatchSaveStockWeeklyParam struct {
	StockWeeklyList []*do.StockWeekly
}
