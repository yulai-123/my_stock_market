package index

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"my_stock_market/infra/mysql"
	do2 "my_stock_market/model/do"
	"my_stock_market/model/po"
	"my_stock_market/repo/interface/index"
	"time"
)

type Impl struct {
	provider *mysql.DBProvider
}

func NewIndexDAL(ctx context.Context) index.DAL {
	return &Impl{
		provider: mysql.GetDBProvider(ctx),
	}
}

func (i *Impl) BatchSaveIndex(ctx context.Context, param index.BatchSaveIndexParam) (*index.BatchSaveIndexResult, error) {
	indexPOList := make([]*po.Index, 0)
	createdAt := int64(time.Now().Unix())
	updatedAt := createdAt
	for _, do := range param.IndexList {
		po, err := do2.IndexDO2PO(ctx, do)
		if err != nil {
			return nil, err
		}
		po.CreatedAt = createdAt
		po.UpdatedAt = updatedAt
		indexPOList = append(indexPOList, po)
	}

	fieldUpdateColumns := []string{"name", "fullname",
		"market", "list_date", "publisher", "index_type", "category",
		"base_date", "base_point", "weight_rule", "desc", "exp_date", "updated_at"}
	err := i.provider.WithContext(ctx).Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns(fieldUpdateColumns)}).
		CreateInBatches(indexPOList, 200).Error
	if err != nil {
		logrus.Errorf("[BatchSaveIndex] create in batches error: %v", err)
		return nil, err
	}

	return &index.BatchSaveIndexResult{}, nil
}

func (i *Impl) GetAllIndex(ctx context.Context, param index.GetAllIndexParam) (*index.GetAllIndexResult, error) {
	result := make([]*do2.Index, 0)

	indexPOList := make([]*po.Index, 0)
	err := i.provider.WithContext(ctx).Find(&indexPOList).Error
	if err != nil {
		logrus.Errorf("[GetAllIndex] find data error: %v", err)
		return nil, err
	}
	for _, po := range indexPOList {
		result = append(result, do2.IndexPO2DO(ctx, po))
	}

	return &index.GetAllIndexResult{IndexList: result}, nil
}
