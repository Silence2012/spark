package utils

import (
	"github.com/tealeg/xlsx"
)

func GenerateExcel(excelTitle []string, excelData []map[int][]string, excelPath string) error {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return err
	}
	//生成excel第一行，一般是表头
	row = sheet.AddRow()
	for _, title := range excelTitle{
		cell = row.AddCell()
		cell.Value = title
	}

	//生成内容
	for index, data := range excelData {
		row = sheet.AddRow()
		rowDataList := data[index]
		for _, content := range rowDataList {
			cell = row.AddCell()
			cell.Value = content
		}
	}

	err = file.Save(excelPath)
	return err
}