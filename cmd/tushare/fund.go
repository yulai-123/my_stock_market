package tushare

import (
	"context"
	"fmt"
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
	startDate := "20210101"
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

		if len(fundDailyResult.FundDailyList) <= 0 {
			return nil
		}

		fundAdjResult, err := s.TuShare.FundAdj(ctx, tushare.FundAdjParam{
			TSCode:    tsCode,
			StartDate: startDate,
			EndDate:   endDate,
		})
		if err != nil {
			return err
		}

		for _, fundDaily := range fundDailyResult.FundDailyList {
			flag := false

			for i := 0; i < len(fundDailyResult.FundDailyList); i++ {
				date, _ := util.AddTime(fundDaily.TradeDate, int64(-1*i))
				for _, adj := range fundAdjResult.FundAdjList {
					if date == adj.TradeDate {
						flag = true
						// 执行后复权
						fundDaily.Close *= adj.AdjFactor
						fundDaily.Close = float64(int(fundDaily.Close*10000)) / 10000
						fundDaily.High *= adj.AdjFactor
						fundDaily.High = float64(int(fundDaily.High*10000)) / 10000
						fundDaily.Low *= adj.AdjFactor
						fundDaily.Low = float64(int(fundDaily.Low*10000)) / 10000
						fundDaily.Open *= adj.AdjFactor
						fundDaily.Open = float64(int(fundDaily.Open*10000)) / 10000
						// 昨收额不准确
						fundDaily.Change *= adj.AdjFactor
						fundDaily.Change = float64(int(fundDaily.Change*10000)) / 10000
						break
					}
				}
				if flag {
					break
				}
			}

			if !flag {
				logrus.Errorf("not find fund_adj, ts_code: %v, trade_date: %v", fundDaily.TSCode, fundDaily.TradeDate)
				err := fmt.Errorf("not find fund_adj")
				return err
			}
		}

		err = s.FundDailyDAL.BatchSaveFundDaily(ctx, fund_daily.BatchSaveFundDailyParam{FundDailyList: fundDailyResult.FundDailyList})
		if err != nil {
			return err
		}

		logrus.Infof("保存 %v 数据成功, 长度: %v, startDate: %v, endDate: %v", tsCode, len(fundDailyResult.FundDailyList), startDate, endDate)

		endDate, err = util.AddTime(startDate, -1)
		if err != nil {
			return err
		}
		startDate, err = util.AddTime(startDate, -1000)
		if err != nil {
			return err
		}

		time.Sleep(250 * time.Millisecond)
	}

	return nil
}
