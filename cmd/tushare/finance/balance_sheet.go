package finance

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/repo/interface/balance_sheet"
	"my_stock_market/service/interface/tushare"
	"time"
)

// SaveAllBalanceSheet 保存所有资产债务表
func (f *Finance) SaveAllBalanceSheet(ctx context.Context) error {
	logrus.Infof("保存资产负债表")
	stockBasicResult, err := f.TuShare.StockBasic(ctx, tushare.StockBasicParam{Limit: 100000})
	if err != nil {
		return err
	}

	logrus.Infof("拉取股票列表成功，长度：%v", len(stockBasicResult.StockList))

	tsCodeList := make([]string, 0)
	for _, stock := range stockBasicResult.StockList {
		tsCodeList = append(tsCodeList, stock.TSCode)
	}

	for i, tsCode := range tsCodeList {
		logrus.Infof("正在保存: %v, 进度: %v/%v", tsCode, i, len(tsCodeList))
		result, err := f.TuShare.BalanceSheet(ctx, tushare.BalanceSheetParam{TSCode: tsCode})
		if err != nil {
			logrus.Errorf("拉取 %v 资产债务表失败：%v", tsCode, err)
			continue
		}
		logrus.Infof("tsCode: %v, 长度: %v", tsCode, len(result.BalanceSheetList))

		err = f.BalanceSheetDAL.BatchSaveBalanceSheet(ctx, balance_sheet.BatchSaveBalanceSheetParam{BalanceSheetList: result.BalanceSheetList})
		if err != nil {
			logrus.Errorf("保存 %v 资产债务表失败：%v", tsCode, err)
			continue
		}
		time.Sleep(time.Second)
	}

	return nil
}
