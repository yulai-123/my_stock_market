package tushare

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"my_stock_market/model/do"
	"my_stock_market/service/interface/tushare"
)

type Impl struct {
}

func NewTuShare(ctx context.Context) tushare.TuShare {
	return Impl{}
}

func (i Impl) TradeCal(ctx context.Context, param tushare.TradeCalParam) (*tushare.TradeCalResult, error) {
	params := map[string]interface{}{}
	if len(param.Exchange) > 0 {
		params["exchange"] = param.Exchange
	}
	if len(param.IsOpen) > 0 {
		params["is_open"] = param.IsOpen
	}
	if len(param.EndDate) > 0 {
		params["end_date"] = param.EndDate
	}
	if len(param.StartDate) > 0 {
		params["start_date"] = param.StartDate
	}
	fieldList := []string{"exchange", "cal_date", "is_open", "pretrade_date"}
	data, err := Post(ctx, TradeCalAPI, params, fieldList)
	if err != nil {
		return nil, err
	}

	fields, ok := data["fields"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion fields failed")
		logrus.Errorf("[StockBasic] force conversion fields failed")
		return nil, err
	}
	items, ok := data["items"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion items failed")
		logrus.Errorf("[StockBasic] force conversion items failed")
		return nil, err
	}

	tradeCalList := make([]*tushare.SingleTradeCal, 0)
	for _, itemTemp := range items {
		item := itemTemp.([]interface{})
		tradeCal := &tushare.SingleTradeCal{}
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
		err = json.Unmarshal(tempJson, tradeCal)
		if err != nil {
			logrus.Errorf("[StockBasic] unmarshal tempJson error: %v", err)
			return nil, err
		}
		tradeCalList = append(tradeCalList, tradeCal)
	}

	return &tushare.TradeCalResult{TradeCalList: tradeCalList}, nil
}

func (i Impl) StockBasic(ctx context.Context, param tushare.StockBasicParam) (*tushare.StockBasicResult, error) {
	params := map[string]interface{}{}
	if param.Limit > 0 {
		params["limit"] = param.Limit
	}
	if param.Offset > 0 {
		params["offset"] = param.Offset
	}
	if len(param.IsHS) > 0 {
		params["is_hs"] = param.IsHS
	}
	if len(param.Exchange) > 0 {
		params["exchange"] = param.Exchange
	}
	if len(param.ListStatus) > 0 {
		params["list_status"] = param.ListStatus
	}
	if len(param.Name) > 0 {
		params["name"] = param.Name
	}
	if len(param.Market) > 0 {
		params["market"] = param.Market
	}
	if len(param.TSCode) > 0 {
		params["ts_code"] = param.TSCode
	}
	fieldList := []string{"ts_code", "symbol", "name", "area",
		"industry", "fullname", "enname", "cnspell", "market",
		"exchange", "curr_type", "list_status", "list_date",
		"delist_date", "is_hs"}
	data, err := Post(ctx, StockBasicAPI, params, fieldList)
	if err != nil {
		return nil, err
	}

	fields, ok := data["fields"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion fields failed")
		logrus.Errorf("[StockBasic] force conversion fields failed")
		return nil, err
	}
	items, ok := data["items"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion items failed")
		logrus.Errorf("[StockBasic] force conversion items failed")
		return nil, err
	}

	stockList := make([]*do.Stock, 0)
	for _, itemTemp := range items {
		item := itemTemp.([]interface{})
		stock := &do.Stock{}
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
		err = json.Unmarshal(tempJson, stock)
		if err != nil {
			logrus.Errorf("[StockBasic] unmarshal tempJson error: %v", err)
			return nil, err
		}
		stockList = append(stockList, stock)
	}

	return &tushare.StockBasicResult{StockList: stockList}, nil
}

