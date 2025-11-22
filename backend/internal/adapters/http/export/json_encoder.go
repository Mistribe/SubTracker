package export

import (
	"encoding/json"
	"fmt"
	"io"
)

// EncodeJSON writes data as JSON to the provided writer.
// It uses encoding/json with streaming to io.Writer.
// Empty arrays produce valid JSON output.
func EncodeJSON(w io.Writer, data interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ") // Pretty print with 2-space indentation

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
