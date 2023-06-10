package cmd

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"my_stock_market/model/do"
	"my_stock_market/repo/interface/balance_sheet"
	"my_stock_market/repo/interface/cashflow"
	"my_stock_market/repo/interface/income"
	"my_stock_market/service/interface/tushare"
	"sort"
	"strconv"
)

type EnterpriseGrowthData struct {
	TSCode string
	Name   string
	Year   int64 //年份

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

// MakeFinancialStatements 计算财务报表
func (s *Stock) MakeFinancialStatements(ctx context.Context) error {

	// 获取所有股票的财务数据
	fullFinancialData, err := s.getFullFinancialData(ctx, []string{})
	if err != nil {
		return err
	}
	logrus.Infof("[MakeFinancialStatements] 获取了%v个公司的财务数据", len(fullFinancialData))
	// 过滤出值得观察的企业
	goodEnterpriseData, err := s.filterGoodEnterpriseData(ctx, fullFinancialData)
	if err != nil {
		return err
	}
	logrus.Infof("[MakeFinancialStatements] 过滤出%v个值得观察的企业", len(goodEnterpriseData))

	// 保存为excel
	err = s.saveAsExcel(ctx, goodEnterpriseData)
	if err != nil {
		return err
	}
	logrus.Infof("[MakeFinancialStatements] 保存为excel成功")

	return nil
}

var (
	ExcelPath = "excel/"
	tableHead = map[string]string{
		"A1": "年份",
		"B1": "营业收入/万",
		"C1": "净利润/万",
		"D1": "现金流量净额/万",
		"E1": "资产总额/万",
		"F1": "负债总额/万",
		"G1": "业务成本/万",
		"H1": "",
		"I1": "净利润率",
		"J1": "总资产周转率",
		"K1": "资产负债率",
		"L1": "权益乘数",
		"M1": "ROE",
		"N1": "毛利率",
		"O1": "",
		"P1": "营业增长",
		"Q1": "净利润增长",
		"R1": "现金流量净额增长",
		"S1": "资产总额增长",
		"T1": "负债总额增长",
		"U1": "业务成本增长",
	}
)

func (s *Stock) saveAsExcel(ctx context.Context, enterpriseDataMap map[string][]*EnterpriseGrowthData) error {
	for _, enterpriseData := range enterpriseDataMap {
		tsCode := enterpriseData[0].TSCode
		name := enterpriseData[0].Name
		excelFilePath := fmt.Sprintf("%v%v-%v.xlsx", ExcelPath, tsCode, name)
		logrus.Infof("[saveAsExcel] 正在保存: tsCode: %v, name: %v, 路径: %v", tsCode, name, excelFilePath)

		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				logrus.Errorf("[saveAsExcel] 关闭excel失败: %v, tsCode: %v", err, tsCode)
			}
		}()

		// 创建一个工作表
		index, err := f.NewSheet("Sheet1")
		if err != nil {
			logrus.Errorf("[saveAsExcel] 创建工作表失败: %v, tsCode: %v", err, tsCode)
			continue
		}
		// 设置单元格的值
		for cell, value := range tableHead {
			f.SetCellValue("Sheet1", cell, value)
		}
		for r, data := range enterpriseData {
			// 使用 data 数据 填充excel
			f.SetCellValue("Sheet1", fmt.Sprintf("A%v", r+2), data.Year)
			f.SetCellValue("Sheet1", fmt.Sprintf("B%v", r+2), fmt.Sprintf("%.2f", data.Revene/10000))
			f.SetCellValue("Sheet1", fmt.Sprintf("C%v", r+2), fmt.Sprintf("%.2f", data.NIncomeAttrP/10000))
			f.SetCellValue("Sheet1", fmt.Sprintf("D%v", r+2), fmt.Sprintf("%.2f", data.NCashflowAct/10000))
			f.SetCellValue("Sheet1", fmt.Sprintf("E%v", r+2), fmt.Sprintf("%.2f", data.TotalAssets/10000))
			f.SetCellValue("Sheet1", fmt.Sprintf("F%v", r+2), fmt.Sprintf("%.2f", data.TotalLiab/10000))
			f.SetCellValue("Sheet1", fmt.Sprintf("G%v", r+2), fmt.Sprintf("%.2f", data.OperCost/10000))
			f.SetCellValue("Sheet1", fmt.Sprintf("H%v", r+2), "")
			f.SetCellValue("Sheet1", fmt.Sprintf("I%v", r+2), fmt.Sprintf("%.2f", data.NetProfitMargin))
			f.SetCellValue("Sheet1", fmt.Sprintf("J%v", r+2), fmt.Sprintf("%.2f", data.TotalAssetTurnover))
			f.SetCellValue("Sheet1", fmt.Sprintf("K%v", r+2), fmt.Sprintf("%.2f", data.AssetLiabilityRatio))
			f.SetCellValue("Sheet1", fmt.Sprintf("L%v", r+2), fmt.Sprintf("%.2f", data.EquityMultiplier))
			f.SetCellValue("Sheet1", fmt.Sprintf("M%v", r+2), fmt.Sprintf("%.2f", data.ROE))
			f.SetCellValue("Sheet1", fmt.Sprintf("N%v", r+2), fmt.Sprintf("%.2f", data.GrossMargin))
			f.SetCellValue("Sheet1", fmt.Sprintf("O%v", r+2), "")
			f.SetCellValue("Sheet1", fmt.Sprintf("P%v", r+2), fmt.Sprintf("%.2f", data.ReveneGrowth))
			f.SetCellValue("Sheet1", fmt.Sprintf("Q%v", r+2), fmt.Sprintf("%.2f", data.NIncomeflowActGrowth))
			f.SetCellValue("Sheet1", fmt.Sprintf("R%v", r+2), fmt.Sprintf("%.2f", data.NCashflowActGrowth))
			f.SetCellValue("Sheet1", fmt.Sprintf("S%v", r+2), fmt.Sprintf("%.2f", data.TotalAssetsGrowth))
			f.SetCellValue("Sheet1", fmt.Sprintf("T%v", r+2), fmt.Sprintf("%.2f", data.TotalLiabGrowth))
			f.SetCellValue("Sheet1", fmt.Sprintf("U%v", r+2), fmt.Sprintf("%.2f", data.OperCostGrowth))
		}

		// 设置工作簿的默认工作表
		f.SetActiveSheet(index)
		// 根据指定路径保存文件
		if err := f.SaveAs(excelFilePath); err != nil {
			logrus.Errorf("[saveAsExcel] 保存excel失败: %v, tsCode: %v", err, tsCode)
			continue
		}
	}

	return nil
}

