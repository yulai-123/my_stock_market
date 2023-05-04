package weekly

import "context"

type DAL interface {
	BatchSaveStockWeekly(ctx context.Context, param BatchSaveStockWeeklyParam) error
}
