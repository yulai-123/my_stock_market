package tushare

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/repo/interface/index"
	"my_stock_market/service/interface/tushare"
)

// SaveIndexList 保存全量股票列表，已经存在的会进行更新
func (s *Stock) SaveAllIndex(ctx context.Context) error {
	indexBasicResult, err := s.TuShare.Index(ctx, tushare.IndexParam{})
	if err != nil {
		return err
	}

	logrus.Infof("拉取股票列表成功，长度：%v", len(indexBasicResult.IndexList))

	_, err = s.IndexDAL.BatchSaveIndex(ctx, index.BatchSaveIndexParam{
		IndexList: indexBasicResult.IndexList,
	})
	if err != nil {
		return err
	}

	logrus.Info("保存数据库成功")
	return nil
}
