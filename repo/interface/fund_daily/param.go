package fund_daily

import "my_stock_market/model/do"

type BatchSaveFundDailyParam struct {
	FundDailyList []*do.FundDaily
}

type BatchGetFundDailyParam struct {
	TSCode    []string
	StartTime string
	EndTime   string
}

type BatchGetFundDailyResult struct {
	FundDailyOfTSCodeMap map[string][]*do.FundDaily
}
