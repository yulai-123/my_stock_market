package tushare

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"my_stock_market/model/do"
	"my_stock_market/service/interface/tushare"
	"reflect"
)

func (i Impl) Cashflow(ctx context.Context, param tushare.CashflowParam) (*tushare.CashflowResult, error) {
	params := map[string]interface{}{}
	params["ts_code"] = param.TSCode
	if len(param.Period) > 0 {
		params["period"] = param.Period
	}

	fieldList := []string{"ts_code", "ann_date", "f_ann_date", "end_date",
		"report_type", "comp_type", "end_type", "n_cashflow_act"}
	data, err := Post(ctx, CashflowAPI, params, fieldList)
	if err != nil {
		return nil, err
	}

	fields, ok := data["fields"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion fields failed")
		logrus.Errorf("[Cashflow] force conversion fields failed")
		return nil, err
	}
	items, ok := data["items"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion items failed")
		logrus.Errorf("[Cashflow] force conversion items failed")
		return nil, err
	}

	endDateMap := make(map[string]bool)
	cashflowList := make([]*do.Cashflow, 0)
	for _, itemTemp := range items {
		item := itemTemp.([]interface{})
		cashflow := &do.Cashflow{}
		if len(fields) != len(item) {
			err := fmt.Errorf("item.len is not equal fields.len")
			logrus.Errorf("[Cashflow] item.len is not equal fields.len")
			return nil, err
		}
		temp := make(map[string]interface{})
		for i, value := range item {
			if value == nil {
				continue
			}
			temp[fields[i].(string)] = value
		}
		tempJson, err := json.Marshal(temp)
		if err != nil {
			logrus.Errorf("[Cashflow] marshal temp error: %v", err)
			return nil, err
		}
		err = json.Unmarshal(tempJson, cashflow)
		if err != nil {
			logrus.Errorf("[Cashflow] unmarshal tempJson error: %v", err)
			return nil, err
		}
		if _, ok := endDateMap[cashflow.EndDate]; ok {
			logrus.Infof("duplicate end_date: %s, ts_code: %v", cashflow.EndDate, cashflow.TSCode)
		}
		endDateMap[cashflow.EndDate] = true
		cashflowList = append(cashflowList, cashflow)
	}

	return &tushare.CashflowResult{CashflowList: cashflowList}, nil
}

func (i Impl) BalanceSheet(ctx context.Context, param tushare.BalanceSheetParam) (*tushare.BalanceSheetResult, error) {
	params := map[string]interface{}{}
	params["ts_code"] = param.TSCode
	if len(param.Period) > 0 {
		params["period"] = param.Period
	}

	fieldList := []string{"ts_code", "ann_date", "f_ann_date", "end_date",
		"report_type", "comp_type", "end_type", "total_assets", "total_liab"}
	data, err := Post(ctx, BalanceSheetAPI, params, fieldList)
	if err != nil {
		return nil, err
	}

	fields, ok := data["fields"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion fields failed")
		logrus.Errorf("[BalanceSheetDAL] force conversion fields failed")
		return nil, err
	}
	items, ok := data["items"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion items failed")
		logrus.Errorf("[BalanceSheetDAL] force conversion items failed, current: %v, target: %v", reflect.TypeOf(data["items"]), "[]interface{}")
		return nil, err
	}

	balanceSheetList := make([]*do.BalanceSheet, 0)
	for _, itemTemp := range items {
		item := itemTemp.([]interface{})
		balanceSheet := &do.BalanceSheet{}
		if len(fields) != len(item) {
			err := fmt.Errorf("item.len is not equal fields.len")
			logrus.Errorf("[BalanceSheetDAL] item.len is not equal fields.len")
			return nil, err
		}
		temp := make(map[string]interface{})
		for i, value := range item {
			if value == nil {
				continue
			}
			temp[fields[i].(string)] = value
		}
		tempJson, err := json.Marshal(temp)
		if err != nil {
			logrus.Errorf("[BalanceSheetDAL] marshal temp error: %v", err)
			return nil, err
		}
		err = json.Unmarshal(tempJson, balanceSheet)
		if err != nil {
			logrus.Errorf("[BalanceSheetDAL] unmarshal tempJson error: %v", err)
			return nil, err
		}
		balanceSheetList = append(balanceSheetList, balanceSheet)
	}

	return &tushare.BalanceSheetResult{BalanceSheetList: balanceSheetList}, nil
}

func (i Impl) Income(ctx context.Context, param tushare.IncomeParam) (*tushare.IncomeResult, error) {
	params := map[string]interface{}{}
	params["ts_code"] = param.TSCode
	if len(param.Period) > 0 {
		params["period"] = param.Period
	}

	fieldList := []string{"ts_code", "ann_date", "f_ann_date", "end_date",
		"report_type", "comp_type", "end_type", "revenue", "n_income_attr_p", "oper_cost"}
	data, err := Post(ctx, IncomeAPI, params, fieldList)
	if err != nil {
		return nil, err
	}

	fields, ok := data["fields"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion fields failed")
		logrus.Errorf("[Income] force conversion fields failed")
		return nil, err
	}
	items, ok := data["items"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion items failed")
		logrus.Errorf("[Income] force conversion items failed")
		return nil, err
	}

	incomeList := make([]*do.Income, 0)
	for _, itemTemp := range items {
		item := itemTemp.([]interface{})
		income := &do.Income{}
		if len(fields) != len(item) {
			err := fmt.Errorf("item.len is not equal fields.len")
			logrus.Errorf("[Income] item.len is not equal fields.len")
			return nil, err
		}
		temp := make(map[string]interface{})
		for i, value := range item {
			if value == nil {
				continue
			}
			temp[fields[i].(string)] = value
		}
		tempJson, err := json.Marshal(temp)
		if err != nil {
			logrus.Errorf("[Income] marshal temp error: %v", err)
			return nil, err
		}
		err = json.Unmarshal(tempJson, income)
		if err != nil {
			logrus.Errorf("[Income] unmarshal tempJson error: %v", err)
			return nil, err
		}
		incomeList = append(incomeList, income)
	}

	return &tushare.IncomeResult{IncomeList: incomeList}, nil
}
