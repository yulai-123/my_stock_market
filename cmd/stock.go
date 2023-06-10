package cmd

import (
	"context"
	"my_stock_market/repo/impl/banlance_sheet"
	cashflow2 "my_stock_market/repo/impl/cashflow"
	daily2 "my_stock_market/repo/impl/daily"
	daily_basic2 "my_stock_market/repo/impl/daily_basic"
	income2 "my_stock_market/repo/impl/income"
	monthly2 "my_stock_market/repo/impl/monthly"
	stock2 "my_stock_market/repo/impl/stock"
	weekly2 "my_stock_market/repo/impl/weekly"
	"my_stock_market/repo/interface/balance_sheet"
	"my_stock_market/repo/interface/cashflow"
	"my_stock_market/repo/interface/daily"
	"my_stock_market/repo/interface/daily_basic"
	"my_stock_market/repo/interface/income"
	"my_stock_market/repo/interface/monthly"
	"my_stock_market/repo/interface/stock"
	"my_stock_market/repo/interface/weekly"
	tushare2 "my_stock_market/service/impl/tushare"
	"my_stock_market/service/interface/tushare"
)

type Stock struct {
	TuShare         tushare.TuShare
	StockDAL        stock.DAL
	StockDailyDAL   daily.DAL
	StockWeeklyDAL  weekly.DAL
	StockMonthlyDAL monthly.DAL
	DailyBasicDAL   daily_basic.DAL
	IncomeDAL       income.Income
	CashflowDAL     cashflow.Cashflow
	BalanceSheetDAL balance_sheet.BalanceSheet
}

func NewStock(ctx context.Context) *Stock {
	return &Stock{
		TuShare:         tushare2.NewTuShare(ctx),
		StockDAL:        stock2.GetStockDAL(ctx),
		StockDailyDAL:   daily2.NewStockDailyDAL(ctx),
		StockWeeklyDAL:  weekly2.NewStockWeeklyDAL(ctx),
		StockMonthlyDAL: monthly2.NewStockMonthlyDAL(ctx),
		DailyBasicDAL:   daily_basic2.NewDailyBasicDAL(ctx),
		IncomeDAL:       income2.NewIncomeDAL(ctx),
		CashflowDAL:     cashflow2.NewCashflowDAL(ctx),
		BalanceSheetDAL: banlance_sheet.NewBalanceSheetDAL(ctx),
	}
}
