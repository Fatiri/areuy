package excel

import (
	"errors"
	"reflect"

	"github.com/xuri/excelize/v2"
)

// GenerateXLSX ...
func GenerateXLSX(sheet string, header []string, data []interface{}) (*excelize.File, error) {
	// Initialize
	xlsx := excelize.NewFile()
	// Set sheet name
	xlsx.SetSheetName("Sheet1", sheet)
	if err := xlsx.SetSheetRow(sheet, "A1", &header); err != nil {
		return nil, err
	}

	for i, val := range data {
		var row []interface{}
		switch reflect.TypeOf(val).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(val)
			for i := 0; i < s.Len(); i++ {
				row = append(row, s.Index(i).Interface())
			}
		default:
			err := errors.New("Data is not an array")
			return nil, err
		}

		addr, err := excelize.JoinCellName("A", (i + 2))
		if err != nil {
			return nil, err
		}

		if err := xlsx.SetSheetRow(sheet, addr, &row); err != nil {
			return nil, err
		}
	}

	return xlsx, nil
}
