package stock

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"my_stock_market/infra/mysql"
	do2 "my_stock_market/model/do"
	"my_stock_market/model/po"
	"my_stock_market/repo/interface/stock"
	"time"
)

type Impl struct {
	provider *mysql.DBProvider
}

func GetStockDAL(ctx context.Context) stock.DAL {
	return &Impl{
		provider: mysql.GetDBProvider(ctx),
	}
}

func (i *Impl) BatchSaveStock(ctx context.Context, param stock.BatchSaveStockParam) (*stock.BatchSaveStockResult, error) {
	stockPOList := make([]*po.Stock, 0)
	createdAt := int64(time.Now().Unix())
	updatedAt := createdAt
	for _, do := range param.StockList {
		po, err := do2.StockDO2PO(ctx, do)
		if err != nil {
			return nil, err
		}
		po.CreatedAt = createdAt
		po.UpdatedAt = updatedAt
		stockPOList = append(stockPOList, po)
	}

	fieldUpdateColumns := []string{"symbol", "name", "area", "industry", "fullname", "cnspell",
		"market", "exchange", "curr_type", "list_status", "list_date", "delist_date", "is_hs", "updated_at"}
	err := i.provider.WithContext(ctx).Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns(fieldUpdateColumns)}).
		CreateInBatches(stockPOList, 200).Error
	if err != nil {
		logrus.Errorf("[BatchSaveStock] create in batches error: %v", err)
		return nil, err
	}

	return &stock.BatchSaveStockResult{}, nil
}

func (i *Impl) GetAllStock(ctx context.Context, param stock.GetAllStockParam) (*stock.GetAllStockResult, error) {
	result := make([]*do2.Stock, 0)

	stockPOList := make([]*po.Stock, 0)
	err := i.provider.WithContext(ctx).Find(&stockPOList).Error
	if err != nil {
		logrus.Errorf("[GetAllStock] find data error: %v", err)
		return nil, err
	}
	for _, po := range stockPOList {
		result = append(result, do2.StockPO2DO(ctx, po))
	}

	return &stock.GetAllStockResult{StockList: result}, nil
}
