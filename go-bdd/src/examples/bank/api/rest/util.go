package rest

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

func unmarshal[T any](r io.Reader) (T, error) {
	var t T
	if err := json.NewDecoder(r).Decode(&t); err != nil {
		return *new(T), errors.Wrap(err, "failed to unmarshal entity")
	}

	return t, nil
}
