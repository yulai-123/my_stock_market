package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/repo/interface/daily_basic"
	"my_stock_market/service/interface/tushare"
	"sync"
	"time"
)

func (s *Stock) SaveAllDailyBasic(ctx context.Context) error {
	// 上海证劵交易所正式营业时间1990.12.19
	tradeCalResult, err := s.TuShare.TradeCal(ctx, tushare.TradeCalParam{
		StartDate: "20100101",
		EndDate:   "20230430",
		IsOpen:    "1",
	})
	if err != nil {
		return err
	}

	logrus.Infof("拉取交易日历成功，长度: %v", len(tradeCalResult.TradeCalList))

	wg := sync.WaitGroup{}
	batch := 1500
	for i := 0; i*batch < len(tradeCalResult.TradeCalList); i++ {
		j := (i+1)*batch - 1
		if j > len(tradeCalResult.TradeCalList) {
			j = len(tradeCalResult.TradeCalList)
		}

		tempTradeCal := tradeCalResult.TradeCalList[i*batch : j]
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := s.saveDailyBasic(ctx, tempTradeCal)
			if err != nil {
				panic(err)
			}
		}()
	}

	wg.Wait()

	logrus.Infof("保存结束")

	return nil
}

func (s *Stock) saveDailyBasic(ctx context.Context, tradeCalList []*tushare.SingleTradeCal) error {
	count := 0
	for _, tradeCal := range tradeCalList {
		logrus.Infof("开始执行: %v", tradeCal.CalDate)
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
