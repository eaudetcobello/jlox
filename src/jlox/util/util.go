package cmdutil

import (
	"os"
	"strconv"
)

func ReadFile(name string) ([]byte, error) {
	bytes, err := os.ReadFile(name)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func IsDigit(str string) bool {
	if _, err := strconv.Atoi(str); err != nil {
		return false
	}
	return true
}
