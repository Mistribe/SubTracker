package export_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/mistribe/subtracker/internal/adapters/http/export"
)

type yamlTestStruct struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

func TestEncodeYAML_BasicEncoding(t *testing.T) {
	data := []yamlTestStruct{
		{ID: "1", Name: "First"},
		{ID: "2", Name: "Second"},
	}

	var buf bytes.Buffer
	err := export.EncodeYAML(&buf, data)
	require.NoError(t, err)

	// Verify it's valid YAML
	var result []yamlTestStruct
	err = yaml.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "1", result[0].ID)
	assert.Equal(t, "First", result[0].Name)
	assert.Equal(t, "2", result[1].ID)
	assert.Equal(t, "Second", result[1].Name)
}

func TestEncodeYAML_EmptyArray(t *testing.T) {
	data := []yamlTestStruct{}

	var buf bytes.Buffer
	err := export.EncodeYAML(&buf, data)
	require.NoError(t, err)

	// Verify it's valid YAML and produces an empty array
	var result []yamlTestStruct
	err = yaml.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Len(t, result, 0)
	assert.NotNil(t, result) // Should be empty array, not null
}

func TestEncodeYAML_WithOptionalFields(t *testing.T) {
	desc := "A description"
	data := []struct {
		ID          string  `yaml:"id"`
		Name        string  `yaml:"name"`
		Description *string `yaml:"description,omitempty"`
	}{
		{ID: "1", Name: "With Desc", Description: &desc},
		{ID: "2", Name: "Without Desc", Description: nil},
	}

	var buf bytes.Buffer
	err := export.EncodeYAML(&buf, data)
	require.NoError(t, err)

	// Verify it's valid YAML
	var result []map[string]interface{}
	err = yaml.Unmarshal(buf.Bytes(), &result)
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

func TestEncodeYAML_ComplexStructure(t *testing.T) {
	data := []struct {
		ID     string   `yaml:"id"`
		Name   string   `yaml:"name"`
		Tags   []string `yaml:"tags"`
		Active bool     `yaml:"active"`
	}{
		{ID: "1", Name: "Test", Tags: []string{"a", "b"}, Active: true},
		{ID: "2", Name: "Another", Tags: []string{}, Active: false},
	}

	var buf bytes.Buffer
	err := export.EncodeYAML(&buf, data)
	require.NoError(t, err)

	// Verify it's valid YAML
	var result []map[string]interface{}
	err = yaml.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify first item
	assert.Equal(t, "1", result[0]["id"])
	assert.Equal(t, "Test", result[0]["name"])
	assert.True(t, result[0]["active"].(bool))
	tags := result[0]["tags"].([]interface{})
	assert.Len(t, tags, 2)
	assert.Equal(t, "a", tags[0])
	assert.Equal(t, "b", tags[1])
}
