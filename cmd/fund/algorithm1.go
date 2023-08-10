package fund

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/model/do"
	"my_stock_market/repo/interface/fund_basic"
	"my_stock_market/repo/interface/fund_daily"
	"sort"
	"strings"
	"sync"
)

var (
	maxPosition     = float64(40000)
	initialPosition = float64(10000)

	singlePosition = float64(1000)

	ExcelPath = "excel/"
	tableHead = map[string]string{
		"A1": "名称",
		"B1": "TSCode",
		"C1": "上市时间",
		"D1": "最大回撤",
		"E1": "最大盈利",
		"F1": "平均盈利",
	}
)

type ResultData struct {
	TSCode     string
	Name       string
	ListDate   string
	MaxDecline float64
	MaxResult  float64
	AvgResult  float64
}

// Algorithm1
/*
1. 实时计算30日均线 -> ma30
2. 初始仓位为0，当遇到低于ma30时，开始买入 10000
3. 单位换算，将1%涨跌幅，换算成 涨跌多少元 -> single
4. 每下跌 single，买入 1000；每上涨 single，卖出 1000
5. 每14天，重新计算一次single
6. 如果出现所有都卖完的情况，回到2
7. 如果出现不断下跌，最多买入到 15000 (跌10个点），然后进入休眠。当价格回升到 跌10个点 上方时，开始恢复，回到4
*/
func (f *Fund) Algorithm1(ctx context.Context) error {
	// 场内，上市中，不然数据量太大
	//fundBasicResult, err := f.TuShare.FundBasic(ctx, tushare.FundBasicParam{
	//	Market: "E",
	//	Status: "L",
	//})
	//if err != nil {
	//	return err
	//}

	fundBasicResult, err := f.FundBasicDAL.GetAllFundBasic(ctx, fund_basic.GetAllFundBasicParam{})
	if err != nil {
		return err
	}

	// 待处理标的
	fundMap := map[string]*do.FundBasic{}

	targetTsCodeList := []string{
		"512040",
	}

	for _, fundBasic := range fundBasicResult.FundBasicList {
		flag := true
		if len(targetTsCodeList) > 0 {
			for _, tsCode := range targetTsCodeList {
				if !strings.Contains(fundBasic.TSCode, tsCode) {
					flag = false
					break
				}
			}
		}
		if !strings.Contains(fundBasic.Name, "ETF") {
			flag = false
		}
		if strings.Compare(fundBasic.ListDate, "20190101") > 0 {
			flag = false
		}
		if flag {
			fundMap[fundBasic.TSCode] = fundBasic
		}
	}
	fundList := make([]*do.FundBasic, 0)
	for _, fundBasic := range fundMap {
		fundList = append(fundList, fundBasic)
	}
	sort.Slice(fundList, func(i, j int) bool {
		return strings.Compare(fundList[i].TSCode, fundList[j].TSCode) <= 0
	})
	logrus.Infof("数据长度: %v", len(fundList))

	wg := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		start := i * 5
		if start > len(fundList) {
			break
		}
		wg.Add(1)
		go func(start, end int) {
			logrus.Infof("开始执行: %v -> %v", start, end)
			for s := start; s < end && s < len(fundList); s++ {
				_, _, _, err := f.algorithm1(ctx, fundList[s], false)
				if err != nil {
					panic(err)
				}
			}
			logrus.Infof("执行完成: %v -> %v", start, end)
			wg.Add(-1)
		}(start, start+5)
	}
	logrus.Infof("开始等待")
	wg.Wait()
	logrus.Infof("等待结束")

	logResultListWithMinProfit := make([]*LogResult, 0)
	logResultMapWithMinProfit.Range(func(key, value interface{}) bool {
		logResultListWithMinProfit = append(logResultListWithMinProfit, value.(*LogResult))
		return true
	})
	sort.Slice(logResultListWithMinProfit, func(i, j int) bool {
		return logResultListWithMinProfit[i].MinProfit > logResultListWithMinProfit[j].MinProfit
	})
	logrus.Infof("--------------开始打印minProfit前十股票-----------")
	for i, logResult := range logResultListWithMinProfit {
		if i >= 10 {
			break
		}
		logrus.Infof("名称：%v, tsCode: %v, 周期长度: %v, 最少盈利: %v，平均盈利: %v, 发行时间: %v, ma: %v, ma2: %v",
			logResult.Name, logResult.TSCode, logResult.TargetTimeLen, logResult.MinProfit,
			logResult.AvgProfit, logResult.ListDate, logResult.MA, logResult.MA2)
	}

	logResultListWithAvgProfit := make([]*LogResult, 0)
	logResultMapWithAvgProfit.Range(func(key, value interface{}) bool {
		logResultListWithAvgProfit = append(logResultListWithAvgProfit, value.(*LogResult))
		return true
	})
	sort.Slice(logResultListWithAvgProfit, func(i, j int) bool {
		return logResultListWithAvgProfit[i].AvgProfit > logResultListWithAvgProfit[j].AvgProfit
	})

	logrus.Infof("--------------开始打印avgProfit前十股票-----------")
	for i, logResult := range logResultListWithAvgProfit {
		if i >= 10 {
			break
		}
		logrus.Infof("名称：%v, tsCode: %v, 周期长度: %v, 最少盈利: %v，平均盈利: %v, 发行时间: %v, ma: %v, ma2: %v",
			logResult.Name, logResult.TSCode, logResult.TargetTimeLen, logResult.MinProfit,
			logResult.AvgProfit, logResult.ListDate, logResult.MA, logResult.MA2)
	}

	logResultListWithTargetTimeLen := make([]*LogResult, 0)
	logResultMapWithTargetTimeLen.Range(func(key, value interface{}) bool {
		logResultListWithTargetTimeLen = append(logResultListWithTargetTimeLen, value.(*LogResult))
		return true
	})
	sort.Slice(logResultListWithTargetTimeLen, func(i, j int) bool {
		return logResultListWithTargetTimeLen[i].TargetTimeLen < logResultListWithTargetTimeLen[j].TargetTimeLen
	})
	logrus.Infof("--------------开始打印targetTimeLen前十股票-----------")
	for i, logResult := range logResultListWithTargetTimeLen {
		if i >= 10 {
			break
		}
		logrus.Infof("名称：%v, tsCode: %v, 周期长度: %v, 最少盈利: %v，平均盈利: %v, 发行时间: %v, ma: %v, ma2: %v",
			logResult.Name, logResult.TSCode, logResult.TargetTimeLen, logResult.MinProfit,
			logResult.AvgProfit, logResult.ListDate, logResult.MA, logResult.MA2)
	}

	return nil
}

