package index

import "context"

type DAL interface {
	BatchSaveIndex(ctx context.Context, param BatchSaveIndexParam) (*BatchSaveIndexResult, error)
	GetAllIndex(ctx context.Context, param GetAllIndexParam) (*GetAllIndexResult, error)
}
