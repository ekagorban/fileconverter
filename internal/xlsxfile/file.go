package xlsxfile

import (
	"errors"
	"fileconverter/internal/model"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"
)

const (
	sheetName    = "sheet1"
	headRowNum   = 1
	rowsNumShift = 2 // 0, 1 - head, 2... - data
)

type File struct {
	File *os.File
}

func New(f *os.File) *File {
	return &File{f}
}

type sheet struct {
	Name string
	Data model.Object
}

type sheets []sheet

func (file *File) ToFile(objects model.Objects) error {
	sheets, err := prepareSheets(objects)
	if err != nil {
		return err
	}

	err = sheets.toFile(file.File)
	if err != nil {
		return err
	}

	return nil
}

func prepareSheets(data []map[string]any) (sheets, error) {
	dataXLSX := make(map[string]interface{})   // "A1" : data1, "A2" : data2, B3: data3
	columnsByFields := make(map[string]string) // "field1 : A, field2 : B"
	column := getNextSymbol("")

	for _, object := range data {
		for field := range object {
			if _, exist := columnsByFields[field]; !exist {
				columnsByFields[field] = column
				cell := fmt.Sprintf("%s%d", column, headRowNum)
				dataXLSX[cell] = field
				column = getNextSymbol(column)
			}
		}

	}

	for i, obj := range data {
		for field, value := range obj {
			column := columnsByFields[field]
			dataXLSX[fmt.Sprintf("%s%d", column, i+rowsNumShift)] = value

		}

	}

	return sheets{
		{
			Name: sheetName,
			Data: dataXLSX,
		},
	}, nil
}

func (sheets sheets) toFile(file *os.File) error {
	if sheets == nil {
		return errors.New("empty data")
	}

	f := excelize.NewFile()

	activeSheetIndex := f.GetActiveSheetIndex()
	f.SetSheetName(f.GetSheetName(activeSheetIndex), sheets[0].Name)

	for i, sheet := range sheets {
		if i != 0 {
			f.NewSheet(sheet.Name)
		}
		for k, v := range sheet.Data {
			err := f.SetCellValue(sheet.Name, k, v)
			if err != nil {
				return fmt.Errorf("f.SetCellValue error: %v", err)
			}
		}

		setAutofitWidth(f, sheet.Name)
	}

	_, err := f.WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}

func getNextSymbol(currSymbols string) string {
	if len(currSymbols) == 0 {
		return "A"
	}
	last := currSymbols[len(currSymbols)-1]
	if len(currSymbols) == 1 {
		if rune(last) == 'Z' {
			return "AA"
		} else {
			return fmt.Sprintf("%v", string(rune(last)+1))
		}
	} else {
		if rune(last) == 'Z' {
			return fmt.Sprintf("%v%v", getNextSymbol(currSymbols[:len(currSymbols)-1]), string('A'))
		} else {
			return fmt.Sprintf("%v%v", currSymbols[:len(currSymbols)-1], string(rune(last)+1))
		}
	}
}

func setAutofitWidth(f *excelize.File, sheet string) {
	cols, err := f.GetCols(sheet)
	if err != nil {
		log.Printf("f.GetCols error: %v", err)
		return
	}

	for idx, col := range cols {
		largestWidth := 0
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 2 // + 2 for margin
			if cellWidth > largestWidth {
				largestWidth = cellWidth
				if largestWidth > 255 {
					largestWidth = 255
					break
				}
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			log.Printf("excelize.ColumnNumberToNameerror: %v", err)
			return
		}

		err = f.SetColWidth(sheet, name, name, float64(largestWidth))
		if err != nil {
			log.Printf("f.SetColWidth error: %v", err)
		}
	}
}