func (i Impl) Daily(ctx context.Context, param tushare.DailyParam) (*tushare.DailyResult, error) {
	params := map[string]interface{}{}
	if len(param.TSCode) > 0 {
		temp := param.TSCode[0]
		for i := 1; i < len(param.TSCode); i++ {
			temp = temp + "," + param.TSCode[i]
		}
		params["ts_code"] = temp
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
	data, err := Post(ctx, StockDailyAPI, params, fieldList)
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

	stockDailyList := make([]*do.StockDaily, 0)
	for _, itemTemp := range items {
		item := itemTemp.([]interface{})
		stockDaily := &do.StockDaily{}
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
		err = json.Unmarshal(tempJson, stockDaily)
		if err != nil {
			logrus.Errorf("[Daily] unmarshal tempJson error: %v", err)
			return nil, err
		}
		stockDailyList = append(stockDailyList, stockDaily)
	}

	return &tushare.DailyResult{StockDailyList: stockDailyList}, nil
}

func (i Impl) Weekly(ctx context.Context, param tushare.WeeklyParam) (*tushare.WeeklyResult, error) {
	params := map[string]interface{}{}
	if len(param.TSCode) > 0 {
		temp := param.TSCode[0]
		for i := 1; i < len(param.TSCode); i++ {
			temp = temp + "," + param.TSCode[i]
		}
		params["ts_code"] = temp
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
	data, err := Post(ctx, StockWeeklyAPI, params, fieldList)
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

	stockWeeklyList := make([]*do.StockWeekly, 0)
	for _, itemTemp := range items {
		item := itemTemp.([]interface{})
		stockWeekly := &do.StockWeekly{}
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
		err = json.Unmarshal(tempJson, stockWeekly)
		if err != nil {
			logrus.Errorf("[Daily] unmarshal tempJson error: %v", err)
			return nil, err
		}
		// 周线和月线的pct_chg都不是百分制的，这里乘以100换算成百分制
		stockWeekly.PctChg *= 100
		stockWeeklyList = append(stockWeeklyList, stockWeekly)
	}

	return &tushare.WeeklyResult{StockWeeklyList: stockWeeklyList}, nil
}

func (i Impl) Monthly(ctx context.Context, param tushare.MonthlyParam) (*tushare.MonthlyResult, error) {
	params := map[string]interface{}{}
	if len(param.TSCode) > 0 {
		temp := param.TSCode[0]
		for i := 1; i < len(param.TSCode); i++ {
			temp = temp + "," + param.TSCode[i]
		}
		params["ts_code"] = temp
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
	data, err := Post(ctx, StockMonthlyAPI, params, fieldList)
	if err != nil {
		return nil, err
	}
	fields, ok := data["fields"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion fields failed")
		logrus.Errorf("[Monthly] force conversion fields failed")
		return nil, err
	}
	items, ok := data["items"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion items failed")
		logrus.Errorf("[Monthly] force conversion items failed")
		return nil, err
	}

	stockMonthlyList := make([]*do.StockMonthly, 0)
	for _, itemTemp := range items {
		item := itemTemp.([]interface{})
		stockMonthly := &do.StockMonthly{}
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
		err = json.Unmarshal(tempJson, stockMonthly)
		if err != nil {
			logrus.Errorf("[Daily] unmarshal tempJson error: %v", err)
			return nil, err
		}
		// 周线和月线的pct_chg都不是百分制的，这里乘以100换算成百分制
		stockMonthly.PctChg *= 100
		stockMonthlyList = append(stockMonthlyList, stockMonthly)
	}

	return &tushare.MonthlyResult{StockMonthlyList: stockMonthlyList}, nil
}

func (i Impl) DailyBasic(ctx context.Context, param tushare.DailyBasicParam) (*tushare.DailyBasicResult, error) {
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
	fieldList := []string{"ts_code", "trade_date", "close", "turnover_rate", "turnover_rate_f", "volume_ratio",
		"pe", "pe_ttm", "pb", "ps", "ps_ttm", "dv_ratio", "dv_ttm",
		"total_share", "float_share", "free_share", "total_mv", "circ_mv"}
	data, err := Post(ctx, DailyBasicAPI, params, fieldList)
	if err != nil {
		return nil, err
	}
	fields, ok := data["fields"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion fields failed")
		logrus.Errorf("[Monthly] force conversion fields failed")
		return nil, err
	}
	items, ok := data["items"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion items failed")
		logrus.Errorf("[Monthly] force conversion items failed")
		return nil, err
	}

	dailyBasicList := make([]*do.DailyBasic, 0)
	for _, itemTemp := range items {
		item := itemTemp.([]interface{})
		dailyBasic := &do.DailyBasic{}
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
		err = json.Unmarshal(tempJson, dailyBasic)
		if err != nil {
			logrus.Errorf("[Daily] unmarshal tempJson error: %v", err)
			return nil, err
		}

		dailyBasicList = append(dailyBasicList, dailyBasic)
	}

	return &tushare.DailyBasicResult{DailyBasicList: dailyBasicList}, nil
}

func (i Impl) Index(ctx context.Context, param tushare.IndexParam) (*tushare.IndexResult, error) {
	params := map[string]interface{}{}

	if len(param.TSCode) > 0 {
		params["ts_code"] = param.TSCode
	}

	fieldList := []string{"ts_code", "name",
		"fullname", "market", "list_date",
		"publisher", "index_type", "category", "base_date",
		"base_point", "weight_rule", "desc", "exp_date"}
	data, err := Post(ctx, IndexAPI, params, fieldList)
	if err != nil {
		return nil, err
	}

	fields, ok := data["fields"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion fields failed")
		logrus.Errorf("[Index] force conversion fields failed")
		return nil, err
	}
	items, ok := data["items"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion items failed")
		logrus.Errorf("[Index] force conversion items failed")
		return nil, err
	}

	indexList := make([]*do.Index, 0)
	for _, itemTemp := range items {
		item := itemTemp.([]interface{})
		index := &do.Index{}
		if len(fields) != len(item) {
			err := fmt.Errorf("item.len is not equal fields.len")
			logrus.Errorf("[Index] item.len is not equal fields.len")
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
			logrus.Errorf("[Index] marshal temp error: %v", err)
			return nil, err
		}
		err = json.Unmarshal(tempJson, index)
		if err != nil {
			logrus.Errorf("[Index] unmarshal tempJson error: %v", err)
			return nil, err
		}
		indexList = append(indexList, index)
	}

	return &tushare.IndexResult{IndexList: indexList}, nil
}

func (i Impl) IndexDaily(ctx context.Context, param tushare.IndexDailyParam) (*tushare.IndexDailyResult, error) {
	params := map[string]interface{}{}

	params["ts_code"] = param.TSCode

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
	data, err := Post(ctx, IndexDailyAPI, params, fieldList)
	if err != nil {
		return nil, err
	}
	fields, ok := data["fields"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion fields failed")
		logrus.Errorf("[IndexDaily] force conversion fields failed")
		return nil, err
	}
	items, ok := data["items"].([]interface{})
	if !ok {
		err := fmt.Errorf("force conversion items failed")
		logrus.Errorf("[IndexDaily] force conversion items failed")
		return nil, err
	}

	indexDailyList := make([]*do.IndexDaily, 0)
	for _, itemTemp := range items {
		item := itemTemp.([]interface{})
		indexDaily := &do.IndexDaily{}
		if len(fields) != len(item) {
			err := fmt.Errorf("item.len is not equal fields.len")
			logrus.Errorf("[IndexDaily] item.len is not equal fields.len")
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
			logrus.Errorf("[IndexDaily] marshal temp error: %v", err)
			return nil, err
		}
		err = json.Unmarshal(tempJson, indexDaily)
		if err != nil {
			logrus.Errorf("[IndexDaily] unmarshal tempJson error: %v", err)
			return nil, err
		}
		indexDailyList = append(indexDailyList, indexDaily)
	}

	return &tushare.IndexDailyResult{IndexDailyList: indexDailyList}, nil
}
