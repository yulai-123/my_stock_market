package income

import "context"

type Income interface {
	BatchSaveIncome(ctx context.Context, param BatchSaveIncomeParam) error
	BatchGetIncome(ctx context.Context, param BatchGetIncomeParam) (*BatchGetIncomeResult, error)
}
