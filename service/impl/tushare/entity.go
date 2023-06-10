package tushare

type APIName = string

var (
	StockBasicAPI   APIName = "stock_basic"
	StockDailyAPI   APIName = "daily"
	StockWeeklyAPI  APIName = "weekly"
	StockMonthlyAPI APIName = "monthly"
	TradeCalAPI     APIName = "trade_cal"
	DailyBasicAPI   APIName = "daily_basic"

	CashflowAPI     APIName = "cashflow"
	BalanceSheetAPI APIName = "balancesheet"
	IncomeAPI       APIName = "income"
)
