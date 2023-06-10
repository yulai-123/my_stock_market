package cashflow

import "my_stock_market/model/do"

type BatchSaveCashflowParam struct {
	CashflowList []*do.Cashflow
}

type BatchGetCashflowParam struct {
	StockCodeList []string
}

type BatchGetCashflowResult struct {
	CashflowMap map[string][]*do.Cashflow
}
