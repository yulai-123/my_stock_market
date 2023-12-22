package stock

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/repo/interface/daily_basic"
	"my_stock_market/service/interface/tushare"
	"time"
)

// SaveAllDailyBasic 保存所有工作日&所有股票的基本指标
func (s *Stock) SaveAllDailyBasic(ctx context.Context) error {
	// 上海证劵交易所正式营业时间1990.12.19
	tradeCalResult, err := s.TuShare.TradeCal(ctx, tushare.TradeCalParam{
		StartDate: "20000101",
		EndDate:   "20240430",
		IsOpen:    "1",
	})
	if err != nil {
		return err
	}

	logrus.Infof("拉取交易日历成功，长度: %v", len(tradeCalResult.TradeCalList))

	err = s.saveDailyBasic(ctx, tradeCalResult.TradeCalList)
	if err != nil {
		return err
	}

	logrus.Infof("保存结束")

	return nil
}

func (s *Stock) saveDailyBasic(ctx context.Context, tradeCalList []*tushare.SingleTradeCal) error {
	count := 0
	for i, tradeCal := range tradeCalList {
		logrus.Infof("开始执行: %v, 进度: %v/%v", tradeCal.CalDate, i, len(tradeCalList))
		dailyBasicResult, err := s.TuShare.DailyBasic(ctx, tushare.DailyBasicParam{
			TradeDate: tradeCal.CalDate,
		})
		if err != nil {
			logrus.Errorf("错误, error: %v, 休眠一秒", err)
			time.Sleep(1 * time.Second)
			count++
			if count >= 10 {
				return err
			}
			continue
		}
		err = s.DailyBasicDAL.BatchSaveDailyBasic(ctx, daily_basic.BatchSaveDailyBasicParam{DailyBasicList: dailyBasicResult.DailyBasicList})
		if err != nil {
			return err
		}
		time.Sleep(300 * time.Millisecond)
	}

	return nil
}

func (s *Stock) TestDailyBasic(ctx context.Context) error {
	dailyBasicResult, err := s.TuShare.DailyBasic(ctx, tushare.DailyBasicParam{
		TradeDate: "20230428",
	})
	if err != nil {
		return err
	}
	logrus.Infof("结果长度：%v", len(dailyBasicResult.DailyBasicList))
	logrus.Infof("%+v", dailyBasicResult.DailyBasicList[0])
	return nil
}
