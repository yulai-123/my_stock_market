package finance

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"my_stock_market/common/excel"
	"my_stock_market/model/do"
	"my_stock_market/repo/interface/balance_sheet"
	"my_stock_market/repo/interface/cashflow"
	"my_stock_market/repo/interface/daily_basic"
	"my_stock_market/repo/interface/income"
	"my_stock_market/repo/interface/stock"
	"sort"
	"strconv"
)

// MakeFinancialStatements 制作一份财务报表
// 包含：
// 1. 企业名称
// 2. tsCode
// 3. 年份-月份
// 4. PE 市盈率
// 5. PB 市净率
// 6. PS 市销率
// 7. Revene 营业收入
// 8. NIncomeAttrP 净利润
// 9. NCashflowAct 现金流量净额
// 10. TotalAssets 资产总额
// 11. TotalLiab 负债总额
// 12. OperCost 业务成本
// 13. NetProfitMargin 净利润率
// 14. TotalAssetTurnover 总资产周转率
// 15. AssetLiabilityRatio 资产负债率
// 16. EquityMultiplier 权益乘数
// 17. ROE
// 18. GrossMargin 毛利率
// 19. ReveneGrowth 营业增长
// 20. NIncomeflowActGrowth 净利润增长
// 21. NCashflowActGrowth 资产总额增长
// 22. TotalAssetsGrowth 资产总额增长
// 23. TotalLiabGrowth 负债总额增长
// 24. OperCostGrowth 业务成本增长
func (f *Finance) MakeFinancialStatements(ctx context.Context) error {
	logrus.Infof("制作财务分析报表")

	// 设置需要查询企业的名单
	tsCodeList := make([]string, 0)
	tsCodeList = []string{
		"600004.SH", // 白云机场
		"605337.SH", // 李子园
		"600016.SH", // 民生银行
		"000651.SZ", // 格力
		"600000.SH", // 浦发
	}

	logrus.Infof("获取到企业名单，长度: %v", len(tsCodeList))

	// 获取企业的财务数据
	fullFinancialData, err := f.getFullFinancialData(ctx, tsCodeList)
	if err != nil {
		return err
	}
	logrus.Infof("获取到企业的财务数据，长度: %v", len(fullFinancialData))

	// 保存为excel
	err = f.saveAsExcel(ctx, fullFinancialData)
	if err != nil {
		return err
	}
	logrus.Infof("成功保存为excel")

	return nil
}

