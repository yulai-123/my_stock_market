package daily_basic

import "context"

type DAL interface {
	BatchSaveDailyBasic(ctx context.Context, param BatchSaveDailyBasicParam) error
	BatchGetDailyBasic(ctx context.Context, param BatchGetDailyBasicParam) (*BatchGetDailyBasicResult, error)
}
