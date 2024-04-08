package reader

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// ReadPlainTextFile read plain text from filepath
func ReadPlainTextFile(path string, sep string, ignoreHeader bool) (contents [][]string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	// handle every row
	rowHandle := func(input string){
		if input == ""{
			return 
		}
		input = strings.Trim(input, "\n")
		// split c
		contents = append(contents, strings.Split(input, sep))
	}

	// get csv content
	reader := bufio.NewReader(f)
	for {
		c, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				rowHandle(c)
				break
			} else {
				return contents, err
			}
		}
		rowHandle(c)
	}
	return
}
