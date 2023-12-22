package stock

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/common/util"
	"my_stock_market/repo/interface/monthly"
	"my_stock_market/service/interface/tushare"
	"sync"
	"time"
)

// SaveAllMonthly
// 股市月线
func (s *Stock) SaveAllMonthly(ctx context.Context) error {
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
			err = s.saveMonthly(ctx, tempTSCode, i+1)
			if err != nil {
				panic(err)
			}
			wg.Add(-1)
		}()
	}

	wg.Wait()
	return nil
}

func (s *Stock) saveMonthly(ctx context.Context, tsCode []string, c int) error {
	logrus.Infof("第 %v 批，tsCode前10: %v", c, tsCode[:10])
	startDate := "19900101"
	for {
		endDate, err := util.AddTime(startDate, 100)
		if err != nil {
			logrus.Errorf("错误: %v", err)
			return err
		}

		monthlyResult, err := s.TuShare.Monthly(ctx, tushare.MonthlyParam{
			TSCode:    tsCode,
			StartDate: startDate,
			EndDate:   endDate,
		})
		if err != nil {
			return err
		}

		err = s.StockMonthlyDAL.BatchSaveStockMonthly(ctx, monthly.BatchSaveStockMonthlyParam{StockMonthlyList: monthlyResult.StockMonthlyList})
		if err != nil {
			return err
		}

		logrus.Infof("保存 %v-%v 数据成功, 第%v批", startDate, endDate, c)
		if len(monthlyResult.StockMonthlyList) <= 500 {
			time.Sleep(1 * time.Second)
		}

		ok, err := util.TimeCompare(endDate, "20240430")
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
