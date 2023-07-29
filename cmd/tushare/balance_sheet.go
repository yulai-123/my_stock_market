package tushare

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/repo/interface/balance_sheet"
	"my_stock_market/service/interface/tushare"
	"time"
)

func (s *Stock) SaveAllBalanceSheet(ctx context.Context) error {
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
		logrus.Infof("正在保存: %v", tsCode)
		result, err := s.TuShare.BalanceSheet(ctx, tushare.BalanceSheetParam{TSCode: tsCode})
		if err != nil {
			logrus.Errorf("拉取 %v 资产债务表失败：%v", tsCode, err)
			continue
		}

		err = s.BalanceSheetDAL.BatchSaveBalanceSheet(ctx, balance_sheet.BatchSaveBalanceSheetParam{BalanceSheetList: result.BalanceSheetList})
		if err != nil {
			logrus.Errorf("保存 %v 资产债务表失败：%v", tsCode, err)
			continue
		}
		time.Sleep(500 * time.Microsecond)
	}

	return nil
}
