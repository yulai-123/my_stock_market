package fund_basic

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"my_stock_market/infra/mysql"
	do2 "my_stock_market/model/do"
	"my_stock_market/model/po"
	fund "my_stock_market/repo/interface/fund_basic"
	"time"
)

type Impl struct {
	provider *mysql.DBProvider
}

func NewFundBasicDAL(ctx context.Context) fund.DAL {
	return &Impl{
		provider: mysql.GetDBProvider(ctx),
	}
}

func (i *Impl) BatchSaveFundBasic(ctx context.Context, param fund.BatchSaveFundBasicParam) (*fund.BatchSaveFundBasicResult, error) {
	fund_basicPOList := make([]*po.FundBasic, 0)
	createdAt := int64(time.Now().Unix())
	updatedAt := createdAt
	for _, do := range param.FundBasicList {
		po, err := do2.FundBasicDO2PO(ctx, do)
		if err != nil {
			return nil, err
		}
		po.CreatedAt = createdAt
		po.UpdatedAt = updatedAt
		fund_basicPOList = append(fund_basicPOList, po)
	}

	fieldUpdateColumns := []string{"name", "management",
		"custodian", "fund_type", "found_date", "due_date", "list_date",
		"issue_date", "delist_date", "issue_amount", "m_fee", "c_fee", "duration_year",
		"p_value", "min_amount", "exp_return", "benchmark", "status", "invest_type", "type",
		"trustee", "purc_startdate", "redm_startdate", "market", "updated_at"}
	err := i.provider.WithContext(ctx).Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns(fieldUpdateColumns)}).
		CreateInBatches(fund_basicPOList, 200).Error
	if err != nil {
		logrus.Errorf("[BatchSaveFundBasic] create in batches error: %v", err)
		return nil, err
	}

	return &fund.BatchSaveFundBasicResult{}, nil
}

func (i *Impl) GetAllFundBasic(ctx context.Context, param fund.GetAllFundBasicParam) (*fund.GetAllFundBasicResult, error) {
	result := make([]*do2.FundBasic, 0)

	fundBasicPOList := make([]*po.FundBasic, 0)
	err := i.provider.WithContext(ctx).Find(&fundBasicPOList).Error
	if err != nil {
		logrus.Errorf("[GetAllFundBasic] find data error: %v", err)
		return nil, err
	}
	for _, po := range fundBasicPOList {
		result = append(result, do2.FundBasicPO2DO(ctx, po))
	}

	return &fund.GetAllFundBasicResult{FundBasicList: result}, nil
}
