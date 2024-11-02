package reader

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

var _ API = (*PlainTextFileImpl)(nil)

// PlainTextFileImpl txt
// NOT SAFE FOR CONCURRENT USE
type PlainTextFileImpl struct {
	path        string
	sep         string
	ignoreNRows int
	// internal use
	reader   *bufio.Reader
	rowIndex int
}

// NewPlainTextFileImpl new
func NewPlainTextFileImpl(path, sep string, ignoreNRows int) (*PlainTextFileImpl, error) {
	if ignoreNRows < 0 {
		ignoreNRows = 0
	}
	if sep == "" {
		sep = "\t"
	}
	impl := &PlainTextFileImpl{
		path:        path,
		sep:         sep,
		ignoreNRows: ignoreNRows,
	}
	f, err := os.Open(impl.path)
	if err != nil {
		return nil, err
	}
	impl.reader = bufio.NewReader(f)
	return impl, nil
}

// ReadAll read all rows
func (impl *PlainTextFileImpl) ReadAll() (contents [][]string, err error) {
	for {
		c, err := impl.readLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if len(c) != 0 {
					contents = append(contents, c)
				}
				break
			} else {
				return nil, err
			}
		}
		contents = append(contents, c)
	}
	return contents, nil
}

func (impl *PlainTextFileImpl) readLine() ([]string, error) {
	impl.rowIndex += 1
	input, err := impl.reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if impl.rowIndex <= impl.ignoreNRows { // ingore n rows
		return impl.readLine()
	}
	if input == "" {
		return nil, err
	}
	input = strings.Trim(input, "\n")
	// split by sep, maybe return io.EOF
	return strings.Split(input, impl.sep), err
}

// ReadAysnc read rows async
func (impl *PlainTextFileImpl) ReadAysnc(bufferSize int) (chan []string, error) {
	result := make(chan []string, bufferSize)
	go func() {
		for {
			c, err := impl.readLine()
			if err != nil {
				if errors.Is(err, io.EOF) {
					if len(c) != 0 {
						result <- c
					}
				}
				break
			}
			result <- c
		}
		close(result)
	}()
	return result, nil
}
