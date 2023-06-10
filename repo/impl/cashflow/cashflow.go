package cashflow

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"my_stock_market/infra/mysql"
	do2 "my_stock_market/model/do"
	"my_stock_market/model/po"
	"my_stock_market/repo/interface/cashflow"
	"time"
)

type Impl struct {
	provider *mysql.DBProvider
}

func NewCashflowDAL(ctx context.Context) cashflow.Cashflow {
	return &Impl{provider: mysql.GetDBProvider(ctx)}
}

func (i Impl) BatchSaveCashflow(ctx context.Context, param cashflow.BatchSaveCashflowParam) error {
	var err error
	for t := 0; t < 60; t++ {
		cashflowPOList := make([]*po.Cashflow, 0)
		createdAt := time.Now().Unix()
		updateAt := createdAt
		for _, do := range param.CashflowList {
			po := do2.CashflowDO2PO(ctx, do)
			po.CreatedAt = createdAt
			po.UpdatedAt = updateAt
			cashflowPOList = append(cashflowPOList, po)
		}
		logrus.Info(len(cashflowPOList))

		fieldUpdateColumns := []string{"ann_date", "f_ann_date", "report_type",
			"comp_type", "end_type", "updated_at", "n_cashflow_act"}
		err = i.provider.WithContext(ctx).Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns(fieldUpdateColumns)}).
			CreateInBatches(cashflowPOList, 200).Error
		if err != nil {
			logrus.Errorf("[BatchSaveCashflow] create in batches error: %v, try: %v", err, t)
			time.Sleep(1 * time.Second)
			continue
		}

		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func (i Impl) BatchGetCashflow(ctx context.Context, param cashflow.BatchGetCashflowParam) (*cashflow.BatchGetCashflowResult, error) {
	result := make(map[string][]*do2.Cashflow)

	incomePOList := make([]*po.Cashflow, 0)
	db := i.provider.WithContext(ctx).Where("ts_code in (?)", param.StockCodeList)
	err := db.Find(&incomePOList).Error
	if err != nil {
		logrus.Errorf("[BatchGetCashflow] find data error: %v", err)
		return nil, err
	}

	for _, po := range incomePOList {
		do := do2.CashflowPO2DO(ctx, po)
		result[do.TSCode] = append(result[do.TSCode], do)
	}

	return &cashflow.BatchGetCashflowResult{CashflowMap: result}, nil
}
