package util

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

func PP(s any) {
	res, err := PrettyStruct(s)
	if err != nil {
		slog.Error("couldn't pp", "error", err)
	}
	fmt.Println(res)
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", fmt.Errorf("error marshalling data: %w", err)
	}
	return string(val), nil
}
