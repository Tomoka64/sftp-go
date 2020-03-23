package csv

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
)

type CSV interface {
	Bytes([][]string) ([]byte, error)
	Read(string) ([][]string, error)
}

type CSVs struct{}

func New() *CSVs {
	return &CSVs{}
}

func (c *CSVs) Bytes(records [][]string) ([]byte, error) {
	var out bytes.Buffer
	csvWriter := csv.NewWriter(&out)
	if err := csvWriter.WriteAll(records); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (c *CSVs) Read(filename string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	var records [][]string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
