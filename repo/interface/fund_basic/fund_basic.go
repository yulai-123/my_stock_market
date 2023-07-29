package fund_basic

import "context"

type DAL interface {
	BatchSaveFundBasic(ctx context.Context, param BatchSaveFundBasicParam) (*BatchSaveFundBasicResult, error)
	GetAllFundBasic(ctx context.Context, param GetAllFundBasicParam) (*GetAllFundBasicResult, error)
}
