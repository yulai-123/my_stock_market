package fund_daily

import "context"

type DAL interface {
	BatchSaveFundDaily(ctx context.Context, param BatchSaveFundDailyParam) error
	BatchGetFundDaily(ctx context.Context, param BatchGetFundDailyParam) (*BatchGetFundDailyResult, error)
}
