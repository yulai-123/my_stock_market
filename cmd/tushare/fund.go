package tushare

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/common/util"
	"my_stock_market/repo/interface/fund_basic"
	"my_stock_market/repo/interface/fund_daily"
	"my_stock_market/service/interface/tushare"
	"time"
)

func (s *Stock) SaveAllFundBasic(ctx context.Context) error {
	// 场内，上市中，不然数据量太大
	fundBasicResult, err := s.TuShare.FundBasic(ctx, tushare.FundBasicParam{
		Market: "E",
		Status: "L",
	})
	if err != nil {
		return err
	}

	logrus.Infof("拉取基金列表成功，长度：%v", len(fundBasicResult.FundBasicList))

	_, err = s.FundBasicDAL.BatchSaveFundBasic(ctx, fund_basic.BatchSaveFundBasicParam{
		FundBasicList: fundBasicResult.FundBasicList,
	})
	if err != nil {
		return err
	}

	logrus.Info("保存数据库成功")
	return nil
}

func (s *Stock) SaveAllFundDaily(ctx context.Context) error {
	fundBasicResult, err := s.TuShare.FundBasic(ctx, tushare.FundBasicParam{
		Market: "E",
		Status: "L",
	})
	if err != nil {
		return err
	}

	logrus.Infof("拉取基金列表成功，长度：%v", len(fundBasicResult.FundBasicList))

	for _, fundBasic := range fundBasicResult.FundBasicList {
		err = s.saveFundDaily(ctx, fundBasic.TSCode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Stock) saveFundDaily(ctx context.Context, tsCode string) error {
	startDate := "20220101"
	endDate := "20240101"

	for i := 0; i < 7; i++ {
		fundDailyResult, err := s.TuShare.FundDaily(ctx, tushare.FundDailyParam{
			TSCode:    tsCode,
			StartDate: startDate,
			EndDate:   endDate,
		})
		if err != nil {
			return err
		}

		err = s.FundDailyDAL.BatchSaveFundDaily(ctx, fund_daily.BatchSaveFundDailyParam{FundDailyList: fundDailyResult.FundDailyList})
		if err != nil {
			return err
		}

		logrus.Infof("保存 %v 数据成功, 长度: %v, startDate: %v, endDate: %v", tsCode, len(fundDailyResult.FundDailyList), startDate, endDate)

		if len(fundDailyResult.FundDailyList) <= 0 {
			return nil
		}

		endDate, err = util.AddTime(startDate, -1)
		if err != nil {
			return err
		}
		startDate, err = util.AddTime(startDate, -730)
		if err != nil {
			return err
		}

		time.Sleep(250 * time.Millisecond)
	}

	return nil
}
