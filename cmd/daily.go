package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/common/util"
	"my_stock_market/repo/interface/daily"
	"my_stock_market/service/interface/tushare"
	"sync"
	"time"
)

// SaveAllDaily 保存所有工作日&所有股票的日线数据
// 先拉取所有ts_code
// 然后根据ts_code拉取每个公司的日线数据，拉取时会先ts_code分批，然后日期分批
// 注意限频
func (s *Stock) SaveAllDaily(ctx context.Context) error {
	stockBasicResult, err := s.TuShare.StockBasic(ctx, tushare.StockBasicParam{Limit: 100000})
	if err != nil {
		return err
	}

	logrus.Infof("拉取股票列表成功，长度：%v", len(stockBasicResult.StockList))

	tsCode := make([]string, 0)
	for _, stock := range stockBasicResult.StockList {
		tsCode = append(tsCode, stock.TSCode)
	}

	wg := sync.WaitGroup{}
	batch := 900
	for i2 := 0; i2*batch < len(tsCode); i2++ {
		i := i2
		j := (i+1)*batch - 1
		if j > len(tsCode) {
			j = len(tsCode)
		}
		tempTSCode := make([]string, 0)
		for _, s := range tsCode[i*batch : j] {
			tempTSCode = append(tempTSCode, s)
		}

		logrus.Infof("开始第 %v 批数据", i+1)

		wg.Add(1)
		go func() {
			err = s.saveDaily(ctx, tempTSCode, i+1)
			if err != nil {
				panic(err)
			}
			wg.Add(-1)
		}()
	}

	wg.Wait()
	return nil
}

func (s *Stock) saveDaily(ctx context.Context, tsCode []string, c int) error {
	logrus.Infof("第 %v 批，tsCode前10: %v", c, tsCode[:10])
	startDate := "20230101"
	for {
		endDate, err := util.AddTime(startDate, 6)
		if err != nil {
			logrus.Errorf("错误: %v", err)
			return err
		}

		dailyResult, err := s.TuShare.Daily(ctx, tushare.DailyParam{
			TSCode:    tsCode,
			StartDate: startDate,
			EndDate:   endDate,
		})
		if err != nil {
			return err
		}

		err = s.StockDailyDAL.BatchSaveStockDaily(ctx, daily.BatchSaveStockDailyParam{StockDailyList: dailyResult.StockDailyList})
		if err != nil {
			return err
		}

		logrus.Infof("保存 %v-%v 数据成功, 第%v批", startDate, endDate, c)
		if len(dailyResult.StockDailyList) <= 500 {
			logrus.Info("休眠1s")
			time.Sleep(1 * time.Second)
		}

		ok, err := util.TimeCompare(endDate, "20230601")
		if err != nil {
			logrus.Errorf("错误: %v", err)
			return err
		}
		if ok {
			break
		}
		startDate, err = util.AddTime(endDate, 1)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Stock) TestTradeCal(ctx context.Context) error {
	TradeCalResult, err := s.TuShare.TradeCal(ctx, tushare.TradeCalParam{
		StartDate: "19900101",
		EndDate:   "19910101",
		IsOpen:    "1",
	})
	if err != nil {
		return err
	}
	logrus.Info(util.ToJsonStr(TradeCalResult.TradeCalList))
	logrus.Info(len(TradeCalResult.TradeCalList))
	return nil
}

func (s *Stock) TestDaily(ctx context.Context) error {
	DailyResult, err := s.TuShare.Daily(ctx, tushare.DailyParam{
		TSCode:    []string{"000001.SZ"},
		TradeDate: "19901219",
	})
	if err != nil {
		return err
	}
	logrus.Info(util.ToJsonStr(DailyResult.StockDailyList))
	return nil
}
