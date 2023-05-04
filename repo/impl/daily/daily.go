package daily

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"my_stock_market/infra/mysql"
	do2 "my_stock_market/model/do"
	"my_stock_market/model/po"
	"my_stock_market/repo/interface/daily"
	"time"
)

type Impl struct {
	provider *mysql.DBProvider
}

func NewStockDailyDAL(ctx context.Context) daily.DAL {
	return &Impl{provider: mysql.GetDBProvider(ctx)}
}

func (i *Impl) BatchGetStockDaily(ctx context.Context, param daily.BatchGetStockDailyParam) (*daily.BatchGetStockDailyResult, error) {
	result := make(map[string][]*do2.StockDaily)

	for _, tsCode := range param.TSCode {
		stockDailyList := make([]*po.StockDaily, 0)
		err := i.provider.WithContext(ctx).Where("ts_code = ?", tsCode).Find(&stockDailyList).Error
		if err != nil {
			logrus.Errorf("[BatchGetStockDaily] find data error: %v", err)
			return nil, err
		}
		stockDailyDOList := make([]*do2.StockDaily, 0)
		for _, po := range stockDailyList {
			stockDailyDOList = append(stockDailyDOList, do2.StockDailyPO2DO(ctx, po))
		}
		result[tsCode] = stockDailyDOList
	}

	return &daily.BatchGetStockDailyResult{StockDailyOfTSCodeMap: result}, nil
}

func (i *Impl) BatchSaveStockDaily(ctx context.Context, param daily.BatchSaveStockDailyParam) error {
	var err error
	for t := 0; t < 10; t++ {
		stockDailyPOList := make([]*po.StockDaily, 0)
		createdAt := int64(time.Now().Unix())
		updatedAt := createdAt
		for _, do := range param.StockDailyList {
			po := do2.StockDailyDO2PO(ctx, do)
			po.CreatedAt = createdAt
			po.UpdatedAt = updatedAt
			stockDailyPOList = append(stockDailyPOList, po)
		}
		logrus.Info(len(stockDailyPOList))

		fieldUpdateColumns := []string{"open", "high", "low", "close",
			"pre_close", "change", "pct_chg", "vol", "amount", "updated_at"}
		err = i.provider.WithContext(ctx).Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns(fieldUpdateColumns)}).
			CreateInBatches(stockDailyPOList, 200).Error
		if err != nil {
			logrus.Errorf("[BatchSaveStockDaily] create in batches error: %v, try: %v", err, t)
			time.Sleep(1 * time.Second)
			continue
		}

		return nil
	}
	if err != nil {
		return err
	}

	return nil
}
