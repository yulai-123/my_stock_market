package tushare

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/repo/interface/index_daily"
	"my_stock_market/service/interface/tushare"
)

// SaveAllIndexDaily 保存所有工作日&所有股票的日线数据
// 先拉取所有ts_code
// 然后根据ts_code拉取每个公司的日线数据，拉取时会先ts_code分批，然后日期分批
// 注意限频
func (s *Stock) SaveAllIndexDaily(ctx context.Context) error {
	indexBasicResult, err := s.TuShare.Index(ctx, tushare.IndexParam{})
	if err != nil {
		return err
	}

	logrus.Infof("拉取指数列表成功，长度：%v", len(indexBasicResult.IndexList))

	for _, index := range indexBasicResult.IndexList {
		err = s.saveIndexDaily(ctx, index.TSCode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Stock) saveIndexDaily(ctx context.Context, tsCode string) error {
	startDate := "19900101"
	endDate := "20230901"
	indexDailyResult, err := s.TuShare.IndexDaily(ctx, tushare.IndexDailyParam{
		TSCode:    tsCode,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return err
	}

	err = s.IndexDailyDAL.BatchSaveIndexDaily(ctx, index_daily.BatchSaveIndexDailyParam{IndexDailyList: indexDailyResult.IndexDailyList})
	if err != nil {
		return err
	}

	logrus.Infof("保存 %v 数据成功, 长度: %v", tsCode, len(indexDailyResult.IndexDailyList))

	return nil
}
