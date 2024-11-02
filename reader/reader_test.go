package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadCsv(t *testing.T) {
	testcases := []struct {
		Path         string
		IgnoreHeader bool
		WantContents [][]string
	}{
		{"tests/1.csv", false, [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}},
		{"tests/1.1.csv", false, [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}},
	}

	for _, testcase := range testcases {
		r, err := ReadCsvFile(testcase.Path, testcase.IgnoreHeader)
		assert.Equal(t, nil, err)
		assert.Equal(t, testcase.WantContents, r)
	}
}

func TestReadPlainText(t *testing.T) {
	testcases := []struct {
		Path         string
		Sep          string
		IgnoreNRows  int
		Async        bool
		WantContents [][]string
	}{
		{"tests/1.txt", "\t", 0, false, [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}},
		{"tests/1.1.txt", "\t", 1, false, [][]string{{"3", "4"}, {"5", "6"}}},
		{"tests/1.txt", "\t", 0, true, [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}},
		{"tests/1.1.txt", "\t", 1, true, [][]string{{"3", "4"}, {"5", "6"}}},
	}

	for _, testcase := range testcases {
		reader, err := NewPlainTextFileImpl(testcase.Path, testcase.Sep, testcase.IgnoreNRows)
		assert.Equal(t, nil, err)
		if testcase.Async {
			var result [][]string
			readChan, _ := reader.ReadAysnc(1)
			for c := range readChan {
				result = append(result, c)
			}
			assert.Equal(t, testcase.WantContents, result)
		} else {
			contents, err := reader.ReadAll()
			assert.Equal(t, nil, err)
			assert.Equal(t, testcase.WantContents, contents)
		}
	}
}
