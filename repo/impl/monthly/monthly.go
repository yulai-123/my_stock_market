package monthly

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"my_stock_market/infra/mysql"
	do2 "my_stock_market/model/do"
	"my_stock_market/model/po"
	"my_stock_market/repo/interface/monthly"
	"time"
)

type Impl struct {
	provider *mysql.DBProvider
}

func NewStockMonthlyDAL(ctx context.Context) monthly.DAL {
	return &Impl{provider: mysql.GetDBProvider(ctx)}
}

func (i Impl) BatchSaveStockMonthly(ctx context.Context, param monthly.BatchSaveStockMonthlyParam) error {
	var err error
	for t := 0; t < 10; t++ {
		stockMonthlyPOList := make([]*po.StockMonthly, 0)
		createdAt := int64(time.Now().Unix())
		updatedAt := createdAt
		for _, do := range param.StockMonthlyList {
			po := do2.StockMonthlyDO2PO(ctx, do)
			po.CreatedAt = createdAt
			po.UpdatedAt = updatedAt
			stockMonthlyPOList = append(stockMonthlyPOList, po)
		}
		logrus.Info(len(stockMonthlyPOList))

		fieldUpdateColumns := []string{"open", "high", "low", "close",
			"pre_close", "change", "pct_chg", "vol", "amount", "updated_at"}
		err = i.provider.WithContext(ctx).Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns(fieldUpdateColumns)}).
			CreateInBatches(stockMonthlyPOList, 200).Error
		if err != nil {
			logrus.Errorf("[BatchSaveStockMonthly] create in batches error: %v, try: %v", err, t)
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
