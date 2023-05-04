package daily

import "context"

type DAL interface {
	BatchSaveStockDaily(ctx context.Context, param BatchSaveStockDailyParam) error
	BatchGetStockDaily(ctx context.Context, param BatchGetStockDailyParam) (*BatchGetStockDailyResult, error)
}
