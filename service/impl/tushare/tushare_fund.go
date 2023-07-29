package tushare

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"my_stock_market/model/do"
	"my_stock_market/service/interface/tushare"
)

func (i Impl) FundBasic(ctx context.Context, param tushare.FundBasicParam) (*tushare.FundBasicResult, error) {
	params := map[string]interface{}{}

	if len(param.Market) > 0 {
		params["market"] = param.Market
	}
	if len(param.Status) > 0 {
		params["status"] = param.Status
	}
	if param.Limit > 0 {
		params["limit"] = param.Limit
	}
	if param.Offset > 0 {
		params["offset"] = param.Offset
	}

	fieldList := []string{"ts_code", "name", "management",
		"custodian", "fund_type", "found_date", "due_date", "list_date",
		"issue_date", "delist_date", "issue_amount", "m_fee", "c_fee", "duration_year",
		"p_value", "min_amount", "exp_return", "benchmark", "status", "invest_type", "type",
		"trustee", "purc_startdate", "redm_startdate", "market"}
	data, err := Post(ctx, FundBasicAPI, params, fieldList)
	if err != nil {
		return nil, err
	}

	fields, ok := data["fields"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion fields failed")
		logrus.Errorf("[FundBasic] force conversion fields failed")
		return nil, err
	}
	items, ok := data["items"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion items failed")
		logrus.Errorf("[FundBasic] force conversion items failed")
		return nil, err
	}

	fundBasicList := make([]*do.FundBasic, 0)
	for _, itemTemp := range items {
		item := itemTemp.([]interface{})
		fundBasic := &do.FundBasic{}
		if len(fields) != len(item) {
			err := fmt.Errorf("item.len is not equal fields.len")
			logrus.Errorf("[StockBasic] item.len is not equal fields.len")
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
			logrus.Errorf("[StockBasic] marshal temp error: %v", err)
			return nil, err
		}
		err = json.Unmarshal(tempJson, fundBasic)
		if err != nil {
			logrus.Errorf("[StockBasic] unmarshal tempJson error: %v", err)
			return nil, err
		}
		fundBasicList = append(fundBasicList, fundBasic)
	}

	return &tushare.FundBasicResult{FundBasicList: fundBasicList}, nil
}

func (i Impl) FundDaily(ctx context.Context, param tushare.FundDailyParam) (*tushare.FundDailyResult, error) {
	params := map[string]interface{}{}
	if len(param.TSCode) > 0 {
		params["ts_code"] = param.TSCode
	}
	if len(param.TradeDate) > 0 {
		params["trade_date"] = param.TradeDate
	}
	if len(param.StartDate) > 0 {
		params["start_date"] = param.StartDate
	}
	if len(param.EndDate) > 0 {
		params["end_date"] = param.EndDate
	}
	fieldList := []string{"ts_code", "trade_date", "open", "high",
		"low", "close", "pre_close", "change", "pct_chg",
		"vol", "amount"}
	data, err := Post(ctx, FundDailyAPI, params, fieldList)
	if err != nil {
		return nil, err
	}
	fields, ok := data["fields"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion fields failed")
		logrus.Errorf("[Daily] force conversion fields failed")
		return nil, err
	}
	items, ok := data["items"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion items failed")
		logrus.Errorf("[Daily] force conversion items failed")
		return nil, err
	}

	fundDailyList := make([]*do.FundDaily, 0)
	for _, itemTemp := range items {
		item := itemTemp.([]interface{})
		fundDaily := &do.FundDaily{}
		if len(fields) != len(item) {
			err := fmt.Errorf("item.len is not equal fields.len")
			logrus.Errorf("[Daily] item.len is not equal fields.len")
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
			logrus.Errorf("[Daily] marshal temp error: %v", err)
			return nil, err
		}
		err = json.Unmarshal(tempJson, fundDaily)
		if err != nil {
			logrus.Errorf("[Daily] unmarshal tempJson error: %v", err)
			return nil, err
		}
		fundDailyList = append(fundDailyList, fundDaily)
	}

	return &tushare.FundDailyResult{FundDailyList: fundDailyList}, nil
}
