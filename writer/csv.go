package writer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

var _ RowWriter = (*RowCsvWriter)(nil)

// RowCsvWriter defines RowWriter with format csv
type RowCsvWriter struct {
	Headers   []string
	writer    *csv.Writer
	columnCnt int
}

// NewRowCsvWriter new a csv writer
func NewRowCsvWriter(path string, header []string) (*RowCsvWriter, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	w := &RowCsvWriter{
		Headers:   header,
		writer:    csv.NewWriter(f),
		columnCnt: len(header),
	}
	return w, nil
}

// NewRowCsvWriterWithF new with io.writer
func NewRowCsvWriterWithF(
	writer io.Writer, header []string,
) *RowCsvWriter {
	return &RowCsvWriter{
		Headers:   header,
		columnCnt: len(header),
		writer:    csv.NewWriter(writer),
	}
}

// WriteLine write one line per time
func (w *RowCsvWriter) WriteLine(line []string) error {
	defer w.writer.Flush()
	err := w.checkLine(line)
	if err != nil {
		return err
	}
	return w.writer.Write(line)
}

// WriteHeader write header to csv
func (w *RowCsvWriter) WriteHeader() error {
	defer w.writer.Flush()
	return w.writer.Write(w.Headers)
}

// WriteLines write multi lines per time
func (w *RowCsvWriter) WriteLines(lines [][]string) error {
	defer w.writer.Flush()
	err := w.checklines(lines)
	if err != nil {
		return err
	}
	return w.writer.WriteAll(lines)
}

func (w *RowCsvWriter) checkLine(line []string) error {
	if len(line) != w.columnCnt {
		return fmt.Errorf("%w, want: %d ,got: %d", ErrLineItemsNotEqualToHeader,
			w.columnCnt,
			len(line),
		)
	}
	return nil
}

func (w *RowCsvWriter) checklines(lines [][]string) error {
	var err error
	for index, line := range lines {
		tErr := w.checkLine(line)
		if tErr != nil {
			err = errors.Join(err,
				fmt.Errorf("%w with line no: %d", tErr, index),
			)
		}
	}
	return err
}
