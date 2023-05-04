package stock

import "my_stock_market/model/do"

type BatchSaveStockParam struct {
	StockList []*do.Stock
}

type BatchSaveStockResult struct {
}

type GetAllStockParam struct {
}

type GetAllStockResult struct {
	StockList []*do.Stock
}
