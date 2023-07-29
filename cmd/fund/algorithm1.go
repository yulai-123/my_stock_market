package fund

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"my_stock_market/common/util"
	"my_stock_market/model/do"
	"my_stock_market/repo/interface/fund_daily"
	"my_stock_market/service/interface/tushare"
	"sort"
	"strconv"
	"strings"
)

var (
	maxPosition     = float64(25000)
	initialPosition = float64(10000)

	singlePosition = float64(1000)
	addDate        = int64(30)

	warningDecline  = float64(-20000)
	minTargetResult = float64(1000)

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
	fundBasicResult, err := f.TuShare.FundBasic(ctx, tushare.FundBasicParam{
		Market: "E",
		Status: "L",
	})
	if err != nil {
		return err
	}

	// 待处理标的
	fundMap := map[string]*do.FundBasic{}

	targetTsCodeList := []string{
		//"159839",
	}
	//targetTsCodeList = []string{"159959"}

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

	resultDataList := make([]ResultData, 0)
	starTSCodeList := make([]string, 0)
	for _, fundBasic := range fundList {
		flag, maxDecline, maxResult, err := f.algorithm1(ctx, fundBasic, false)
		if err != nil {
			return err
		}
		if flag {
			starTSCodeList = append(starTSCodeList, fundBasic.TSCode)
			fromYear, err := strconv.Atoi(fundBasic.ListDate[:4])
			if err != nil {
				return err
			}
			yearCount := 2023 - fromYear
			avgResult := maxResult / float64(yearCount)
			//logrus.Infof("重点关注：%v, tsCode: %v, 最大回撤: %v, 平均盈利: %v, 最大盈利: %v, 上市时间: %v", fundBasic.Name, fundBasic.TSCode, maxDecline, avgResult, maxResult, fundBasic.ListDate)

			resultDataList = append(resultDataList, ResultData{
				TSCode:     fundBasic.TSCode,
				Name:       fundBasic.Name,
				ListDate:   fundBasic.ListDate,
				MaxDecline: maxDecline,
				MaxResult:  maxResult,
				AvgResult:  avgResult,
			})
		}
	}

	excelFilePath := fmt.Sprintf("%v%v.xlsx", ExcelPath, "ETF数据")
	f2 := excelize.NewFile()
	defer func() {
		if err := f2.Close(); err != nil {
			logrus.Errorf("[saveAsExcel] 关闭excel失败: %v", err)
		}
	}()
	// 创建一个工作表
	index, err := f2.NewSheet("Sheet1")
	if err != nil {
		return err
	}
	// 设置单元格的值
	for cell, value := range tableHead {
		f2.SetCellValue("Sheet1", cell, value)
	}

	for r, resultData := range resultDataList {
		f2.SetCellValue("Sheet1", fmt.Sprintf("A%v", r+2), resultData.Name)
		f2.SetCellValue("Sheet1", fmt.Sprintf("B%v", r+2), resultData.TSCode)
		f2.SetCellValue("Sheet1", fmt.Sprintf("C%v", r+2), resultData.ListDate)
		f2.SetCellValue("Sheet1", fmt.Sprintf("D%v", r+2), resultData.MaxDecline)
		f2.SetCellValue("Sheet1", fmt.Sprintf("E%v", r+2), resultData.MaxResult)
		f2.SetCellValue("Sheet1", fmt.Sprintf("F%v", r+2), resultData.AvgResult)

	}
	// 设置工作簿的默认工作表
	f2.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := f2.SaveAs(excelFilePath); err != nil {
		return err
	}

	return nil
}