type DailyResult struct {
	Position     float64
	Number       int
	Single       float64
	CurrentPrice float64
	Money        float64
	DayCount     int
	Result       float64
	TradeDate    string
	RealResult   float64
}

func (f *Fund) algorithm1(ctx context.Context, fundBasic *do.FundBasic, needLog bool) (bool, float64, float64, error) {
	fundDailyList, err := f.getAllFundDaily(ctx, fundBasic.TSCode, "20170601")
	if err != nil {
		return false, 0, 0, err
	}

	for i := 1; i <= 100; i++ {
		if fundDailyList[len(fundDailyList)-i].Amount < 100 {
			logrus.Infof("成交量太低，名称: %v, tsCode: %v, amount: %v", fundBasic.Name, fundBasic.TSCode, fundDailyList[i].Amount)
			return false, 0, 0, err
		}
	}

	for ma := 1; ma <= 150; ma++ {
		for ma2 := 1; ma2 <= ma; ma2++ {

			// 初始仓位为 0
			position := float64(0)
			number := int(0)
			single := float64(0)
			// 当前购买价，用来计算涨跌的基准
			currentPrice := float64(0)
			// 钱
			money := float64(0)
			dayCount := 0

			dailyResultList := make([]*DailyResult, 0)

			for i, fundDaily := range fundDailyList {
				// 因为要计算30日均线，跳过前30天
				if i < ma {
					continue
				}
				// 获取30日均值
				ma30 := f.getMA(ctx, fundDailyList, i, ma)
				ma5 := f.getMA(ctx, fundDailyList, i, ma2)

				// 认为是0仓位
				if position < 10 {
					if fundDaily.Close < ma30 && fundDaily.Close > ma5 {
						tempNumber := int(initialPosition / fundDaily.Close)
						tempNumber = f.getNumber(tempNumber)

						number = number + tempNumber
						position = position + float64(tempNumber)*fundDaily.Close
						single = f.getSingle(fundDaily.Close)
						currentPrice = fundDaily.Close

						money = money - float64(tempNumber)*fundDaily.Close*1.0001
						dayCount = 0
					}

					dailyResultList = append(dailyResultList, &DailyResult{
						Position:     position,
						Number:       number,
						Single:       single,
						CurrentPrice: currentPrice,
						Money:        money,
						DayCount:     dayCount,
						Result:       position + money,
						TradeDate:    fundDaily.TradeDate,
						RealResult:   money + fundDaily.Close*float64(number),
					})

					continue
				}

				isSleep := false
				if position >= maxPosition {
					isSleep = true
				}

				if !isSleep || fundDaily.Close > currentPrice {
					change := fundDaily.Close - currentPrice
					upFlag := false
					if change > 0 {
						x := int(change / single)
						if x > 0 && position > 10 {
							upFlag = true
							t := 0
							for i := 1; i <= x; i++ {
								newPrice := currentPrice + single*float64(i)
								tempNumber := int(singlePosition / newPrice)
								tempNumber = f.getNumber(tempNumber)
								// 没有这么多，卖不出
								if number < tempNumber {
									break
								}
								t++
								number = number - tempNumber
								position = position - float64(tempNumber)*currentPrice

								money += float64(tempNumber) * newPrice * 0.9999
							}

							// 实际卖出 t 单位
							currentPrice += float64(t) * single
						}
					} else {
						// 价格下降多少单位
						x := int(change / single)
						if x > 0 && position < maxPosition {
							t := 0
							// 仓位上升x单位
							for i := 1; i <= x; i++ {
								newPrice := currentPrice - single*float64(i)
								tempNumber := int(singlePosition / newPrice)
								tempNumber = f.getNumber(tempNumber)

								number = number + tempNumber
								position = position + float64(tempNumber)*newPrice
								money = money - float64(tempNumber)*newPrice*1.0001
								t++
								if position >= maxPosition {
									break
								}
							}

							// 买入x单位
							currentPrice = currentPrice - float64(t)*single
						}
					}

					// 余额很少时，不管价格，直接按尾盘清仓
					tempNumber := int(singlePosition / single)
					tempNumber = f.getNumber(tempNumber)
					if number < tempNumber && upFlag {
						money += float64(number) * fundDaily.Close

						number = 0
						position = 0
					}
				}

				dailyResultList = append(dailyResultList, &DailyResult{
					Position:     position,
					Number:       number,
					Single:       single,
					CurrentPrice: currentPrice,
					Money:        money,
					DayCount:     dayCount,
					Result:       position + money,
					TradeDate:    fundDaily.TradeDate,
					RealResult:   money + fundDaily.Close*float64(number),
				})

				dayCount++
				if dayCount%5 == 0 {
					single = f.getSingle(fundDaily.Close)
				}
			}

			for i := 1; i <= 15; i++ {
				targetTimeLen := -10000
				minProfit := float64(10000)
				avgProfit := float64(0)

				timeLen := i * 20
				minProfit = float64(10000)

				daySum := 0
				avgProfit = 0

				for current := timeLen; current < len(dailyResultList); current++ {
					start := current - timeLen
					profit := dailyResultList[current].RealResult - dailyResultList[start].RealResult

					if profit < minProfit {
						minProfit = profit
					}
					daySum++
					avgProfit += profit
				}
				targetTimeLen = timeLen
				avgProfit = avgProfit / float64(daySum)
				minProfit /= float64(targetTimeLen)
				avgProfit /= float64(targetTimeLen)

				if minProfit > 2 {
					logrus.Infof("周期: %v, ma: %v, ma2: %v, minProfit: %v, avgProfit: %v",
						timeLen, ma, ma2, minProfit, avgProfit)
				}

				if true {
					logResult, exist := logResultMapWithMinProfit.Load(fundBasic.TSCode)
					if !exist || (exist && minProfit > logResult.(*LogResult).MinProfit) {
						logResultMapWithMinProfit.Store(fundBasic.TSCode, &LogResult{
							Name:          fundBasic.Name,
							TSCode:        fundBasic.TSCode,
							MinProfit:     minProfit,
							AvgProfit:     avgProfit,
							ListDate:      fundBasic.ListDate,
							TargetTimeLen: targetTimeLen,
							MA:            ma,
							MA2:           ma2,
						})
					}
					logResult, exist = logResultMapWithAvgProfit.Load(fundBasic.TSCode)
					if !exist || (exist && avgProfit > logResult.(*LogResult).AvgProfit) {
						logResultMapWithAvgProfit.Store(fundBasic.TSCode, &LogResult{
							Name:          fundBasic.Name,
							TSCode:        fundBasic.TSCode,
							MinProfit:     minProfit,
							AvgProfit:     avgProfit,
							ListDate:      fundBasic.ListDate,
							TargetTimeLen: targetTimeLen,
							MA:            ma,
							MA2:           ma2,
						})
					}
					logResult, exist = logResultMapWithTargetTimeLen.Load(fundBasic.TSCode)
					if !exist || (exist && targetTimeLen < logResult.(*LogResult).TargetTimeLen) {
						logResultMapWithTargetTimeLen.Store(fundBasic.TSCode, &LogResult{
							Name:          fundBasic.Name,
							TSCode:        fundBasic.TSCode,
							MinProfit:     minProfit,
							AvgProfit:     avgProfit,
							ListDate:      fundBasic.ListDate,
							TargetTimeLen: targetTimeLen,
							MA:            ma,
							MA2:           ma2,
						})
					}

				}
			}

		}
	}

	return false, 0, 0, nil
}

