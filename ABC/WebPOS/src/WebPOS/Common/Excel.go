package Common

import (
	"github.com/goframework/gf/exterror"
	"github.com/goframework/xlsx"
	"os"
)

const (
	SIMPLE_EXCEL_DEFAULT_FONT                   = "ＭＳ Ｐゴシック"
	SIMPLE_EXCEL_DEFAULT_FONT_SIZE              = 10
	SIMPLE_EXCEL_DEFAULT_BORDER                 = "thin"
	SIMPLE_EXCEL_DEFAULT_HEADER_FILL_STYLE      = "solid"
	SIMPLE_EXCEL_DEFAULT_HEADER_FILL_COLOR      = "000000"
	SIMPLE_EXCEL_DEFAULT_HEADER_FILL_BACKGROUND = "dddddd"
	SIMPLE_EXCEL_DEFAULT_VERTICAL_ALIGNMENT     = "center"

	EXCEL_MAX_ROWS = 1048576
)

type SimpleExcelFile struct {
	xlsxFile        *xlsx.File
	xlsxSheet       *xlsx.Sheet
	filePath        string
	headerCellStyle *xlsx.Style
	cellStyle       *xlsx.Style
	fileWriter      *os.File
	rowCount        int
}

func NewSimpleExcelFile(filePath string, sheetName string) (*SimpleExcelFile, error) {

	fileWriter, err := os.Create(filePath)
	if err != nil {
		return nil, exterror.WrapExtError(err)
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet(sheetName)

	if err != nil {
		return nil, exterror.WrapExtError(err)
	}

	headerCellStyle := xlsx.NewStyle()
	headerCellStyle.Font = *xlsx.NewFont(SIMPLE_EXCEL_DEFAULT_FONT_SIZE, SIMPLE_EXCEL_DEFAULT_FONT)
	headerCellStyle.ApplyBorder = true
	headerCellStyle.Border = *xlsx.NewBorder(SIMPLE_EXCEL_DEFAULT_BORDER, SIMPLE_EXCEL_DEFAULT_BORDER, SIMPLE_EXCEL_DEFAULT_BORDER, SIMPLE_EXCEL_DEFAULT_BORDER)
	headerCellStyle.Fill = *xlsx.NewFill(SIMPLE_EXCEL_DEFAULT_HEADER_FILL_STYLE, SIMPLE_EXCEL_DEFAULT_HEADER_FILL_BACKGROUND, SIMPLE_EXCEL_DEFAULT_HEADER_FILL_COLOR)
	headerCellStyle.ApplyFill = true
	headerCellStyle.Alignment.Vertical = SIMPLE_EXCEL_DEFAULT_VERTICAL_ALIGNMENT

	cellStyle := xlsx.NewStyle()
	cellStyle.Font = *xlsx.NewFont(SIMPLE_EXCEL_DEFAULT_FONT_SIZE, SIMPLE_EXCEL_DEFAULT_FONT)
	cellStyle.ApplyBorder = true
	cellStyle.Border = *xlsx.NewBorder(SIMPLE_EXCEL_DEFAULT_BORDER, SIMPLE_EXCEL_DEFAULT_BORDER, SIMPLE_EXCEL_DEFAULT_BORDER, SIMPLE_EXCEL_DEFAULT_BORDER)
	cellStyle.Alignment.Vertical = SIMPLE_EXCEL_DEFAULT_VERTICAL_ALIGNMENT

	sef := SimpleExcelFile{
		file,
		sheet,
		filePath,
		headerCellStyle,
		cellStyle,
		fileWriter,
		0,
	}

	return &sef, nil
}

func (this *SimpleExcelFile) WriteHeader(header []string) {
	if this.xlsxSheet != nil {
		if this.rowCount < EXCEL_MAX_ROWS {
			row := this.xlsxSheet.AddRow()
			this.rowCount++
			row.WriteSlice(&header, -1)
			for _, cell := range row.Cells {
				cell.SetStyle(this.headerCellStyle)
			}
		}
	}
}

func (this *SimpleExcelFile) WriteData(rowData []interface{}) {
	if this.xlsxSheet != nil {
		if this.rowCount < EXCEL_MAX_ROWS {
			this.rowCount++

			row := this.xlsxSheet.AddRow()
			row.WriteSlice(&rowData, -1)
			for _, cell := range row.Cells {
				cell.SetStyle(this.cellStyle)
			}
		}
	}
}

func (this *SimpleExcelFile) WriteDataStruct(rowData interface{}) {
	if this.xlsxSheet != nil {
		if this.rowCount < EXCEL_MAX_ROWS {
			this.rowCount++
			row := this.xlsxSheet.AddRow()
			row.WriteStruct(rowData, -1)
			for _, cell := range row.Cells {
				cell.SetStyle(this.cellStyle)
			}
		}
	}
}

func (this *SimpleExcelFile) SetColWidth(listColWidth []int) {
	if this.xlsxSheet != nil {
		for i, v := range listColWidth {
			this.xlsxSheet.SetColWidth(i, i, float64(v))
		}
	}
}

// Save file and close
func (this *SimpleExcelFile) Close() {
	if this.fileWriter != nil && this.xlsxFile != nil {
		this.xlsxFile.Write(this.fileWriter)
		this.fileWriter.Close()
		this.fileWriter = nil
		this.xlsxFile = nil
		this.xlsxSheet = nil
	}
}

// Delete file without saving
func (this *SimpleExcelFile) Destroy() {
	if this.fileWriter != nil && this.xlsxFile != nil {
		this.fileWriter.Close()
		os.Remove(this.filePath)

		this.fileWriter = nil
		this.xlsxFile = nil
		this.xlsxSheet = nil
	}
}