func (f *Finance) getFullFinancialData(ctx context.Context, tsCodeList []string) (map[string][]*EnterpriseGrowthData, error) {
	stockBasicResult, err := f.StockDAL.GetAllStock(ctx, stock.GetAllStockParam{})
	if err != nil {
		return nil, err
	}
	stockBasicMap := make(map[string]*do.Stock)
	for _, stock := range stockBasicResult.StockList {
		stockBasicMap[stock.TSCode] = stock
	}

	financialDataMap := make(map[string][]*EnterpriseGrowthData)

	for _, tsCode := range tsCodeList {
		stock, exist := stockBasicMap[tsCode]
		if !exist {
			logrus.Warnf("不存在企业, tsCode: %v", tsCode)
			continue
		}

		tsCode := stock.TSCode
		name := stock.Name

		cashflowMap, err := f.getCashflow(ctx, tsCode)
		if err != nil {
			logrus.Warnf("[MakeFinancialStatements] cashflow 不存在, tsCode: %v", tsCode)
			continue
		}
		balanceSheetMap, err := f.getBalanceSheet(ctx, tsCode)
		if err != nil {
			logrus.Warnf("[MakeFinancialStatements] balanceSheet 不存在, tsCode: %v", tsCode)
			continue
		}
		incomeMap, err := f.getIncome(ctx, tsCode)
		if err != nil {
			logrus.Warnf("[MakeFinancialStatements] income 不存在, tsCode: %v", tsCode)
			continue
		}
		stockBasic, err := f.getStockBasic(ctx, tsCode)
		if err != nil {
			logrus.Warnf("[MakeFinancialStatements] stockBasic 不存在, tsCode: %v, err: %v", tsCode, err)
		}

		// 填充初始数据
		eList := make([]*EnterpriseGrowthData, 0)
		for year, c := range cashflowMap {
			b, ok := balanceSheetMap[year]
			if !ok {
				logrus.Warnf("[MakeFinancialStatements] balanceSheet 不存在, tsCode: %v, year: %v", tsCode, year)
				continue
			}
			i, ok := incomeMap[year]
			if !ok {
				logrus.Warnf("[MakeFinancialStatements] income 不存在, tsCode: %v, year: %v", tsCode, year)
				continue
			}
			s, ok := stockBasic[year]
			if !ok {
				logrus.Warnf("[MakeFinancialStatements] stockBasic 不存在, tsCode: %v, year: %v", tsCode, year)
				s = &do.DailyBasic{}
			}

			e := &EnterpriseGrowthData{
				TSCode: tsCode,
				Name:   name,
				Year:   year,

				Revene:       i.Revenue,
				NIncomeAttrP: i.NIncomeAttrP,
				NCashflowAct: c.NCashflowAct,
				TotalAssets:  b.TotalAssets,
				TotalLiab:    b.TotalLiab,
				OperCost:     i.OperCost,

				PE:      s.PE,
				PB:      s.PB,
				PS:      s.PS,
				DVRatio: s.DVRatio,
				TotalMV: s.TotalMV,
				//Close:   s.Close,
			}

			eList = append(eList, e)
		}

		// 按年份排序
		sort.Slice(eList, func(i, j int) bool {
			return eList[i].Year < eList[j].Year
		})

		//过滤年份中断的异常数据
		flag := false
		for i := 1; i < len(eList); i++ {
			if eList[i].Year != eList[i-1].Year+1 {
				logrus.Warnf("[MakeFinancialStatements] 有年份缺失, tsCode: %v, year: %v", tsCode, eList[i].Year)
				flag = true
				continue
			}
		}
		if flag {
			logrus.Warnf("[MakeFinancialStatements] 有年份缺失, tsCode: %v, name: %v", tsCode, name)
			continue
		}

		// 填充关键数据
		for i := 0; i < len(eList); i++ {
			e := eList[i]
			e.NetProfitMargin = e.NIncomeAttrP / e.Revene * 100
			e.AssetLiabilityRatio = e.TotalLiab / e.TotalAssets * 100
			e.EquityMultiplier = 1.0 / (1.0 - e.AssetLiabilityRatio/100)
			e.GrossMargin = (e.Revene - e.OperCost) / e.Revene * 100

			if i > 0 {
				e.TotalAssetTurnover = e.Revene / (e.TotalAssets + eList[i-1].TotalAssets) * 2 * 100
				e.ROE = (e.NetProfitMargin / 100) * (e.TotalAssetTurnover / 100) * e.EquityMultiplier * 100
			} else {
				e.TotalAssetTurnover = -9999
				e.ROE = -9999
			}
		}

		// 填充数据增长
		for i := 1; i < len(eList); i++ {
			e := eList[i]
			e.ReveneGrowth = (e.Revene - eList[i-1].Revene) / eList[i-1].Revene * 100
			e.NIncomeflowActGrowth = (e.NIncomeAttrP - eList[i-1].NIncomeAttrP) / eList[i-1].NIncomeAttrP * 100
			e.NCashflowActGrowth = (e.NCashflowAct - eList[i-1].NCashflowAct) / eList[i-1].NCashflowAct * 100
			e.TotalAssetsGrowth = (e.TotalAssets - eList[i-1].TotalAssets) / eList[i-1].TotalAssets * 100
			e.TotalLiabGrowth = (e.TotalLiab - eList[i-1].TotalLiab) / eList[i-1].TotalLiab * 100
			e.OperCostGrowth = (e.OperCost - eList[i-1].OperCost) / eList[i-1].OperCost * 100
		}

		financialDataMap[tsCode] = eList
	}

	return financialDataMap, nil
}

