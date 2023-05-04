package stock

import "context"

type DAL interface {
	BatchSaveStock(ctx context.Context, param BatchSaveStockParam) (*BatchSaveStockResult, error)
	GetAllStock(ctx context.Context, param GetAllStockParam) (*GetAllStockResult, error)
}
