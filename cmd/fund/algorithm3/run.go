package algorithm3

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/common/util"
	"my_stock_market/model/do"
	"my_stock_market/repo/interface/fund_basic"
	"my_stock_market/repo/interface/fund_daily"
	"sort"
	"strings"
)

func (a *Algorithm3) Run(ctx context.Context) error {
	fundList, err := a.GetTargetFund(ctx)
	if err != nil {
		return err
	}
	//chart := chart2.NewChart(ctx, "网格交易")

	for _, fundBasic := range fundList {
		err := a.Exec(ctx, fundBasic)
		if err != nil {
			return err
		}
	}

	return nil
}

type BuyTime struct {
	BuyDate  string
	SellDate string
}

// Exec
// 1. 从20170601开始回测，跌x点买入x元？
func (a *Algorithm3) Exec(ctx context.Context, fundBasic *do.FundBasic) error {
	fundDailyList, err := a.getAllFundDaily(ctx, fundBasic.TSCode, "20170601")
	if err != nil {
		return err
	}

	haveFund := false
	last := 0
	buyList := make([]*BuyTime, 0)
	currentPrice := float64(0)

	for index, fundDaily := range fundDailyList {
		if index == 0 {
			continue
		}
		maxPrice := float64(-10000)
		for j := last; j < index; j++ {
			if maxPrice < fundDailyList[j].Close {
				maxPrice = fundDailyList[j].Close
			}
		}
		if !haveFund && fundDaily.Close < 0.5*maxPrice {
			haveFund = true
			buyList = append(buyList, &BuyTime{
				BuyDate: fundDaily.TradeDate,
			})
			currentPrice = fundDaily.Close
		}

		if haveFund && fundDaily.Close > 1.3*currentPrice {
			haveFund = false
			buyList[len(buyList)-1].SellDate = fundDaily.TradeDate
			last = index
		}
	}

	logrus.Infof("tsCode: %v, 名称: %v, 记录: %v", fundBasic.TSCode, fundBasic.Name, util.ToJsonStr(buyList))

	return nil

}

func (a *Algorithm3) GetTargetFund(ctx context.Context) ([]*do.FundBasic, error) {
	fundBasicResult, err := a.FundBasicDAL.GetAllFundBasic(ctx, fund_basic.GetAllFundBasicParam{})
	if err != nil {
		return nil, err
	}

	// 待处理标的
	fundMap := map[string]*do.FundBasic{}

	targetTsCodeList := []string{
		//"512040",
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

func (a *Algorithm3) getAllFundDaily(ctx context.Context, tsCode, startDate string) ([]*do.FundDaily, error) {
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
