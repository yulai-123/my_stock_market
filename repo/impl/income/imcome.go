package income

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"my_stock_market/infra/mysql"
	do2 "my_stock_market/model/do"
	"my_stock_market/model/po"
	"my_stock_market/repo/interface/income"
	"time"
)

type Impl struct {
	provider *mysql.DBProvider
}

func NewIncomeDAL(ctx context.Context) income.Income {
	return &Impl{provider: mysql.GetDBProvider(ctx)}
}

func (i Impl) BatchSaveIncome(ctx context.Context, param income.BatchSaveIncomeParam) error {
	var err error
	for t := 0; t < 10; t++ {
		incomePOList := make([]*po.Income, 0)
		createdAt := time.Now().Unix()
		updateAt := createdAt
		for _, do := range param.IncomeList {
			po := do2.IncomeDO2PO(ctx, do)
			po.CreatedAt = createdAt
			po.UpdatedAt = updateAt
			incomePOList = append(incomePOList, po)
		}

		fieldUpdateColumns := []string{"ann_date", "f_ann_date", "report_type",
			"comp_type", "end_type", "updated_at", "revenue", "n_income_attr_p", "oper_cost"}
		err = i.provider.WithContext(ctx).Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns(fieldUpdateColumns)}).
			CreateInBatches(incomePOList, 200).Error
		if err != nil {
			logrus.Errorf("[BatchSaveIncome] create in batches error: %v, try: %v", err, t)
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

func (i Impl) BatchGetIncome(ctx context.Context, param income.BatchGetIncomeParam) (*income.BatchGetIncomeResult, error) {
	result := make(map[string][]*do2.Income)

	incomePOList := make([]*po.Income, 0)
	db := i.provider.WithContext(ctx).Where("ts_code in (?)", param.StockCodeList)
	err := db.Find(&incomePOList).Error
	if err != nil {
		logrus.Errorf("[BatchGetIncome] find data error: %v", err)
		return nil, err
	}

	for _, po := range incomePOList {
		do := do2.IncomePO2DO(ctx, po)
		result[do.TSCode] = append(result[do.TSCode], do)
	}

	return &income.BatchGetIncomeResult{IncomeMap: result}, nil
}
