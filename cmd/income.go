package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/repo/interface/income"
	"my_stock_market/service/interface/tushare"
	"time"
)

// SaveAllIncome 保存所有利润表
func (s *Stock) SaveAllIncome(ctx context.Context) error {
	stockBasicResult, err := s.TuShare.StockBasic(ctx, tushare.StockBasicParam{Limit: 100000})
	if err != nil {
		return err
	}

	logrus.Infof("拉取股票列表成功，长度：%v", len(stockBasicResult.StockList))

	tsCodeList := make([]string, 0)
	for _, stock := range stockBasicResult.StockList {
		tsCodeList = append(tsCodeList, stock.TSCode)
	}

	for _, tsCode := range tsCodeList {
		result, err := s.TuShare.Income(ctx, tushare.IncomeParam{TSCode: tsCode})
		if err != nil {
			logrus.Errorf("拉取 %v 利润表失败：%v", tsCode, err)
			return err
		}

		err = s.IncomeDAL.BatchSaveIncome(ctx, income.BatchSaveIncomeParam{IncomeList: result.IncomeList})
		if err != nil {
			logrus.Errorf("保存 %v 利润表失败：%v", tsCode, err)
			return err
		}
		time.Sleep(time.Second)
	}

	return nil
}
