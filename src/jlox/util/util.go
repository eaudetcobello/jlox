package cmdutil

import (
	"os"
)

func ReadFile(name string) ([]byte, error) {
	bytes, err := os.ReadFile(name)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
