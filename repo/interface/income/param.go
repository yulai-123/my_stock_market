package income

import "my_stock_market/model/do"

type BatchSaveIncomeParam struct {
	IncomeList []*do.Income
}

type BatchGetIncomeParam struct {
	StockCodeList []string
}

type BatchGetIncomeResult struct {
	IncomeMap map[string][]*do.Income
}