var (
	logResultMapWithMinProfit     = sync.Map{}
	logResultMapWithAvgProfit     = sync.Map{}
	logResultMapWithTargetTimeLen = sync.Map{}
)

type LogResult struct {
	Name          string
	TSCode        string
	MinProfit     float64
	AvgProfit     float64
	ListDate      string
	TargetTimeLen int
	MA            int
	MA2           int
}

// getNumber 因为购买要求，每次都要以100股为单位，因此也向上取整
func (f *Fund) getNumber(tempNumber int) int {
	return (tempNumber/100 + 1) * 100
}

// getSingle 因为ETF单价都很低，变动1分钱，可能就变动一个点
// 线上允许关注点最小单位是0.001元，因此要对single进行向上取整
func (f *Fund) getSingle(closePrice float64) float64 {
	a := int(closePrice * 1000 / 100)
	// 向上+1，保证大于一个点
	a += 1
	return float64(a) / 1000
}

func (f *Fund) getMA(ctx context.Context, dailyList []*do.FundDaily, currentDay int, ma int) float64 {

	sum := float64(0)

	for j := 1; j <= ma; j++ {
		sum = sum + dailyList[currentDay-j].Close
	}

	sum /= float64(ma)

	return sum
}

func (f *Fund) getAllFundDaily(ctx context.Context, tsCode, startDate string) ([]*do.FundDaily, error) {
	result := make([]*do.FundDaily, 0)

	//endTime, _ := util.AddTime(startDate, 100)

	dailyResult, err := f.FundDailyDAL.BatchGetFundDaily(ctx, fund_daily.BatchGetFundDailyParam{
		TSCode:    []string{tsCode},
		StartTime: startDate,
		//EndTime:   endTime,
	})
	if err != nil {
		return nil, err
	}

	result = dailyResult.FundDailyOfTSCodeMap[tsCode]
	sort.Slice(result, func(i, j int) bool {
		return strings.Compare(result[i].TradeDate, result[j].TradeDate) <= 0
	})

	return result, nil
}
