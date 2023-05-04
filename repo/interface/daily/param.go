package daily

import "my_stock_market/model/do"

type BatchSaveStockDailyParam struct {
	StockDailyList []*do.StockDaily
}

type BatchGetStockDailyParam struct {
	TSCode []string
}

type BatchGetStockDailyResult struct {
	StockDailyOfTSCodeMap map[string][]*do.StockDaily
}
