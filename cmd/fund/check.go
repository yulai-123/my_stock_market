package fund

import (
	"context"
	"github.com/sirupsen/logrus"
)

func (f *Fund) Check(ctx context.Context) error {
	// 确定ts_code
	tsCode := "512040.SH"
	maDayCount, ma2DayCount := 49, 5
	initialMoney := float64(10000)
	singleMoney := float64(200)
	maxPosition := float64(100000)
	isSleep := false

	// 获取日线数据
	fundDailyList, err := f.getAllFundDaily(ctx, tsCode, "20170601")
	if err != nil {
		return err
	}

	// 开始运行
	position := float64(0)     //当前仓位
	money := float64(0)        //当前金钱
	currentNumber := 0         //当前数量
	currentPrice := float64(0) //当前单价
	single := float64(0)       // 单价变动最小单位

	for currentDay, fundDaily := range fundDailyList {
		// 为了方便计算ma
		if currentDay < maDayCount {
			continue
		}

		//if strings.Compare(fundDaily.TradeDate, "20230301") < 0 {
		//	continue
		//}

		ma := f.getMA(ctx, fundDailyList, currentDay, maDayCount)
		ma2 := f.getMA(ctx, fundDailyList, currentDay, ma2DayCount)
		ma = 10000000
		ma2 = -1

		dayCount := 0

		// 空仓时，会根据 ma&m2 进行判断是否买入
		// 在当前这天的末尾判断是否买入
		if position <= 10 {
			// 买入条件
			if fundDaily.Close < ma && fundDaily.Close > ma2 {
				// 买入数量
				tempNumber := int(initialMoney / fundDaily.Close)
				realNumber := f.getNumber2(tempNumber)

				// 买入结果
				money -= float64(realNumber) * fundDaily.Close * 1.0001 // 花费金额
				position += float64(realNumber) * fundDaily.Close       // 仓位上涨
				currentNumber += realNumber
				currentPrice = fundDaily.Close         // 仓位价格，会随仓位变动，价格变动
				single = f.getSingle2(fundDaily.Close) // 单价变动最小单位

				result := money + position
				realResult := money + float64(currentNumber)*fundDaily.Close
				logrus.Infof("第 %v 天，执行买入操作, ma: %v, ma2: %v, close: %v, 买入, realNumber: %v, "+
					"结果: money: %v, position: %v, currentNumber: %v, single: %v, result: %v, realResult: %v",
					currentDay+1, ma, ma2, fundDaily.Close, realNumber, money, position,
					currentNumber, single, result, realResult)
				dayCount = 0

				continue
			}

			result := money + position
			realResult := money + float64(currentNumber)*fundDaily.Close
			logrus.Infof("第 %v 天，空仓，未执行操作, ma: %v, ma2: %v, close: %v, "+
				"结果: money: %v, position: %v, currentNumber: %v, single: %v, result: %v, realResult: %v",
				currentDay+1, ma, ma2, fundDaily.Close, money, position,
				currentNumber, single, result, realResult)

			continue
		}

		dayCount++

		// 有仓位的时候，根据收盘价进行计算买卖
		// 因为网格交易利于震荡。根据收盘价进行计算，只会算出收益更少的情况，不会出现虚高

		change := fundDaily.Close - currentPrice
		if change > 0 {
			if isSleep {
				isSleep = false
				logrus.Infof("第 %v 天，解除休眠期", currentDay+1)
			}
			if change < single {
				result := money + position
				realResult := money + float64(currentNumber)*fundDaily.Close
				logrus.Infof("第 %v 天，上涨，单位太小，未执行操作, ma: %v, ma2: %v, close: %v, "+
					"结果: money: %v, position: %v, currentNumber: %v, single: %v, result: %v, realResult: %v, change: %v",
					currentDay+1, ma, ma2, fundDaily.Close, money, position,
					currentNumber, single, result, realResult, change)
				continue
			}

			changeCount := int(change / single)
			sumNumber := 0
			singleNumber := int(singleMoney / currentPrice)
			realNumber := f.getNumber2(singleNumber)
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

			// 剩余数量不多，一次按最后价格卖出
			if currentNumber < realNumber {
				money += currentPrice * float64(currentNumber) * 0.9999
				sumNumber += currentNumber
				currentNumber = 0
				position = 0

				result := money + position
				realResult := money + float64(currentNumber)*fundDaily.Close
				logrus.Infof("第 %v 天，上涨，执行卖出操作, 全卖出, 卖出：sumNumber: %v, ma: %v, ma2: %v, close: %v, "+
					"结果: money: %v, position: %v, currentNumber: %v, single: %v, result: %v, realResult: %v",
					currentDay+1, sumNumber, ma, ma2, fundDaily.Close, money, position,
					currentNumber, single, result, realResult)

				continue
			}

			result := money + position
			realResult := money + float64(currentNumber)*fundDaily.Close
			logrus.Infof("第 %v 天，上涨，执行卖出操作, 有剩余, 卖出：sumNumber: %v, ma: %v, ma2: %v, close: %v, "+
				"结果: money: %v, position: %v, currentNumber: %v, single: %v, result: %v, realResult: %v",
				currentDay+1, sumNumber, ma, ma2, fundDaily.Close, money, position,
				currentNumber, single, result, realResult)
		} else {
			if isSleep {
				result := money + position
				realResult := money + float64(currentNumber)*fundDaily.Close
				logrus.Infof("第 %v 天，下跌，正在休眠期，不操作, ma: %v, ma2: %v, close: %v, "+
					"结果: money: %v, position: %v, currentNumber: %v, single: %v, result: %v, realResult: %v",
					currentDay+1, ma, ma2, fundDaily.Close, money, position,
					currentNumber, single, result, realResult)
				continue
			}

			change *= -1
			if change < single {
				result := money + position
				realResult := money + float64(currentNumber)*fundDaily.Close
				logrus.Infof("第 %v 天，下跌，单位太小，未执行操作, ma: %v, ma2: %v, close: %v, "+
					"结果: money: %v, position: %v, currentNumber: %v, single: %v, result: %v, realResult: %v, change: %v",
					currentDay+1, ma, ma2, fundDaily.Close, money, position,
					currentNumber, single, result, realResult, change)
				continue
			}

			changeCount := int(change / single)
			sumNumber := 0
			singleNumber := int(singleMoney / currentPrice)
			realNumber := f.getNumber2(singleNumber)

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

			if position > maxPosition {
				isSleep = true

				result := money + position
				realResult := money + float64(currentNumber)*fundDaily.Close
				logrus.Infof("第 %v 天，下跌，执行买入操作, 达到最大仓位，进入休眠期，买入：sumNumber: %v, ma: %v, ma2: %v, close: %v, "+
					"结果: money: %v, position: %v, currentNumber: %v, single: %v, result: %v, realResult: %v",
					currentDay+1, sumNumber, ma, ma2, fundDaily.Close, money, position,
					currentNumber, single, result, realResult)
				continue
			}

			result := money + position
			realResult := money + float64(currentNumber)*fundDaily.Close
			logrus.Infof("第 %v 天，下跌，执行买入操作, 买入：sumNumber: %v, ma: %v, ma2: %v, close: %v, "+
				"结果: money: %v, position: %v, currentNumber: %v, single: %v, result: %v, realResult: %v",
				currentDay+1, sumNumber, ma, ma2, fundDaily.Close, money, position,
				currentNumber, single, result, realResult)
		}

		if dayCount%5 == 0 {
			single = f.getSingle2(fundDaily.Close)
		}

	}

	return nil
}

// getNumber2 因为购买要求，每次都要以100股为单位，因此也向上取整
func (f *Fund) getNumber2(tempNumber int) int {
	return (tempNumber/100 + 1) * 100
}

// getSingle2 因为ETF单价都很低，变动1分钱，可能就变动一个点
// 线上允许关注点最小单位是0.001元，因此要对single进行向上取整
func (f *Fund) getSingle2(closePrice float64) float64 {
	//a := int(closePrice * 1000 / 200)
	//// 向上+1，保证大于一个点
	//a += 1
	//return float64(a) / 1000
	return 0.002
}
