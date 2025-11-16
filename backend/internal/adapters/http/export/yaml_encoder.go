package export

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// EncodeYAML writes data as YAML to the provided writer.
// It uses gopkg.in/yaml.v3 with streaming to io.Writer.
// Empty arrays produce valid YAML output.
func EncodeYAML(w io.Writer, data interface{}) error {
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2) // 2-space indentation

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode YAML: %w", err)
	}

	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to close YAML encoder: %w", err)
	}

	return nil
}
