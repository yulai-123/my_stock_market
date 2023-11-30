package algorithm2

import (
	"context"
	"github.com/sirupsen/logrus"
	chart2 "my_stock_market/common/chart"
	"my_stock_market/model/do"
	"my_stock_market/repo/interface/fund_basic"
	"my_stock_market/repo/interface/fund_daily"
	"sort"
	"strings"
)

func (a *Algorithm2) Run(ctx context.Context) error {
	fundList, err := a.GetTargetFund(ctx)
	if err != nil {
		return err
	}
	chart := chart2.NewChart(ctx, "网格交易")

	for _, fundBasic := range fundList {
		dataList, err := a.Exec(ctx, fundBasic)
		if err != nil {
			return err
		}
		chart.AddResultData(dataList)
	}

	chart.WithMoney = true
	chart.Render(ctx)

	return nil
}

func (a *Algorithm2) Exec(ctx context.Context, fundBasic *do.FundBasic) ([]*do.ResultData, error) {
	result := make([]*do.ResultData, 0)

	maDayCount, ma2DayCount := 20, 5
	initialMoney := float64(10000)
	singleMoney := float64(1000)
	maxPosition := float64(40000)
	isSleep := false

	// 获取日线数据
	fundDailyList, err := a.getAllFundDaily(ctx, fundBasic.TSCode, "20170601")
	if err != nil {
		return nil, err
	}

	// 开始运行
	position := float64(0)     //当前仓位
	money := float64(0)        //当前金钱
	currentNumber := 0         //当前数量
	currentPrice := float64(0) //当前单价
	single := float64(0)       // 单价变动最小单位
	level := 0

	for currentDay, fundDaily := range fundDailyList {
		// 为了方便计算ma
		if currentDay < 30 {
			continue
		}

		ma := a.getMA(ctx, fundDailyList, currentDay, maDayCount)
		ma2 := a.getMA(ctx, fundDailyList, currentDay, ma2DayCount)

		dayCount := 0

		// 空仓时，会根据 ma&m2 进行判断是否买入
		// 在当前这天的末尾判断是否买入
		if position <= 10 {
			// 买入条件
			if fundDaily.Close < ma && fundDaily.Close > ma2 {
				// 买入数量
				tempNumber := int(initialMoney / fundDaily.Close)
				realNumber := a.getNumber(tempNumber)

				// 买入结果
				money -= float64(realNumber) * fundDaily.Close * 1.0001 // 花费金额
				position += float64(realNumber) * fundDaily.Close       // 仓位上涨
				currentNumber += realNumber
				currentPrice = fundDaily.Close        // 仓位价格，会随仓位变动，价格变动
				single = a.getSingle(fundDaily.Close) // 单价变动最小单位
				dayCount = 0
			}

			result = append(result, &do.ResultData{
				Result:    money + float64(currentNumber)*fundDaily.Close,
				Money:     money,
				TradeDate: fundDaily.TradeDate,
				TSCode:    fundDaily.TSCode,
			})

			continue
		}

		dayCount++

		// 有仓位的时候，根据收盘价进行计算买卖
		// 因为网格交易利于震荡。根据收盘价进行计算，只会算出收益更少的情况，不会出现虚高

		change := fundDaily.Close - currentPrice
		if change > 0 {
			if isSleep {
				isSleep = false
			}
			if change < single {
				result = append(result, &do.ResultData{
					Result:    money + float64(currentNumber)*fundDaily.Close,
					Money:     money,
					TradeDate: fundDaily.TradeDate,
					TSCode:    fundDaily.TSCode,
				})
				continue
			}

			changeCount := int(change / single)
			sumNumber := 0
			singleNumber := int(singleMoney / currentPrice)
			realNumber := a.getNumber(singleNumber)
			for i := 1; i <= changeCount; i++ {
				if realNumber > currentNumber {
					break
				}

				// 卖出
				currentPrice = currentPrice + single
				money += currentPrice * float64(realNumber) * 0.9999
				currentNumber -= realNumber
				position -= currentPrice * float64(realNumber)
				sumNumber += realNumber
			}

			// 比过去30天最高值每低5%，则卖出
			minClose := float64(100000)
			for i := 1; i <= 30; i++ {
				if fundDailyList[currentDay-i].Close < minClose {
					minClose = fundDailyList[currentDay-i].Close
				}
			}
			closeChange := (fundDaily.Close - minClose) / minClose
			//closeChangeCount := closeChange / 0.05
			if closeChange > float64(level*-1)*0.05 && level != 0 {

				tempMoney := float64(level * -1 * 10000)
				tempNumber := int(tempMoney / currentPrice)
				tempNumber = a.getNumber(tempNumber)

				if currentNumber < tempNumber {
					money += currentPrice * float64(currentNumber) * 0.9999
					currentNumber = 0
					position = 0
					level = 0
				} else {
					money += currentPrice * float64(tempNumber) * 0.9999
					currentNumber -= tempNumber
					position -= currentPrice * float64(tempNumber)
					level = 0
				}

				logrus.Infof("tsCode: %v, tradeDate: %v, closeChange: %v, level: %v", fundDaily.TSCode, fundDaily.TradeDate, closeChange, level)
			}

			// 剩余数量不多，一次按最后价格卖出
			if currentNumber < realNumber {
				money += currentPrice * float64(currentNumber) * 0.9999
				currentNumber = 0
				position = 0

				result = append(result, &do.ResultData{
					Result:    money + float64(currentNumber)*fundDaily.Close,
					Money:     money,
					TradeDate: fundDaily.TradeDate,
					TSCode:    fundDaily.TSCode,
				})

				continue
			}

			result = append(result, &do.ResultData{
				Result:    money + float64(currentNumber)*fundDaily.Close,
				Money:     money,
				TradeDate: fundDaily.TradeDate,
				TSCode:    fundDaily.TSCode,
			})
		} else {
			change *= -1

			if isSleep || change < single {
				result = append(result, &do.ResultData{
					Result:    money + float64(currentNumber)*fundDaily.Close,
					Money:     money,
					TradeDate: fundDaily.TradeDate,
					TSCode:    fundDaily.TSCode,
				})
				continue
			}

			changeCount := int(change / single)
			sumNumber := 0
			singleNumber := int(singleMoney / currentPrice)
			realNumber := a.getNumber(singleNumber)

			for i := 1; i <= changeCount; i++ {
				// 如果大于最大仓位，进入休眠
				if position > maxPosition {
					break
				}

				currentPrice = currentPrice - single
				money -= currentPrice * float64(realNumber) * 1.0001
				currentNumber += realNumber
				sumNumber += realNumber
				position += currentPrice * float64(realNumber)
			}

			// 比过去30天最高值每低5%，则卖出
			maxClose := float64(-100000)
			for i := 1; i <= 30; i++ {
				if fundDailyList[currentDay-i].Close > maxClose {
					maxClose = fundDailyList[currentDay-i].Close
				}
			}
			closeChange := (maxClose - fundDaily.Close) / maxClose
			if closeChange > float64(level)*-1*0.05 {
				// 没低5个点，买入一万
				for ; closeChange > float64(level)*-1*0.05; level -= 1 {
					// 如果大于最大仓位，进入休眠
					if position > maxPosition {
						break
					}

					tempNumber := int(10000 / currentPrice)
					tempNumber = a.getNumber(tempNumber)

					money -= currentPrice * float64(tempNumber) * 1.0001
					currentNumber += tempNumber
					position += currentPrice * float64(tempNumber)
				}
				logrus.Infof("tsCode: %v, tradeDate: %v, closeChange: %v, level: %v, position: %v", fundDaily.TSCode, fundDaily.TradeDate, closeChange, level, position)
			}

			if position > maxPosition {
				isSleep = true
			}

			result = append(result, &do.ResultData{
				Result:    money + float64(currentNumber)*fundDaily.Close,
				Money:     money,
				TradeDate: fundDaily.TradeDate,
				TSCode:    fundDaily.TSCode,
			})
		}

		if dayCount%5 == 0 {
			single = a.getSingle(fundDaily.Close)
		}

	}

	return result, nil
}

