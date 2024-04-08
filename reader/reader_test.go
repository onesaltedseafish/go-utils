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
		IgnoreHeader bool
		WantContents [][]string
	}{
		{"tests/1.txt", "\t", false, [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}},
		{"tests/1.1.txt", "\t", false, [][]string{{"1", "2"}, {"3", "4"}, {"5", "6"}}},
	}

	for _, testcase := range testcases {
		r, err := ReadPlainTextFile(testcase.Path, testcase.Sep, testcase.IgnoreHeader)
		assert.Equal(t, nil, err)
		assert.Equal(t, testcase.WantContents, r)
	}
}
