package yaml

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Encode a slice of interface{} into a multi-document yaml byte string.
// Any errors encountered while encoding any of the yaml docs will cause the
// function to fail fast and return the a nil []byte and the error encountered.
// E.g if 10 docs are provided and 3 documents cause errors, a nil slice and an
// errors is returned on the first error encountered.
func Encode(docs ...interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	for _, doc := range docs {
		if err := encoder.Encode(doc); err != nil {
			return nil, fmt.Errorf(`yaml encoding %+v: %w`, doc, err)
		}
	}

	return buf.Bytes(), nil
}

// EncodeNoFailFast returns a slice of interface{} as a multi-document yaml
// byte string.
// Any errors that are encountered while encoding any of the documents provided
// are appended to the error but will NOT force the function to fail fast.
// E.g. if 10 docs are provided and 3 cause an encoding error, the function will
// return a byte slice of yaml containing 7 docs and an error containing 3
// errors appended together.
func EncodeNoFailFast(docs ...interface{}) ([]byte, error) {
	var finalErr error
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	for _, doc := range docs {
		if err := encoder.Encode(doc); err != nil {
			finalErr = fmt.Errorf(`yaml encoding %+v: %w: %v`, doc, err, finalErr)
		}
	}

	return buf.Bytes(), finalErr
}
