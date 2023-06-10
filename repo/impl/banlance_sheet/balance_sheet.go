package banlance_sheet

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
	"my_stock_market/infra/mysql"
	do2 "my_stock_market/model/do"
	"my_stock_market/model/po"
	"my_stock_market/repo/interface/balance_sheet"
	"time"
)

type Impl struct {
	provider *mysql.DBProvider
}

func NewBalanceSheetDAL(ctx context.Context) balance_sheet.BalanceSheet {
	return &Impl{provider: mysql.GetDBProvider(ctx)}
}

func (i Impl) BatchSaveBalanceSheet(ctx context.Context, param balance_sheet.BatchSaveBalanceSheetParam) error {
	var err error
	for t := 0; t < 10; t++ {
		balanceSheetPOList := make([]*po.BalanceSheet, 0)
		createdAt := time.Now().Unix()
		updateAt := createdAt
		for _, do := range param.BalanceSheetList {
			po := do2.BalanceSheetDo2PO(ctx, do)
			po.CreatedAt = createdAt
			po.UpdatedAt = updateAt
			balanceSheetPOList = append(balanceSheetPOList, po)
		}
		logrus.Info(len(balanceSheetPOList))

		fieldUpdateColumns := []string{"ann_date", "f_ann_date", "report_type", "comp_type",
			"end_type", "updated_at", "total_assets", "total_liab"}
		err = i.provider.WithContext(ctx).Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns(fieldUpdateColumns)}).
			CreateInBatches(balanceSheetPOList, 200).Error
		if err != nil {
			logrus.Errorf("[BatchSaveBalanceSheet] create in batches error: %v, try: %v", err, t)
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

func (i Impl) BatchGetBalanceSheet(ctx context.Context, param balance_sheet.BatchGetBalanceSheetParam) (*balance_sheet.BatchGetBalanceSheetResult, error) {
	result := make(map[string][]*do2.BalanceSheet)

	balanceSheetList := make([]*po.BalanceSheet, 0)
	db := i.provider.WithContext(ctx).Where("ts_code in (?)", param.StockCodeList)
	err := db.Find(&balanceSheetList).Error
	if err != nil {
		logrus.Errorf("[BatchGetBalanceSheet] find data error: %v", err)
		return nil, err
	}

	for _, po := range balanceSheetList {
		do := do2.BalanceSheetPO2DO(ctx, po)
		result[do.TSCode] = append(result[do.TSCode], do)
	}

	return &balance_sheet.BatchGetBalanceSheetResult{BalanceSheetMap: result}, nil
}
