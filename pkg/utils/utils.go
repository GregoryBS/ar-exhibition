package utils

import (
	"encoding/json"
	"io"
)

type JSONError struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func DecodeJSON(body io.Reader, dst interface{}) error {
	return json.NewDecoder(body).Decode(dst)
}

func EncodeJSON(src interface{}) []byte {
	result, err := json.Marshal(src)
	if err != nil {
		return nil
	}
	return result
}
