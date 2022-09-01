package crawlTest

import (
	"errors"
	"github.com/tealeg/xlsx"
)

type headerColumn struct {
	Field string // 字段，数据映射到的数据字段名
	Title string // 标题，表格中的列名称
}

func createFile() (*xlsx.File, *xlsx.Sheet, *[]*headerColumn, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("MOVIESINFO") //表实例
	if err != nil {
		return nil, nil, nil, err
	}
	headers := []*headerColumn{
		{Field: "Name", Title: "电影名称"},
		{Field: "Year", Title: "年代"},
		{Field: "From", Title: "产地"},
		{Field: "Kind", Title: "类别"},
		{Field: "Language", Title: "语言"},
		{Field: "Subtitle", Title: "字幕"},
		{Field: "ReleaseDate", Title: "上映日期"},
		{Field: "Score", Title: "IMDb评分"},
		{Field: "Link", Title: "磁力链接"},
	}
	style := map[string]float64{
		"Name":        2.0,
		"Year":        2.0,
		"From":        2.0,
		"Kind":        2.0,
		"Language":    2.0,
		"Subtitle":    2.0,
		"ReleaseDate": 2.0,
		"Score":       2.0,
		"Link":        2.0,
	}
	sheet, _ = setHeader(sheet, headers, style)
	return file, sheet, &headers, nil
}

func saveFile(file *xlsx.File, name string) error {
	err := file.Save(name)
	if err != nil {
		return err
	}
	return nil
}

func setHeader(sheet *xlsx.Sheet, header []*headerColumn, width map[string]float64) (*xlsx.Sheet, error) {
	if len(header) == 0 {
		return nil, errors.New("Excel.SetHeader 错误: 表头不能为空")
	}

	// 表头样式
	style := xlsx.NewStyle()

	font := xlsx.DefaultFont()
	font.Bold = true

	alignment := xlsx.DefaultAlignment()
	alignment.Vertical = "center"

	style.Font = *font
	style.Alignment = *alignment

	style.ApplyFont = true
	style.ApplyAlignment = true

	// 设置表头字段
	row := sheet.AddRow()
	row.SetHeightCM(1.0)
	row_w := make([]string, 0)
	for _, column := range header {
		row_w = append(row_w, column.Field)
		cell := row.AddCell()
		cell.Value = column.Title
		cell.SetStyle(style) //设置单元样式
	}

	// 表格列，宽度
	if len(row_w) > 0 {
		for k, v := range row_w {
			if width[v] > 0.0 {
				sheet.SetColWidth(k, k, width[v]*10)
			}
		}
	}
	return sheet, nil
}