func (s *Stock) getFullFinancialData(ctx context.Context, tsCodeList []string) (map[string][]*EnterpriseGrowthData, error) {
	stockBasicResult, err := s.TuShare.StockBasic(ctx, tushare.StockBasicParam{Limit: 100000})
	if err != nil {
		return nil, err
	}

	financialDataMap := make(map[string][]*EnterpriseGrowthData)

	for _, stockBasic := range stockBasicResult.StockList {
		exist := false
		for _, tsCode := range tsCodeList {
			if stockBasic.TSCode == tsCode {
				exist = true
				break
			}
		}
		if len(tsCodeList) > 0 && !exist {
			continue
		}

		tsCode := stockBasic.TSCode
		name := stockBasic.Name

		cashflowMap, err := s.getCashflow(ctx, tsCode)
		if err != nil {
			logrus.Errorf("[MakeFinancialStatements] cashflow 不存在, tsCode: %v", tsCode)
			continue
		}
		balanceSheetMap, err := s.getBalanceSheet(ctx, tsCode)
		if err != nil {
			logrus.Errorf("[MakeFinancialStatements] balanceSheet 不存在, tsCode: %v", tsCode)
			continue
		}
		incomeMap, err := s.getIncome(ctx, tsCode)
		if err != nil {
			logrus.Errorf("[MakeFinancialStatements] income 不存在, tsCode: %v", tsCode)
			continue
		}

		// 填充初始数据
		eList := make([]*EnterpriseGrowthData, 0)
		for year, c := range cashflowMap {
			b, ok := balanceSheetMap[year]
			if !ok {
				logrus.Errorf("[MakeFinancialStatements] balanceSheet 不存在, tsCode: %v, year: %v", tsCode, year)
				continue
			}
			i, ok := incomeMap[year]
			if !ok {
				logrus.Errorf("[MakeFinancialStatements] income 不存在, tsCode: %v, year: %v", tsCode, year)
				continue
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
				logrus.Errorf("[MakeFinancialStatements] 有年份缺失, tsCode: %v, year: %v", tsCode, eList[i].Year)
				flag = true
				continue
			}
		}
		if flag {
			logrus.Infof("[MakeFinancialStatements] 有年份缺失, tsCode: %v, name: %v", tsCode, name)
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

// filterGoodEnterpriseData 按照策略过滤值得观察的企业
// 1. 上市时间达到3年，财报长度达到三年
// 2. 如果上市时间超过5年，最近5年ROE都要在15%以上
// 3. 如果上市时间不超过5年，则最近每年ROE都要在15%以上
func (s *Stock) filterGoodEnterpriseData(ctx context.Context, enterpriseDataMap map[string][]*EnterpriseGrowthData) (map[string][]*EnterpriseGrowthData, error) {
	result := make(map[string][]*EnterpriseGrowthData)

	for _, eList := range enterpriseDataMap {
		if len(eList) < 3 {
			continue
		}

		flag := true

		if len(eList) > 5 {
			for i := len(eList) - 1; i >= len(eList)-5; i-- {
				if eList[i].ROE < 15 {
					flag = false
					continue
				}
			}
		} else {
			for i := 1; i < len(eList); i++ {
				if eList[i].ROE < 15 {
					flag = false
					continue
				}
			}
		}

		if flag {
			result[eList[0].TSCode] = eList
		}
	}

	return result, nil
}

func (s *Stock) getCashflow(ctx context.Context, tsCode string) (map[int64]*do.Cashflow, error) {
	cashflowResult, err := s.CashflowDAL.BatchGetCashflow(ctx, cashflow.BatchGetCashflowParam{StockCodeList: []string{tsCode}})
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

func (s *Stock) getBalanceSheet(ctx context.Context, tsCode string) (map[int64]*do.BalanceSheet, error) {
	balanceSheetResult, err := s.BalanceSheetDAL.BatchGetBalanceSheet(ctx, balance_sheet.BatchGetBalanceSheetParam{StockCodeList: []string{tsCode}})
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

func (s *Stock) getIncome(ctx context.Context, tsCode string) (map[int64]*do.Income, error) {
	incomeResult, err := s.IncomeDAL.BatchGetIncome(ctx, income.BatchGetIncomeParam{StockCodeList: []string{tsCode}})
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
