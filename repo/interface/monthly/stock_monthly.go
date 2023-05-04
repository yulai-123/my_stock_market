package monthly

import "context"

type DAL interface {
	BatchSaveStockMonthly(ctx context.Context, param BatchSaveStockMonthlyParam) error
}
