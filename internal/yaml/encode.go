package yaml

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Encode a slice of interface{} into a multi-document yaml byte string.
// Errors if any interface{} cannot be encoded.
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
