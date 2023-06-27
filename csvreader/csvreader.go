package csvreader

import (
	"bytes"
	"github.com/jfyne/csvd"
)

type CsvReader interface {
	Open(file *bytes.Buffer)
	ReadLine() (records [][]string, err error)
}

type csvReader struct {
	File *bytes.Buffer
}

func NewCsvReader() CsvReader {
	return &csvReader{}
}

func (c *csvReader) Open(file *bytes.Buffer) {
	c.File = file
}

func (c *csvReader) ReadLine() (records [][]string, err error) {
	reader := csvd.NewReader(c.File)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, err
}

