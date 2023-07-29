package tushare

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/repo/interface/stock"
	"my_stock_market/service/interface/tushare"
)

// SaveStockList 保存全量股票列表，已经存在的会进行更新
func (s *Stock) SaveStockList(ctx context.Context) error {
	stockBasicResult, err := s.TuShare.StockBasic(ctx, tushare.StockBasicParam{Limit: 100000})
	if err != nil {
		return err
	}

	logrus.Infof("拉取股票列表成功，长度：%v", len(stockBasicResult.StockList))

	_, err = s.StockDAL.BatchSaveStock(ctx, stock.BatchSaveStockParam{
		StockList: stockBasicResult.StockList,
	})
	if err != nil {
		return err
	}

	logrus.Info("保存数据库成功")
	return nil
}

func (s *Stock) TestStockList(ctx context.Context) error {
	stockBasicResult, err := s.TuShare.StockBasic(ctx, tushare.StockBasicParam{
		//ListStatus: "L",
		Limit: 10000,
	})
	if err != nil {
		return err
	}

	//logrus.Info(util.ToJsonStr(stockBasicResult.StockList))
	logrus.Info(len(stockBasicResult.StockList))

	return nil
}
