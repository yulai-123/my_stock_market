package index

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/repo/interface/index_daily"
	"my_stock_market/service/interface/tushare"
)

// SaveAllIndexDaily
// 股市-日线
func (i2 *Index) SaveAllIndexDaily(ctx context.Context) error {
	indexBasicResult, err := i2.TuShare.Index(ctx, tushare.IndexParam{})
	if err != nil {
		return err
	}

	logrus.Infof("拉取指数列表成功，长度：%v", len(indexBasicResult.IndexList))

	for _, index := range indexBasicResult.IndexList {
		err = i2.saveIndexDaily(ctx, index.TSCode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i2 *Index) saveIndexDaily(ctx context.Context, tsCode string) error {
	startDate := "19900101"
	endDate := "20230901"
	indexDailyResult, err := i2.TuShare.IndexDaily(ctx, tushare.IndexDailyParam{
		TSCode:    tsCode,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return err
	}

	err = i2.IndexDailyDAL.BatchSaveIndexDaily(ctx, index_daily.BatchSaveIndexDailyParam{IndexDailyList: indexDailyResult.IndexDailyList})
	if err != nil {
		return err
	}

	logrus.Infof("保存 %v 数据成功, 长度: %v", tsCode, len(indexDailyResult.IndexDailyList))

	return nil
}
