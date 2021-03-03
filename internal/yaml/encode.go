package yaml

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Encode a slice of interface{} into a multi-document yaml string.
// Errors if any interface{} cannot be encoded.
func Encode(docs []interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	for _, doc := range docs {
		if err := encoder.Encode(doc); err != nil {
			return "", fmt.Errorf(`yaml encoding %+v: %w`, doc, err)
		}
	}

	return buf.String(), nil
}
