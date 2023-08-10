package fund

func PrintExcel() {
	//logrus.Infof("--------------开始打印avgProfit前十股票-----------")
	//for i, logResult := range logResultListWithAvgProfit {
	//	if i >= 10 {
	//		break
	//	}
	//	logrus.Infof("名称：%v, tsCode: %v, 周期长度: %v, 最少盈利: %v，平均盈利: %v, 发行时间: %v, ma: %v, ma2: %v",
	//		logResult.Name, logResult.TSCode, logResult.TargetTimeLen, logResult.MinProfit,
	//		logResult.AvgProfit, logResult.ListDate, logResult.MA, logResult.MA2)
	//}
	//
	//logResultListWithTargetTimeLen := make([]*LogResult, 0)
	//logResultMapWithTargetTimeLen.Range(func(key, value interface{}) bool {
	//	logResultListWithTargetTimeLen = append(logResultListWithTargetTimeLen, value.(*LogResult))
	//	return true
	//})
	//sort.Slice(logResultListWithTargetTimeLen, func(i, j int) bool {
	//	return logResultListWithTargetTimeLen[i].TargetTimeLen < logResultListWithTargetTimeLen[j].TargetTimeLen
	//})
	//logrus.Infof("--------------开始打印targetTimeLen前十股票-----------")
	//for i, logResult := range logResultListWithTargetTimeLen {
	//	if i >= 10 {
	//		break
	//	}
	//	logrus.Infof("名称：%v, tsCode: %v, 周期长度: %v, 最少盈利: %v，平均盈利: %v, 发行时间: %v, ma: %v, ma2: %v",
	//		logResult.Name, logResult.TSCode, logResult.TargetTimeLen, logResult.MinProfit,
	//		logResult.AvgProfit, logResult.ListDate, logResult.MA, logResult.MA2)
	//}
	//
	//excelFilePath := fmt.Sprintf("%v%v_%v.xlsx", ExcelPath, "ETF数据", time.Now().Format("01-02T15:04:05"))
	//f2 := excelize.NewFile()
	//defer func() {
	//	if err := f2.Close(); err != nil {
	//		logrus.Errorf("[saveAsExcel] 关闭excel失败: %v", err)
	//	}
	//}()
	//// 创建一个工作表
	//index, err := f2.NewSheet("Sheet1")
	//if err != nil {
	//	return err
	//}
	//// 设置单元格的值
	//for cell, value := range tableHead {
	//	f2.SetCellValue("Sheet1", cell, value)
	//}
	//
	//for r, resultData := range resultDataList {
	//	f2.SetCellValue("Sheet1", fmt.Sprintf("A%v", r+2), resultData.Name)
	//	f2.SetCellValue("Sheet1", fmt.Sprintf("B%v", r+2), resultData.TSCode)
	//	f2.SetCellValue("Sheet1", fmt.Sprintf("C%v", r+2), resultData.ListDate)
	//	f2.SetCellValue("Sheet1", fmt.Sprintf("D%v", r+2), resultData.MaxDecline)
	//	f2.SetCellValue("Sheet1", fmt.Sprintf("E%v", r+2), resultData.MaxResult)
	//	f2.SetCellValue("Sheet1", fmt.Sprintf("F%v", r+2), resultData.AvgResult)
	//
	//}
	//// 设置工作簿的默认工作表
	//f2.SetActiveSheet(index)
	//// 根据指定路径保存文件
	//if err := f2.SaveAs(excelFilePath); err != nil {
	//	return err
	//}
	//
	//return nil
}
