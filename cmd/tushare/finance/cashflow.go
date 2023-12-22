package finance

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/repo/interface/cashflow"
	"my_stock_market/service/interface/tushare"
	"time"
)

// SaveAllCashflow 保存现金流量表
func (f *Finance) SaveAllCashflow(ctx context.Context) error {
	logrus.Infof("保存现金流量表")

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
		result, err := f.TuShare.Cashflow(ctx, tushare.CashflowParam{TSCode: tsCode})
		if err != nil {
			logrus.Errorf("拉取 %v 现金流量表失败：%v", tsCode, err)
			return err
		}

		logrus.Infof("tsCode: %v, 长度: %v", tsCode, len(result.CashflowList))

		err = f.CashflowDAL.BatchSaveCashflow(ctx, cashflow.BatchSaveCashflowParam{CashflowList: result.CashflowList})
		if err != nil {
			logrus.Errorf("保存 %v 现金流量表失败：%v", tsCode, err)
			return err
		}
		time.Sleep(time.Second)
	}

	return nil
}
