package index

import "my_stock_market/model/do"

type BatchSaveIndexParam struct {
	IndexList []*do.Index
}

type BatchSaveIndexResult struct {
}

type GetAllIndexParam struct {
}

type GetAllIndexResult struct {
	IndexList []*do.Index
}
