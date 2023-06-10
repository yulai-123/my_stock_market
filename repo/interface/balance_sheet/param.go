package balance_sheet

import "my_stock_market/model/do"

type BatchSaveBalanceSheetParam struct {
	BalanceSheetList []*do.BalanceSheet
}

type BatchGetBalanceSheetParam struct {
	StockCodeList []string
}

type BatchGetBalanceSheetResult struct {
	BalanceSheetMap map[string][]*do.BalanceSheet
}
