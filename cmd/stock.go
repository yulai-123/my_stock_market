package cmd

import (
	"context"
	daily2 "my_stock_market/repo/impl/daily"
	daily_basic2 "my_stock_market/repo/impl/daily_basic"
	monthly2 "my_stock_market/repo/impl/monthly"
	stock2 "my_stock_market/repo/impl/stock"
	weekly2 "my_stock_market/repo/impl/weekly"
	"my_stock_market/repo/interface/daily"
	"my_stock_market/repo/interface/daily_basic"
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
}

func NewStock(ctx context.Context) *Stock {
	return &Stock{
		TuShare:         tushare2.NewTuShare(ctx),
		StockDAL:        stock2.GetStockDAL(ctx),
		StockDailyDAL:   daily2.NewStockDailyDAL(ctx),
		StockWeeklyDAL:  weekly2.NewStockWeeklyDAL(ctx),
		StockMonthlyDAL: monthly2.NewStockMonthlyDAL(ctx),
		DailyBasicDAL:   daily_basic2.NewDailyBasicDAL(ctx),
	}
}
