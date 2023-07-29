package index_daily

import "my_stock_market/model/do"

type BatchSaveIndexDailyParam struct {
	IndexDailyList []*do.IndexDaily
}

type BatchGetIndexDailyParam struct {
	TSCode    []string
	StartTime string
	EndTime   string
}

type BatchGetIndexDailyResult struct {
	IndexDailyOfTSCodeMap map[string][]*do.IndexDaily
}
