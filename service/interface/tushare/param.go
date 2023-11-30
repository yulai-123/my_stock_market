package tushare

import "my_stock_market/model/do"

type StockBasicParam struct {
	IsHS       string
	ListStatus string
	Exchange   string
	TSCode     string
	Market     string
	Limit      int64
	Offset     int64
	Name       string
}

type StockBasicResult struct {
	StockList []*do.Stock
}

type DailyParam struct {
	TSCode    []string
	TradeDate string
	StartDate string
	EndDate   string
}

type DailyResult struct {
	StockDailyList []*do.StockDaily
}

type WeeklyParam struct {
	TSCode    []string
	TradeDate string
	StartDate string
	EndDate   string
}

type WeeklyResult struct {
	StockWeeklyList []*do.StockWeekly
}

type MonthlyParam struct {
	TSCode    []string
	TradeDate string
	StartDate string
	EndDate   string
}

type MonthlyResult struct {
	StockMonthlyList []*do.StockMonthly
}

type TradeCalParam struct {
	Exchange  string
	StartDate string
	EndDate   string
	IsOpen    string
}

type SingleTradeCal struct {
	Exchange     string `json:"exchange"`
	CalDate      string `json:"cal_date"`
	IsOpen       int64  `json:"is_open"`
	PreTradeDate string `json:"pretrade_date"`
}

type TradeCalResult struct {
	TradeCalList []*SingleTradeCal
}

type DailyBasicParam struct {
	TSCode    string
	TradeDate string
	StartDate string
	EndDate   string
}

type DailyBasicResult struct {
	DailyBasicList []*do.DailyBasic
}

type CashflowParam struct {
	TSCode string
	Period string
}

type CashflowResult struct {
	CashflowList []*do.Cashflow
}

type BalanceSheetParam struct {
	TSCode string
	Period string
}

type BalanceSheetResult struct {
	BalanceSheetList []*do.BalanceSheet
}

type IncomeParam struct {
	TSCode string
	Period string
}

type IncomeResult struct {
	IncomeList []*do.Income
}

type IndexParam struct {
	TSCode string
}

type IndexResult struct {
	IndexList []*do.Index
}

type IndexDailyParam struct {
	TSCode    string
	TradeDate string
	StartDate string
	EndDate   string
}

type IndexDailyResult struct {
	IndexDailyList []*do.IndexDaily
}

type FundBasicParam struct {
	Market string
	Status string
	Limit  int64
	Offset int64
}

type FundBasicResult struct {
	FundBasicList []*do.FundBasic
}

type FundDailyParam struct {
	TSCode    string
	TradeDate string
	StartDate string
	EndDate   string
}

type FundDailyResult struct {
	FundDailyList []*do.FundDaily
}

type FundAdjParam struct {
	TSCode    string
	StartDate string
	EndDate   string
}

type FundAdj struct {
	TSCode    string  `json:"ts_code"`
	TradeDate string  `json:"trade_date"`
	AdjFactor float64 `json:"adj_factor"`
}

type FundAdjResult struct {
	FundAdjList []*FundAdj
}
