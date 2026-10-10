package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// Decode parses one JSON value, optionally rejecting unknown fields and always rejecting trailing values.
func Decode(data []byte, target any, disallowUnknownFields bool) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if disallowUnknownFields {
		decoder.DisallowUnknownFields()
	}

	if err := decoder.Decode(target); err != nil {
		return err
	}

	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("unexpected trailing JSON value")
		}
		return err
	}

	return nil
}
