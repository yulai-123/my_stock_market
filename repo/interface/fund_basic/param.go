package fund_basic

import "my_stock_market/model/do"

type BatchSaveFundBasicParam struct {
	FundBasicList []*do.FundBasic
}

type BatchSaveFundBasicResult struct {
}

type GetAllFundBasicParam struct {
}

type GetAllFundBasicResult struct {
	FundBasicList []*do.FundBasic
}