func (f *Fund) algorithm1(ctx context.Context, fundBasic *do.FundBasic, needLog bool) (bool, float64, float64, error) {
	tsCode := fundBasic.TSCode

	starFlag := true
	maxDecline := float64(0)
	maxResult := float64(0)

	if needLog {
		logrus.Infof("开始计算 tsCode: %v, 基金名称: %v", tsCode, fundBasic.Name)
	}
	// 从不同时间开始
	startDate := "20100101"
	for ; strings.Compare(startDate, "20240101") < 0; startDate, _ = util.AddTime(startDate, addDate) {
		currentMaxDecline := float64(0)
		// 获取数据
		fundDailyList, err := f.getAllFundDaily(ctx, tsCode, startDate)
		if err != nil {
			return false, 0, 0, err
		}

		// 运行时间少于两年，不考虑；低于两毛的，不考虑
		if len(fundDailyList) < 30 || fundDailyList[0].Close < 0.3 {
			continue
		}

		// 初始仓位为 0
		position := float64(0)
		number := int(0)
		single := float64(0)
		// 当前购买价，用来计算涨跌的基准
		currentPrice := float64(0)
		// 钱
		money := float64(0)
		dayCount := 0
		for i, fundDaily := range fundDailyList {
			if money+position < maxDecline {
				maxDecline = money + position
			}
			if money+position > maxResult {
				maxResult = money + position
			}
			if money+position < currentMaxDecline {
				currentMaxDecline = money + position
			}

			//logrus.Infof("---------tradeDate: %v, position: %v, number: %v, single: %v, currentPrice: %v, money: %v, dayCount: %v", fundDaily.TradeDate, position, number, single, currentPrice, money, dayCount)
			// 因为要计算30日均线，跳过前30天
			if i < 30 {
				continue
			}
			// 获取30日均值
			ma30 := f.getMA30(ctx, fundDailyList, i)

			// 认为是0仓位
			if position < 10 {
				if fundDaily.Close < ma30 {
					// 购买10000
					tempNumber := int(initialPosition / fundDaily.Close)
					tempNumber = f.getNumber(tempNumber)

					number = number + tempNumber
					position = position + float64(tempNumber)*fundDaily.Close
					single = f.getSingle(fundDaily.Close)
					currentPrice = fundDaily.Close

					money -= float64(tempNumber) * fundDaily.Close
					dayCount = 0
				}
				continue
			}

			// 根据最低点买入
			// 会买入
			decline := currentPrice - fundDaily.Low
			// 价格下降多少单位
			x := int(decline / single)
			if x > 0 && position < maxPosition {
				t := 0
				// 仓位上升x单位
				for i := 1; i <= x; i++ {
					newPrice := currentPrice - single*float64(i)
					tempNumber := int(singlePosition / newPrice)
					tempNumber = f.getNumber(tempNumber)

					number = number + tempNumber
					position = position + float64(tempNumber)*newPrice
					money -= float64(tempNumber) * newPrice
					t++
					if position >= maxPosition {
						break
					}
				}

				// 买入x单位
				currentPrice -= float64(t) * single
			}

			// 根据最高点卖出
			rise := fundDaily.High - currentPrice
			// 价格上涨x单位
			x = int(rise / single)
			if x > 0 && position > 10 {
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

					money += float64(tempNumber) * newPrice
				}

				// 实际卖出 t 单位
				currentPrice += float64(t) * single
			}

			// 余额很少时，不管价格，直接按尾盘清仓
			tempNumber := int(singlePosition / single)
			tempNumber = f.getNumber(tempNumber)
			if number < tempNumber {
				money += float64(number) * fundDaily.Close

				number = 0
				position = 0
			}

			dayCount++
			if dayCount%5 == 0 {
				single = f.getSingle(fundDaily.Close)
			}
		}

		if needLog {
			logrus.Infof("ts_code: %v, 基金: %v, actualStartDate: %v, result: %v, currentMaxDecline: %v, maxDecline: %v, money: %v, position: %v, startDate: %v", tsCode, fundBasic.Name, fundDailyList[0].TradeDate, money+position, currentMaxDecline, maxDecline, money, position, startDate)
		}
		//if maxDecline < warningDecline {
		//	starFlag = false
		//}
	}

	//平均每年盈利至少10%
	//fromYear, err := strconv.Atoi(fundBasic.ListDate[:4])
	//if err != nil {
	//	return false, 0, 0, err
	//}
	//yearCount := 2023 - fromYear
	//if maxResult/float64(yearCount) < minTargetResult {
	//	starFlag = false
	//}

	return starFlag, maxDecline, maxResult, nil
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

func (f *Fund) getMA30(ctx context.Context, dailyList []*do.FundDaily, currentDay int) float64 {
	sum := float64(0)

	for j := 0; j < 30; j++ {
		sum = sum + dailyList[currentDay-j].Close
	}

	sum /= 30

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
