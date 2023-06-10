package cmd

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/repo/interface/cashflow"
	"my_stock_market/service/interface/tushare"
)

func (s *Stock) SaveAllCashflow(ctx context.Context) error {
	stockBasicResult, err := s.TuShare.StockBasic(ctx, tushare.StockBasicParam{Limit: 100000})
	if err != nil {
		return err
	}

	logrus.Infof("拉取股票列表成功，长度：%v", len(stockBasicResult.StockList))

	tsCodeList := make([]string, 0)
	for _, stock := range stockBasicResult.StockList {
		tsCodeList = append(tsCodeList, stock.TSCode)
	}

	tsCodeList = []string{"689009.SH"}

	for _, tsCode := range tsCodeList {
		logrus.Infof("正在保存: %v", tsCode)
		result, err := s.TuShare.Cashflow(ctx, tushare.CashflowParam{TSCode: tsCode})
		if err != nil {
			logrus.Errorf("拉取 %v 现金流量表失败：%v", tsCode, err)
			return err
		}

		err = s.CashflowDAL.BatchSaveCashflow(ctx, cashflow.BatchSaveCashflowParam{CashflowList: result.CashflowList})
		if err != nil {
			logrus.Errorf("保存 %v 现金流量表失败：%v", tsCode, err)
			return err
		}
	}

	return nil
}
