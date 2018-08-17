package server

import (
	"bufio"
	"strings"
)

func readFromReader(reader *bufio.Reader) (string, error) {
	text, err := reader.ReadString('\n')
	return strings.TrimSpace(text), err
}
