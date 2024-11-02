package reader

import (
	"encoding/csv"
	"os"
)

// CsvImpl impls CSV file
type CsvImpl struct {
	path string
	// internal
	reader *csv.Reader
}

// NewCsvImpl new csv impl
func NewCsvImpl(path string) (*CsvImpl, error) {
	impl := &CsvImpl{
		path: path,
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// get csv content
	impl.reader = csv.NewReader(f)
	return impl, nil
}

// ReadAll read all rows
func (impl *CsvImpl) ReadAll() ([][]string, error) {
	return impl.reader.ReadAll()
}

// ReadAysnc read rows async
func (impl *CsvImpl) ReadAysnc(bufferSize int) (chan []string, error) {
	recordChan := make(chan []string, bufferSize)
	go func() {
		for {
			line, err := impl.reader.Read()
			if err != nil {
				break
			}
			recordChan <- line
		}
		close(recordChan)
	}()
	return recordChan, nil
}
