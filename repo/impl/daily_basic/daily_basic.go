package daily_basic

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"my_stock_market/infra/mysql"
	do2 "my_stock_market/model/do"
	"my_stock_market/model/po"
	"my_stock_market/repo/interface/daily_basic"
	"time"
)

type Impl struct {
	provider *mysql.DBProvider
}

func NewDailyBasicDAL(ctx context.Context) daily_basic.DAL {
	return &Impl{provider: mysql.GetDBProvider(ctx)}
}

func (i *Impl) BatchGetDailyBasic(ctx context.Context, param daily_basic.BatchGetDailyBasicParam) (*daily_basic.BatchGetDailyBasicResult, error) {
	result := make(map[string][]*do2.DailyBasic)

	for _, tsCode := range param.TSCodeList {
		dailyBasicList := make([]*po.DailyBasic, 0)
		db := i.provider.WithContext(ctx).Where("ts_code = ?", tsCode)
		if len(param.TradeDate) > 0 {
			db = db.Where("trade_date in (?)", param.TradeDate)
		}

		err := db.Find(&dailyBasicList).Error
		if err != nil {
			logrus.Errorf("[BatchGetDailyBasic] find data error: %v", err)
			return nil, err
		}

		dailyBasicDOList := make([]*do2.DailyBasic, 0)
		for _, po := range dailyBasicList {
			dailyBasicDOList = append(dailyBasicDOList, do2.DailyBasicPO2DO(ctx, po))
		}
		result[tsCode] = dailyBasicDOList
	}

	return &daily_basic.BatchGetDailyBasicResult{DailyBasicMap: result}, nil
}

func (i *Impl) BatchSaveDailyBasic(ctx context.Context, param daily_basic.BatchSaveDailyBasicParam) error {
	var err error
	for t := 0; t < 10; t++ {
		dailyBasicPOList := make([]*po.DailyBasic, 0)
		createdAt := int64(time.Now().Unix())
		updatedAt := createdAt
		for _, do := range param.DailyBasicList {
			po := do2.DailyBasicDo2PO(ctx, do)
			po.CreatedAt = createdAt
			po.UpdatedAt = updatedAt
			dailyBasicPOList = append(dailyBasicPOList, po)
		}

		fieldUpdateColumns := []string{"close", "turnover_rate", "turnover_rate_f", "volume_ratio",
			"pe", "pe_ttm", "pb", "ps", "ps_ttm", "dv_ratio", "dv_ttm",
			"total_share", "float_share", "free_share", "total_mv", "circ_mv", "updated_at"}
		err = i.provider.WithContext(ctx).Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns(fieldUpdateColumns)}).
			CreateInBatches(dailyBasicPOList, 200).Error
		if err != nil {
			logrus.Errorf("[BatchSaveDailyBasic] create in batches error: %v, try: %v", err, t)
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
