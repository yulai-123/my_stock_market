package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/common/util"
	"my_stock_market/model/do"
	"my_stock_market/repo/interface/daily"
	"my_stock_market/repo/interface/stock"
	"my_stock_market/service/interface/tushare"
	"sort"
	"strings"
	"sync"
	"time"
)

// SaveAllDaily 保存所有工作日&所有股票的日线数据
// 先拉取所有工作日历
// 然后根据每一个工作日历获取所有股票数据进行保存
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

// CheckAlgorithm1 测试算法1
// 1. 在当天收盘时，如果涨幅超过9.5，则买入
// 2. 第二天收盘时涨幅为正，继续持有；第二天收盘时涨幅为负，直接出
// 3. 如果继续持有，一直循环第二天，直到结束
// 特殊点：
// 1> 只记录主板和中小板，不记录科创板和创业板
func (s *Stock) CheckAlgorithm1(ctx context.Context) error {
	for plateNumber := 0; plateNumber <= 5; plateNumber++ {
		for inPct := 9.5; inPct <= 10; inPct += 0.1 {
			startTime := "20100101"

			for {
				//获取所有股票
				result, err := s.StockDAL.GetAllStock(ctx, stock.GetAllStockParam{})
				if err != nil {
					return err
				}

				stockCount := 0
				changeSum := float64(0)
				changeSumAvg := float64(0)
				daySum := 0

				moneySum := 0

				for _, stock := range result.StockList {
					if stock.Market != "主板" && stock.Market != "中小板" {
						continue
					}
					stockCount++
					dailyList, err := s.getDaily(ctx, stock.TSCode)
					if err != nil {
						return err
					}

					for dayNumber, d := range dailyList {
						if strings.Compare(d.TradeDate, startTime) < 0 {
							continue
						}
						//小于9.5，不买入
						if !s.IsBuy(ctx, plateNumber, dailyList, dayNumber, inPct) {
							continue
						}
						// 总涨跌额
						change := float64(100)
						p := 1
						// 买入数量
						money := getBuyPrice(ctx, d)
						inMoney := money
						for ; (p + dayNumber) < len(dailyList); p++ {
							j := dayNumber + p
							money = money * (100 + dailyList[j].PctChg) / 100
							// 持有时收盘涨幅为负，卖出
							if dailyList[j].PctChg < 0 {
								break
							}
						}
						//修正最后一天
						if dayNumber+1 >= len(dailyList) {
							p = 0
						}
						daySum += p
						changeSum += change
						// 加入手续费
						changeSum -= 0.6
						changeSumAvg = changeSum / float64(daySum)
					}
				}

				logrus.Infof("只在%v板买入：入口: %v, 开始时间: %v, 计算股票数量: %v, 持有天数: %v, 平均总涨跌幅: %v, 手续费: %v", plateNumber, inPct, startTime, stockCount, daySum, changeSumAvg, 0.6)

				startTime, _ = util.AddTime(startTime, 365)
				ok, _ := util.TimeCompare(startTime, "20230428")
				if ok {
					break
				}
			}
		}
	}

	return nil
}

func getBuyPrice(ctx context.Context, d *do.StockDaily) float64 {
	//预期2000
	if d.Close > 20 {
		return d.Close * 100
	}
	a := d.Close * 100
	for a < float64(2000) {
		a += d.Close * 100
	}

	return a
}

var dailyMap = make(map[string][]*do.StockDaily)

func (s *Stock) getDaily(ctx context.Context, tsCode string) ([]*do.StockDaily, error) {
	dailyList, ok := dailyMap[tsCode]
	if ok {
		return dailyList, nil
	}
	dailyResult, err := s.StockDailyDAL.BatchGetStockDaily(ctx,
		daily.BatchGetStockDailyParam{TSCode: []string{tsCode}})
	if err != nil {
		return nil, err
	}
	dailyList = dailyResult.StockDailyOfTSCodeMap[tsCode]
	sort.Slice(dailyList, func(i, j int) bool {
		return strings.Compare(dailyList[i].TradeDate, dailyList[j].TradeDate) < 0
	})
	dailyMap[tsCode] = dailyList

	return dailyList, nil
}

