// Package reader define easy to use file reader
package reader

// API defines Reader methods
type API interface {
	// ReadAll read all rows
	ReadAll() ([][]string, error)
	// ReadAysnc read rows async
	ReadAysnc(bufferSize int) (chan []string, error)
}
