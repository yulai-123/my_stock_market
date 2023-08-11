package excel

import (
	"context"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"time"
)

type Excel struct {
	excelFile     *excelize.File
	excelFilePath string
}

func NewExcel(ctx context.Context, excelName string) (*Excel, error) {
	excel := &Excel{}

	excelFilePath := fmt.Sprintf("%v%v_%v.xlsx", "excel/", excelName, time.Now().Format("01-02T15:04:05"))
	excelFile := excelize.NewFile()
	index, err := excelFile.NewSheet("Sheet1")
	if err != nil {
		return nil, err
	}
	excelFile.SetActiveSheet(index)

	excel.excelFile = excelFile
	excel.excelFilePath = excelFilePath

	return excel, nil
}

func (e *Excel) SetHeaderList(ctx context.Context, headerList []string) {
	for index, value := range headerList {
		cell := string('A'+index) + "1"
		e.excelFile.SetCellValue("Sheet1", cell, value)
	}
}

func (e *Excel) InsertData(ctx context.Context, dataList [][]interface{}) {
	for r, data := range dataList {
		for index, value := range data {
			cell := string('A'+index) + strconv.Itoa(r+2)
			e.excelFile.SetCellValue("Sheet1", cell, value)
		}
	}
}

func (e *Excel) SaveExcel(ctx context.Context) error {
	err := e.excelFile.SaveAs(e.excelFilePath)
	if err != nil {
		return err
	}

	err = e.excelFile.Close()
	if err != nil {
		return err
	}

	return nil
}
