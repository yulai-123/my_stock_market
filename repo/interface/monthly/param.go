package monthly

import "my_stock_market/model/do"

type BatchSaveStockMonthlyParam struct {
	StockMonthlyList []*do.StockMonthly
}