func (f *Finance) saveAsExcel(ctx context.Context, enterpriseDataMap map[string][]*EnterpriseGrowthData) error {
	for _, enterpriseData := range enterpriseDataMap {
		if len(enterpriseData) == 0 {
			continue
		}

		e, err := excel.NewExcel(ctx, fmt.Sprintf("%v_%v", enterpriseData[0].Name, enterpriseData[0].TSCode))
		if err != nil {
			logrus.Warnf("保存excel失败，tsCode: %v, name: %v, err: %v", enterpriseData[0].TSCode, enterpriseData[0].Name, err)
			continue
		}

		e.SetHeaderList(ctx, tableHead)
		excelDataList := make([][]interface{}, 0)
		for _, data := range enterpriseData {
			excelData := make([]interface{}, 0)
			excelData = append(excelData, data.Year)
			excelData = append(excelData, fmt.Sprintf("%.2f", data.Revene/10000))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.NIncomeAttrP/10000))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.NCashflowAct/10000))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.TotalAssets/10000))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.TotalLiab/10000))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.OperCost/10000))
			excelData = append(excelData, "")
			excelData = append(excelData, fmt.Sprintf("%.2f", data.NetProfitMargin))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.TotalAssetTurnover))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.AssetLiabilityRatio))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.EquityMultiplier))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.ROE))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.GrossMargin))
			excelData = append(excelData, "")
			excelData = append(excelData, fmt.Sprintf("%.2f", data.ReveneGrowth))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.NIncomeflowActGrowth))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.NCashflowActGrowth))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.TotalAssetsGrowth))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.TotalLiabGrowth))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.OperCostGrowth))
			excelData = append(excelData, "")
			excelData = append(excelData, fmt.Sprintf("%.2f", data.PE))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.PB))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.PS))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.DVRatio))
			excelData = append(excelData, fmt.Sprintf("%.2f", data.TotalMV))
			//excelData = append(excelData, fmt.Sprintf("%.2f", data.Close))

			excelDataList = append(excelDataList, excelData)
		}

		e.InsertData(ctx, excelDataList)
		err = e.SaveExcel(ctx)
		if err != nil {
			logrus.Warnf("保存excel失败, tsCode: %v, name: %v, err: %v", enterpriseData[0].TSCode, enterpriseData[0].Name, err)
		}
	}

	return nil
}

var (
	ExcelPath = "excel/"
	tableHead = []string{
		"年份", "营业收入/万", "净利润/万", "现金流量净额/万", "资产总额/万", "负债总额/万", "业务成本/万", "",
		"净利润率", "总资产周转率", "资产负债率", "权益乘数", "ROE", "毛利率", "",
		"营业增长", "净利润增长", "现金流量净额增长", "资产总额增长", "负债总额增长", "业务成本增长", "",
		"市盈率", "市净率", "市销率", "股息率", "流通总市值",
	}
)

type EnterpriseGrowthData struct {
	TSCode string
	Name   string
	Year   int64 //年份

	PE      float64
	PB      float64
	PS      float64
	DVRatio float64 // 股息率
	TotalMV float64 // 流通总市值
	Close   float64 // 收盘价

	Revene       float64 //营业收入
	NIncomeAttrP float64 //净利润
	NCashflowAct float64 //现金流量净额
	TotalAssets  float64 //资产总额
	TotalLiab    float64 //负债总额
	OperCost     float64 //业务成本

	NetProfitMargin     float64 //净利润率
	TotalAssetTurnover  float64 //总资产周转率
	AssetLiabilityRatio float64 //资产负债率
	EquityMultiplier    float64 //权益乘数
	ROE                 float64 //ROE
	GrossMargin         float64 //毛利率

	ReveneGrowth         float64 //营业增长
	NIncomeflowActGrowth float64 //净利润增长
	NCashflowActGrowth   float64 //资产总额增长
	TotalAssetsGrowth    float64 //资产总额增长
	TotalLiabGrowth      float64 //负债总额增长
	OperCostGrowth       float64 //业务成本增长
}

