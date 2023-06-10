package tushare

import "context"

type TuShare interface {
	StockBasic(ctx context.Context, param StockBasicParam) (*StockBasicResult, error)
	Daily(ctx context.Context, param DailyParam) (*DailyResult, error)
	Weekly(ctx context.Context, param WeeklyParam) (*WeeklyResult, error)
	Monthly(ctx context.Context, param MonthlyParam) (*MonthlyResult, error)
	TradeCal(ctx context.Context, param TradeCalParam) (*TradeCalResult, error)
	DailyBasic(ctx context.Context, param DailyBasicParam) (*DailyBasicResult, error)
	Cashflow(ctx context.Context, param CashflowParam) (*CashflowResult, error)
	BalanceSheet(ctx context.Context, param BalanceSheetParam) (*BalanceSheetResult, error)
	Income(ctx context.Context, param IncomeParam) (*IncomeResult, error)
}
