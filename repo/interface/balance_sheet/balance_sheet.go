package balance_sheet

import "context"

type BalanceSheet interface {
	BatchSaveBalanceSheet(ctx context.Context, param BatchSaveBalanceSheetParam) error
	BatchGetBalanceSheet(ctx context.Context, param BatchGetBalanceSheetParam) (*BatchGetBalanceSheetResult, error)
}
