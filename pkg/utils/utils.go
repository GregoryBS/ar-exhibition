package utils

import (
	"ar_exhibition/pkg/domain"
	"encoding/json"
	"io"
	"math/rand"
	"strings"
	"time"
)

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

func SplitPic(pics string) []string {
	buf := strings.Split(pics, ",")
	for i := range buf {
		buf[i] = ImageService + buf[i]
	}
	return buf
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
