package utils

import (
	"ar_exhibition/pkg/domain"
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

func MapJSON(obj map[string]string) []*domain.Param {
	result := make([]*domain.Param, 0)
	for k, v := range obj {
		result = append(result, &domain.Param{Type: k, Value: v})
	}
	return result
}