func (s *Stock) IsBuy(ctx context.Context, plateNumber int, dailyList []*do.StockDaily, dayNumber int, inPct float64) bool {
	if dailyList[dayNumber].Close > 50 {
		return false
	}

	if plateNumber <= 0 {
		if dailyList[dayNumber].PctChg < inPct {
			return false
		}
		return true
	}

	if plateNumber == 1 {
		if dailyList[dayNumber].PctChg < inPct {
			return false
		}
		if dayNumber != 0 && dailyList[dayNumber-1].PctChg >= inPct {
			return false
		}
		return true
	}

	if plateNumber == 2 {
		if dailyList[dayNumber].PctChg < inPct {
			return false
		}
		if dayNumber == 0 {
			return false
		}
		if dailyList[dayNumber-1].PctChg < inPct {
			return false
		}
		return true
	}

	if plateNumber == 3 {
		if dailyList[dayNumber].PctChg < inPct {
			return false
		}
		if dayNumber < 2 {
			return false
		}
		if dailyList[dayNumber-1].PctChg < inPct || dailyList[dayNumber-2].PctChg < inPct {
			return false
		}
		return true
	}

	if plateNumber == 4 {
		if dailyList[dayNumber].PctChg < inPct {
			return false
		}
		if dayNumber < 3 {
			return false
		}
		if dailyList[dayNumber-1].PctChg < inPct || dailyList[dayNumber-2].PctChg < inPct || dailyList[dayNumber-3].PctChg < inPct {
			return false
		}
		return true
	}

	if plateNumber == 5 {
		if dailyList[dayNumber].PctChg < inPct {
			return false
		}
		if dayNumber < 4 {
			return false
		}
		if dailyList[dayNumber-1].PctChg < inPct || dailyList[dayNumber-2].PctChg < inPct || dailyList[dayNumber-3].PctChg < inPct || dailyList[dayNumber-4].PctChg < inPct {
			return false
		}
		return true
	}

	return false
}

// Algorithm1
// 获取这两天收盘超过9.5的股票
// 1. 要求主板或中小板
// 2. 不考虑今天停牌的
// 3. 隐形，手动筛选，股价超过50的
func (s *Stock) GetTsCodeWithAlgorithm1(ctx context.Context) error {
	now := time.Now()
	currentDay := now.Format("20060102")
	logrus.Info("今天是: ", currentDay)

	//获取所有股票
	result, err := s.StockDAL.GetAllStock(ctx, stock.GetAllStockParam{})
	if err != nil {
		return err
	}

	count := 0
	for _, stock := range result.StockList {
		if stock.Market != "主板" && stock.Market != "中小板" {
			continue
		}

		dailyList, err := s.getDaily(ctx, stock.TSCode)
		if err != nil {
			return err
		}

		if len(dailyList) < 3 {
			logrus.Errorf("%v", util.ToJsonStr(dailyList))
			panic("应该都超过三个")
		}

		i, j := len(dailyList)-1, len(dailyList)-2
		dayi, dayj := dailyList[i], dailyList[j]
		if dayi.TradeDate != currentDay {
			continue
		}
		//if dayi.TradeDate != "20230504" || dayj.TradeDate != "20230428" || dayk.TradeDate != "20230427" {
		//	logrus.Errorf("%v, %v, %v", util.ToJsonStr(dayi), util.ToJsonStr(dayj), util.ToJsonStr(dayk))
		//	panic("日期不对")
		//}

		if dayi.PctChg < 9.5 || dayj.PctChg < 9.5 {
			continue
		}
		count++

		logrus.Infof("第 %v 支股票, tsCode: %v, 名字: %v, 第一天: %v, 第二天: %v", count, stock.TSCode, stock.Name, dayi.PctChg, dayj.PctChg)
	}

	return nil
}