func (a *Algorithm2) GetTargetFund(ctx context.Context) ([]*do.FundBasic, error) {
	fundBasicResult, err := a.FundBasicDAL.GetAllFundBasic(ctx, fund_basic.GetAllFundBasicParam{})
	if err != nil {
		return nil, err
	}

	// 待处理标的
	fundMap := map[string]*do.FundBasic{}

	targetTsCodeList := []string{
		"512040",
		//"512600",
	}

	for _, fundBasic := range fundBasicResult.FundBasicList {
		flag := true
		if len(targetTsCodeList) > 0 {
			flag = false
			for _, tsCode := range targetTsCodeList {
				if strings.Contains(fundBasic.TSCode, tsCode) {
					flag = true
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

	return fundList, nil
}

func (a *Algorithm2) getAllFundDaily(ctx context.Context, tsCode, startDate string) ([]*do.FundDaily, error) {
	result := make([]*do.FundDaily, 0)

	//endTime, _ := util.AddTime(startDate, 100)

	dailyResult, err := a.FundDailyDAL.BatchGetFundDaily(ctx, fund_daily.BatchGetFundDailyParam{
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

// getNumber 因为购买要求，每次都要以100股为单位，因此也向上取整
func (a *Algorithm2) getNumber(tempNumber int) int {
	return (tempNumber/100 + 1) * 100
}

// getSingle 因为ETF单价都很低，变动1分钱，可能就变动一个点
// 线上允许关注点最小单位是0.001元，因此要对single进行向上取整
func (a *Algorithm2) getSingle(closePrice float64) float64 {
	b := int(closePrice * 1000 / 100)
	// 向上+1，保证大于一个点
	b += 1
	return float64(b) / 1000
}

func (a *Algorithm2) getMA(ctx context.Context, dailyList []*do.FundDaily, currentDay int, ma int) float64 {

	sum := float64(0)

	for j := 1; j <= ma; j++ {
		sum = sum + dailyList[currentDay-j].Close
	}

	sum /= float64(ma)

	return sum
}
