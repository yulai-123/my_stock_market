package daily

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"my_stock_market/infra/mysql"
	do2 "my_stock_market/model/do"
	"my_stock_market/model/po"
	"my_stock_market/repo/interface/index_daily"
	"time"
)

type Impl struct {
	provider *mysql.DBProvider
}

func NewIndexDailyDAL(ctx context.Context) index_daily.DAL {
	return &Impl{provider: mysql.GetDBProvider(ctx)}
}

func (i *Impl) BatchGetIndexDaily(ctx context.Context, param index_daily.BatchGetIndexDailyParam) (*index_daily.BatchGetIndexDailyResult, error) {
	result := make(map[string][]*do2.IndexDaily)

	for _, tsCode := range param.TSCode {
		indexDailyList := make([]*po.IndexDaily, 0)

		db := i.provider.WithContext(ctx).Where("ts_code = ?", tsCode)

		if len(param.StartTime) > 0 {
			db = db.Where("trade_date >= ?", param.StartTime)
		}
		if len(param.EndTime) > 0 {
			db = db.Where("trade_date <= ?", param.EndTime)
		}

		err := db.Find(&indexDailyList).Error
		if err != nil {
			logrus.Errorf("[BatchGetIndexDaily] find data error: %v", err)
			return nil, err
		}
		indexDailyDOList := make([]*do2.IndexDaily, 0)
		for _, po := range indexDailyList {
			indexDailyDOList = append(indexDailyDOList, do2.IndexDailyPO2DO(ctx, po))
		}
		result[tsCode] = indexDailyDOList
	}

	return &index_daily.BatchGetIndexDailyResult{IndexDailyOfTSCodeMap: result}, nil
}

func (i *Impl) BatchSaveIndexDaily(ctx context.Context, param index_daily.BatchSaveIndexDailyParam) error {
	var err error
	for t := 0; t < 10; t++ {
		indexDailyPOList := make([]*po.IndexDaily, 0)
		createdAt := int64(time.Now().Unix())
		updatedAt := createdAt
		for _, do := range param.IndexDailyList {
			po := do2.IndexDailyDO2PO(ctx, do)
			po.CreatedAt = createdAt
			po.UpdatedAt = updatedAt
			indexDailyPOList = append(indexDailyPOList, po)
		}
		logrus.Info(len(indexDailyPOList))

		fieldUpdateColumns := []string{"open", "high", "low", "close",
			"pre_close", "change", "pct_chg", "vol", "amount", "updated_at"}
		err = i.provider.WithContext(ctx).Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns(fieldUpdateColumns)}).
			CreateInBatches(indexDailyPOList, 200).Error
		if err != nil {
			logrus.Errorf("[BatchSaveIndexDaily] create in batches error: %v, try: %v", err, t)
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
