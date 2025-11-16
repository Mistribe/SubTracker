package export_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mistribe/subtracker/internal/adapters/http/export"
)

type jsonTestStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestEncodeJSON_BasicEncoding(t *testing.T) {
	data := []jsonTestStruct{
		{ID: "1", Name: "First"},
		{ID: "2", Name: "Second"},
	}

	var buf bytes.Buffer
	err := export.EncodeJSON(&buf, data)
	require.NoError(t, err)

	// Verify it's valid JSON
	var result []jsonTestStruct
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "1", result[0].ID)
	assert.Equal(t, "First", result[0].Name)
	assert.Equal(t, "2", result[1].ID)
	assert.Equal(t, "Second", result[1].Name)
}

func TestEncodeJSON_EmptyArray(t *testing.T) {
	data := []jsonTestStruct{}

	var buf bytes.Buffer
	err := export.EncodeJSON(&buf, data)
	require.NoError(t, err)

	// Verify it's valid JSON and produces an empty array
	var result []jsonTestStruct
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Len(t, result, 0)
	assert.NotNil(t, result) // Should be empty array, not null
}

func TestEncodeJSON_WithOptionalFields(t *testing.T) {
	desc := "A description"
	data := []struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
	}{
		{ID: "1", Name: "With Desc", Description: &desc},
		{ID: "2", Name: "Without Desc", Description: nil},
	}

	var buf bytes.Buffer
	err := export.EncodeJSON(&buf, data)
	require.NoError(t, err)

	// Verify it's valid JSON
	var result []map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// First item has description
	assert.Equal(t, "1", result[0]["id"])
	assert.Equal(t, "With Desc", result[0]["name"])
	assert.Equal(t, "A description", result[0]["description"])

	// Second item doesn't have description (omitempty)
	assert.Equal(t, "2", result[1]["id"])
	assert.Equal(t, "Without Desc", result[1]["name"])
	_, hasDesc := result[1]["description"]
	assert.False(t, hasDesc)
}
