package reader

import (
	"encoding/csv"
	"os"
)

// ReadCsvFile read csv from filepath
func ReadCsvFile(path string, ignoreHeader bool) (contents [][]string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	// get csv content
	reader := csv.NewReader(f)
	contents, err = reader.ReadAll()
	return
}
