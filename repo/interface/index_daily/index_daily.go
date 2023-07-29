package index_daily

import "context"

type DAL interface {
	BatchSaveIndexDaily(ctx context.Context, param BatchSaveIndexDailyParam) error
	BatchGetIndexDaily(ctx context.Context, param BatchGetIndexDailyParam) (*BatchGetIndexDailyResult, error)
}
