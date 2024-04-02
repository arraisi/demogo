package utils

import (
	"bytes"
	"encoding/csv"
)

func WriteCSV(data [][]string) (result bytes.Buffer, err error) {
	var buffer bytes.Buffer
	csvWriter := csv.NewWriter(&buffer)
	csvWriter.Comma = ';'

	for _, d := range data {
		if err = csvWriter.Write(d); err != nil {
			return
		}
	}

	csvWriter.Flush()
	if err = csvWriter.Error(); err != nil {
		return bytes.Buffer{}, err
	}

	return buffer, nil
}
