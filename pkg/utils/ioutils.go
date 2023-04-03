package utils

import (
	"encoding/json"
	"os"
)

func WriteToJsonFile(path string, obj any) error {
	value, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return os.WriteFile(path, value, os.ModePerm)
}
