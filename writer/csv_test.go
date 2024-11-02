package writer

import (
	"testing"

	"github.com/onesaltedseafish/go-utils/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCsvWriter(t *testing.T) {
	path := "test.csv"
	csvWriter, err := NewRowCsvWriter(
		path, []string{"a", "b"},
	)
	assert.Equal(t, nil, err)
	_ = csvWriter.WriteHeader()
	err = csvWriter.checkLine([]string{"1", "2"})
	assert.Equal(t, nil, err)
	err = csvWriter.checkLine([]string{"1"})
	require.ErrorIs(t, err, ErrLineItemsNotEqualToHeader)
	err = csvWriter.checklines([][]string{{"3", "4"}, {"5", "6"}})
	assert.Equal(t, nil, err)
	err = csvWriter.checklines([][]string{{"3", "4"}, {"5"}, {"2", "4", "5"}})
	require.ErrorIs(t, err, ErrLineItemsNotEqualToHeader)
	// test write
	err = csvWriter.WriteLine([]string{"1", "2"})
	assert.Equal(t, nil, err)
	err = csvWriter.WriteLines([][]string{{"3", "4"}, {"5", "6"}})
	assert.Equal(t, nil, err)

	reader, err := reader.NewCsvImpl(path)
	assert.Equal(t, nil, err)
	contents, err := reader.ReadAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, contents, [][]string{
		{"a", "b"},
		{"1", "2"},
		{"3", "4"},
		{"5", "6"},
	})
}
