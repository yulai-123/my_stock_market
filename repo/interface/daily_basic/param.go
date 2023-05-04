package daily_basic

import "my_stock_market/model/do"

type BatchSaveDailyBasicParam struct {
	DailyBasicList []*do.DailyBasic
}

type BatchGetDailyBasicParam struct {
	TSCodeList []string
	TradeDate  []string
}

type BatchGetDailyBasicResult struct {
	DailyBasicMap map[string][]*do.DailyBasic
}
