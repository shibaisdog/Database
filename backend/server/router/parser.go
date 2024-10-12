package router

import (
	"bytes"
	"encoding/json"
)

func ParseJSON[T any](f []byte, t *T) error {
	decoder := json.NewDecoder(bytes.NewReader(f))
	return decoder.Decode(t)
}
