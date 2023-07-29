package fund_daily

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"my_stock_market/infra/mysql"
	do2 "my_stock_market/model/do"
	"my_stock_market/model/po"
	"my_stock_market/repo/interface/fund_daily"
	"time"
)

type Impl struct {
	provider *mysql.DBProvider
}

func NewFundDailyDAL(ctx context.Context) fund_daily.DAL {
	return &Impl{provider: mysql.GetDBProvider(ctx)}
}

func (i *Impl) BatchGetFundDaily(ctx context.Context, param fund_daily.BatchGetFundDailyParam) (*fund_daily.BatchGetFundDailyResult, error) {
	result := make(map[string][]*do2.FundDaily)

	for _, tsCode := range param.TSCode {
		fundDailyList := make([]*po.FundDaily, 0)

		db := i.provider.WithContext(ctx).Where("ts_code = ?", tsCode)

		if len(param.StartTime) > 0 {
			db = db.Where("trade_date >= ?", param.StartTime)
		}
		if len(param.EndTime) > 0 {
			db = db.Where("trade_date <= ?", param.EndTime)
		}

		err := db.Find(&fundDailyList).Error
		if err != nil {
			logrus.Errorf("[BatchGetFundDaily] find data error: %v", err)
			return nil, err
		}
		fundDailyDOList := make([]*do2.FundDaily, 0)
		for _, po := range fundDailyList {
			fundDailyDOList = append(fundDailyDOList, do2.FundDailyPO2DO(ctx, po))
		}
		result[tsCode] = fundDailyDOList
	}

	return &fund_daily.BatchGetFundDailyResult{FundDailyOfTSCodeMap: result}, nil
}

func (i *Impl) BatchSaveFundDaily(ctx context.Context, param fund_daily.BatchSaveFundDailyParam) error {
	var err error
	for t := 0; t < 10; t++ {
		fundDailyPOList := make([]*po.FundDaily, 0)
		createdAt := int64(time.Now().Unix())
		updatedAt := createdAt
		for _, do := range param.FundDailyList {
			po := do2.FundDailyDO2PO(ctx, do)
			po.CreatedAt = createdAt
			po.UpdatedAt = updatedAt
			fundDailyPOList = append(fundDailyPOList, po)
		}
		logrus.Info(len(fundDailyPOList))

		fieldUpdateColumns := []string{"open", "high", "low", "close",
			"pre_close", "change", "pct_chg", "vol", "amount", "updated_at"}
		err = i.provider.WithContext(ctx).Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns(fieldUpdateColumns)}).
			CreateInBatches(fundDailyPOList, 200).Error
		if err != nil {
			logrus.Errorf("[BatchSaveFundDaily] create in batches error: %v, try: %v", err, t)
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