func (f *Finance) getCashflow(ctx context.Context, tsCode string) (map[int64]*do.Cashflow, error) {
	cashflowResult, err := f.CashflowDAL.BatchGetCashflow(ctx, cashflow.BatchGetCashflowParam{StockCodeList: []string{tsCode}})
	if err != nil {
		return nil, err
	}

	cashflowList, ok := cashflowResult.CashflowMap[tsCode]
	if !ok {
		logrus.Errorf("[getCashflow] 获取现金流表失败, tsCode: %v", tsCode)
		err := fmt.Errorf("fetch cashflow failed, tsCode: %v", tsCode)
		return nil, err
	}

	result := make(map[int64]*do.Cashflow)
	for _, c := range cashflowList {
		if c.EndDate[4:] != "1231" {
			continue
		}
		year, err := strconv.ParseInt(c.EndDate[0:4], 10, 64)
		if err != nil {
			return nil, err
		}
		result[year] = c
	}

	return result, nil
}

func (f *Finance) getBalanceSheet(ctx context.Context, tsCode string) (map[int64]*do.BalanceSheet, error) {
	balanceSheetResult, err := f.BalanceSheetDAL.BatchGetBalanceSheet(ctx, balance_sheet.BatchGetBalanceSheetParam{StockCodeList: []string{tsCode}})
	if err != nil {
		return nil, err
	}

	balanceSheetList, ok := balanceSheetResult.BalanceSheetMap[tsCode]
	if !ok {
		logrus.Errorf("[getBalanceSheet] 获取资产负债表失败, tsCode: %v", tsCode)
		err := fmt.Errorf("fetch balanceSheet failed, tsCode: %v", tsCode)
		return nil, err
	}

	result := make(map[int64]*do.BalanceSheet)
	for _, b := range balanceSheetList {
		if b.EndDate[4:] != "1231" {
			continue
		}
		year, err := strconv.ParseInt(b.EndDate[0:4], 10, 64)
		if err != nil {
			return nil, err
		}
		result[year] = b
	}

	return result, nil
}

func (f *Finance) getIncome(ctx context.Context, tsCode string) (map[int64]*do.Income, error) {
	incomeResult, err := f.IncomeDAL.BatchGetIncome(ctx, income.BatchGetIncomeParam{StockCodeList: []string{tsCode}})
	if err != nil {
		return nil, err
	}

	incomeList, ok := incomeResult.IncomeMap[tsCode]
	if !ok {
		logrus.Errorf("[getIncome] 获取利润表失败, tsCode: %v", tsCode)
		err := fmt.Errorf("fetch income failed, tsCode: %v", tsCode)
		return nil, err
	}

	result := make(map[int64]*do.Income)
	for _, i := range incomeList {
		if i.EndDate[4:] != "1231" {
			continue
		}
		year, err := strconv.ParseInt(i.EndDate[0:4], 10, 64)
		if err != nil {
			return nil, err
		}
		result[year] = i
	}

	return result, nil
}

func (f *Finance) getStockBasic(ctx context.Context, tsCode string) (map[int64]*do.DailyBasic, error) {
	basicResult, err := f.DailyBasicDAL.BatchGetDailyBasic(ctx, daily_basic.BatchGetDailyBasicParam{
		TSCodeList: []string{tsCode},
	})
	if err != nil {
		return nil, err
	}

	dailyBasicList, exist := basicResult.DailyBasicMap[tsCode]
	if !exist {
		err := fmt.Errorf("不存在DailyBasic, tsCode: %v", tsCode)
		return nil, err
	}

	result := make(map[int64]*do.DailyBasic, 0)
	for _, i := range dailyBasicList {
		if i.TradeDate[4:6] != "05" {
			continue
		}
		year, err := strconv.ParseInt(i.TradeDate[0:4], 10, 64)
		if err != nil {
			return nil, err
		}
		result[year] = i
	}

	return result, nil
}
