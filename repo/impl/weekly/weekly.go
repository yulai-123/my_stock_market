package weekly

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"my_stock_market/infra/mysql"
	do2 "my_stock_market/model/do"
	"my_stock_market/model/po"
	"my_stock_market/repo/interface/weekly"
	"time"
)

type Impl struct {
	provider *mysql.DBProvider
}

func NewStockWeeklyDAL(ctx context.Context) weekly.DAL {
	return &Impl{provider: mysql.GetDBProvider(ctx)}
}

func (i Impl) BatchSaveStockWeekly(ctx context.Context, param weekly.BatchSaveStockWeeklyParam) error {
	var err error
	for t := 0; t < 10; t++ {
		stockWeeklyPOList := make([]*po.StockWeekly, 0)
		createdAt := int64(time.Now().Unix())
		updatedAt := createdAt
		for _, do := range param.StockWeeklyList {
			po := do2.StockWeeklyDO2PO(ctx, do)
			po.CreatedAt = createdAt
			po.UpdatedAt = updatedAt
			stockWeeklyPOList = append(stockWeeklyPOList, po)
		}
		logrus.Info(len(stockWeeklyPOList))

		fieldUpdateColumns := []string{"open", "high", "low", "close",
			"pre_close", "change", "pct_chg", "vol", "amount", "updated_at"}
		err = i.provider.WithContext(ctx).Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns(fieldUpdateColumns)}).
			CreateInBatches(stockWeeklyPOList, 200).Error
		if err != nil {
			logrus.Errorf("[BatchSaveStockWeekly] create in batches error: %v, try: %v", err, t)
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
