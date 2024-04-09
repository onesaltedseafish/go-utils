// Package writer defines useful writer
package writer

import "errors"

var (
	// ErrLineItemsNotEqualToHeader defines 
	ErrLineItemsNotEqualToHeader = errors.New("line column cnt not equal to header")
)

// RowWriter defines writer like csv
// which keeps evey line has equal length
type RowWriter interface {
	WriteLine([]string) error
	WriteLines([][]string) error
	WriteHeader() error
}
