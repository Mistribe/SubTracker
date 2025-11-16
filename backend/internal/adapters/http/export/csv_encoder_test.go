package export_test

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mistribe/subtracker/internal/adapters/http/export"
)

type testStruct struct {
	ID    string   `csv:"id"`
	Name  string   `csv:"name"`
	Count int      `csv:"count"`
	Tags  []string `csv:"tags"`
}

type testStructWithPointers struct {
	ID          string  `csv:"id"`
	Name        string  `csv:"name"`
	Description *string `csv:"description"`
	URL         *string `csv:"url"`
}

func TestEncodeCSV_BasicEncoding(t *testing.T) {
	data := []testStruct{
		{ID: "1", Name: "First", Count: 10, Tags: []string{"tag1", "tag2"}},
		{ID: "2", Name: "Second", Count: 20, Tags: []string{"tag3"}},
	}

	var buf bytes.Buffer
	err := export.EncodeCSV(&buf, data)
	require.NoError(t, err)

	// Parse the CSV to verify it's valid
	reader := csv.NewReader(&buf)
	records, err := reader.ReadAll()
	require.NoError(t, err)
	require.Len(t, records, 3) // header + 2 data rows

	// Verify headers
	assert.Equal(t, []string{"id", "name", "count", "tags"}, records[0])

	// Verify first row
	assert.Equal(t, "1", records[1][0])
	assert.Equal(t, "First", records[1][1])
	assert.Equal(t, "10", records[1][2])
	assert.Equal(t, "tag1, tag2", records[1][3])

	// Verify second row
	assert.Equal(t, "2", records[2][0])
	assert.Equal(t, "Second", records[2][1])
	assert.Equal(t, "20", records[2][2])
	assert.Equal(t, "tag3", records[2][3])
}

func TestEncodeCSV_SliceFields(t *testing.T) {
	data := []testStruct{
		{ID: "1", Name: "Test", Count: 5, Tags: []string{"a", "b", "c"}},
		{ID: "2", Name: "Empty", Count: 0, Tags: []string{}},
		{ID: "3", Name: "Single", Count: 1, Tags: []string{"single"}},
	}

	var buf bytes.Buffer
	err := export.EncodeCSV(&buf, data)
	require.NoError(t, err)

	reader := csv.NewReader(&buf)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Verify slice field with multiple values
	assert.Equal(t, "a, b, c", records[1][3])

	// Verify empty slice
	assert.Equal(t, "", records[2][3])

	// Verify single value slice
	assert.Equal(t, "single", records[3][3])
}

func TestEncodeCSV_NilPointerFields(t *testing.T) {
	desc := "A description"
	data := []testStructWithPointers{
		{ID: "1", Name: "With Description", Description: &desc, URL: nil},
		{ID: "2", Name: "No Optional Fields", Description: nil, URL: nil},
	}

	var buf bytes.Buffer
	err := export.EncodeCSV(&buf, data)
	require.NoError(t, err)

	reader := csv.NewReader(&buf)
	records, err := reader.ReadAll()
	require.NoError(t, err)
	require.Len(t, records, 3)

	// Verify headers
	assert.Equal(t, []string{"id", "name", "description", "url"}, records[0])

	// First row has description
	assert.Equal(t, "1", records[1][0])
	assert.Equal(t, "With Description", records[1][1])
	assert.Equal(t, "A description", records[1][2])
	assert.Equal(t, "", records[1][3]) // nil URL

	// Second row has no optional fields
	assert.Equal(t, "2", records[2][0])
	assert.Equal(t, "No Optional Fields", records[2][1])
	assert.Equal(t, "", records[2][2]) // nil description
	assert.Equal(t, "", records[2][3]) // nil URL
}

func TestEncodeCSV_EmptyDataset(t *testing.T) {
	data := []testStruct{}

	var buf bytes.Buffer
	err := export.EncodeCSV(&buf, data)
	require.NoError(t, err)

	// Should still have headers
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	require.Len(t, lines, 1)
	assert.Equal(t, "id,name,count,tags", lines[0])
}

func TestEncodeCSV_InvalidInput(t *testing.T) {
	t.Run("not a slice", func(t *testing.T) {
		var buf bytes.Buffer
		err := export.EncodeCSV(&buf, "not a slice")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a slice")
	})
}
