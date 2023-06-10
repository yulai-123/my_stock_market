package cashflow

import "context"

type Cashflow interface {
	BatchSaveCashflow(ctx context.Context, param BatchSaveCashflowParam) error
	BatchGetCashflow(ctx context.Context, param BatchGetCashflowParam) (*BatchGetCashflowResult, error)
}
