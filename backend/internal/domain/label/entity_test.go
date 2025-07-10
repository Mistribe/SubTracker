package label_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oleexo/subtracker/internal/domain/label"
)

func TestNewLabel_Success(t *testing.T) {
	// Given
	id := uuid.New()
	name := "Test Label"
	isDefault := true
	color := "#FF0000"
	createdAt := time.Now()
	updatedAt := time.Now()

	// When
	lbl := label.NewLabelWithoutValidation(id, name, isDefault, color, createdAt, updatedAt)

	// Then
	require.Nil(t, lbl.Validate(), "Label should be valid")

	assert.Equal(t, id, lbl.Id())
	assert.Equal(t, name, lbl.Name())
	assert.Equal(t, isDefault, lbl.IsDefault())
	assert.Equal(t, color, lbl.Color())
	assert.Equal(t, createdAt, lbl.CreatedAt())
	assert.Equal(t, updatedAt, lbl.UpdatedAt())
}

func TestNewLabel_WithValidColors(t *testing.T) {
	testCases := []struct {
		name  string
		color string
	}{
		{"6-digit hex", "#FF0000"},
		{"3-digit hex", "#F00"},
		{"lowercase hex", "#ff0000"},
		{"mixed case hex", "#Ff0000"},
		{"green color", "#00FF00"},
		{"blue color", "#0000FF"},
		{"white color", "#FFFFFF"},
		{"black color", "#000000"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			id := uuid.New()
			name := "Test Label"
			isDefault := false
			createdAt := time.Now()
			updatedAt := time.Now()

			// When
			lbl := label.NewLabelWithoutValidation(id, name, isDefault, tc.color, createdAt, updatedAt)

			// Then
			assert.Nil(t, lbl.Validate(), "Color %s should be valid", tc.color)
		})
	}
}

func TestNewLabel_WithInvalidColors(t *testing.T) {
	testCases := []struct {
		name  string
		color string
	}{
		{"empty color", ""},
		{"missing hash", "FF0000"},
		{"too short", "#F0"},
		{"too long", "#FF00000"},
		{"invalid characters", "#GGGGGG"},
		{"spaces", "#FF 00 00"},
		{"invalid format", "red"},
		{"hash only", "#"},
		{"4 digits", "#F000"},
		{"5 digits", "#FF000"},
		{"7 digits", "#FF00000"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			id := uuid.New()
			name := "Test Label"
			isDefault := false
			createdAt := time.Now()
			updatedAt := time.Now()

			// When
			lbl := label.NewLabelWithoutValidation(id, name, isDefault, tc.color, createdAt, updatedAt)

			// Then
			assert.NotNil(t, lbl.Validate(), "Color %s should be invalid", tc.color)
			assert.Equal(t, label.ErrLabelColorInvalid, lbl.Validate())
		})
	}
}

func TestNewLabel_WithInvalidNames(t *testing.T) {
	testCases := []struct {
		name        string
		labelName   string
		expectedErr error
	}{
		{"empty name", "", label.ErrLabelNameEmpty},
		{"whitespace only", "   ", label.ErrLabelNameEmpty},
		{"tabs only", "\t\t", label.ErrLabelNameEmpty},
		{"newlines only", "\n\n", label.ErrLabelNameEmpty},
		{"mixed whitespace", " \t\n ", label.ErrLabelNameEmpty},
		{"too long name", string(make([]rune, 101)), label.ErrLabelNameTooLong},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			id := uuid.New()
			isDefault := false
			color := "#FF0000"
			createdAt := time.Now()
			updatedAt := time.Now()

			// When
			lbl := label.NewLabelWithoutValidation(id, tc.labelName, isDefault, color, createdAt, updatedAt)

			// Then
			assert.NotNil(t, lbl.Validate())
			assert.Equal(t, tc.expectedErr, lbl.Validate())
		})
	}
}

func TestNewLabel_WithValidNames(t *testing.T) {
	testCases := []struct {
		name      string
		labelName string
	}{
		{"single character", "A"},
		{"normal name", "Work"},
		{"name with spaces", "Work Label"},
		{"name with numbers", "Label 123"},
		{"name with special chars", "Label @#$%"},
		{"exactly 100 characters", string(make([]rune, 100))},
		{"name with leading/trailing spaces", "  Work  "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			id := uuid.New()
			isDefault := false
			color := "#FF0000"
			createdAt := time.Now()
			updatedAt := time.Now()

			// When
			lbl := label.NewLabelWithoutValidation(id, tc.labelName, isDefault, color, createdAt, updatedAt)

			// Then
			assert.Nil(t, lbl.Validate(), "Name '%s' should be valid", tc.labelName)
		})
	}
}

func TestNewLabelWithoutValidation(t *testing.T) {
	// Given
	id := uuid.New()
	name := "" // Invalid name
	isDefault := true
	color := "invalid" // Invalid color
	createdAt := time.Now()
	updatedAt := time.Now()

	// When
	lbl := label.NewLabelWithoutValidation(id, name, isDefault, color, createdAt, updatedAt)

	// Then
	assert.Equal(t, id, lbl.Id())
	assert.Equal(t, name, lbl.Name())
	assert.Equal(t, isDefault, lbl.IsDefault())
	assert.Equal(t, color, lbl.Color())
	assert.Equal(t, createdAt, lbl.CreatedAt())
	assert.Equal(t, updatedAt, lbl.UpdatedAt())
}

func TestLabel_Getters(t *testing.T) {
	// Given
	id := uuid.New()
	name := "Test Label"
	isDefault := true
	color := "#FF0000"
	createdAt := time.Now()
	updatedAt := time.Now()

	lbl := label.NewLabelWithoutValidation(id, name, isDefault, color, createdAt, updatedAt)

	// When & Then
	assert.Equal(t, id, lbl.Id())
	assert.Equal(t, name, lbl.Name())
	assert.Equal(t, isDefault, lbl.IsDefault())
	assert.Equal(t, color, lbl.Color())
	assert.Equal(t, createdAt, lbl.CreatedAt())
	assert.Equal(t, updatedAt, lbl.UpdatedAt())
}

func TestLabel_SetName(t *testing.T) {
	// Given
	lbl := label.NewLabelWithoutValidation(uuid.New(), "Original", false, "#FF0000", time.Now(), time.Now())
	newName := "New Name"

	// When
	lbl.SetName(newName)

	// Then
	assert.Equal(t, newName, lbl.Name())
}

func TestLabel_SetIsDefault(t *testing.T) {
	// Given
	lbl := label.NewLabelWithoutValidation(uuid.New(), "Test", false, "#FF0000", time.Now(), time.Now())

	// When
	lbl.SetIsDefault(true)

	// Then
	assert.True(t, lbl.IsDefault())

	// When
	lbl.SetIsDefault(false)

	// Then
	assert.False(t, lbl.IsDefault())
}

func TestLabel_SetColor(t *testing.T) {
	// Given
	lbl := label.NewLabelWithoutValidation(uuid.New(), "Test", false, "#FF0000", time.Now(), time.Now())
	newColor := "#00FF00"

	// When
	lbl.SetColor(newColor)

	// Then
	assert.Equal(t, newColor, lbl.Color())
}

func TestLabel_SetUpdatedAt(t *testing.T) {
	// Given
	originalTime := time.Now()
	lbl := label.NewLabelWithoutValidation(uuid.New(), "Test", false, "#FF0000", time.Now(), originalTime)
	newTime := time.Now().Add(time.Hour)

	// When
	lbl.SetUpdatedAt(newTime)

	// Then
	assert.Equal(t, newTime, lbl.UpdatedAt())
}

func TestLabel_Validate_Success(t *testing.T) {
	// Given
	lbl := label.NewLabelWithoutValidation(uuid.New(), "Valid Name", false, "#FF0000", time.Now(), time.Now())

	// When
	err := lbl.Validate()

	// Then
	assert.NoError(t, err)
}

func TestLabel_Validate_NameErrors(t *testing.T) {
	testCases := []struct {
		name        string
		labelName   string
		expectedErr error
	}{
		{"empty name", "", label.ErrLabelNameEmpty},
		{"whitespace only", "   ", label.ErrLabelNameEmpty},
		{"too long name", string(make([]rune, 101)), label.ErrLabelNameTooLong},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Given
			lbl := label.NewLabelWithoutValidation(uuid.New(), tc.labelName, false, "#FF0000", time.Now(), time.Now())

			// When
			err := lbl.Validate()

			// Then
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestLabel_Validate_ColorErrors(t *testing.T) {
	testCases := []string{
		"",
		"invalid",
		"#GG0000",
		"FF0000",
		"#F0",
		"#FF00000",
	}

	for _, color := range testCases {
		t.Run("invalid color: "+color, func(t *testing.T) {
			// Given
			lbl := label.NewLabelWithoutValidation(uuid.New(), "Valid Name", false, color, time.Now(), time.Now())

			// When
			err := lbl.Validate()

			// Then
			assert.Equal(t, label.ErrLabelColorInvalid, err)
		})
	}
}

func TestLabel_Validate_AfterSetters(t *testing.T) {
	// Given
	lbl := label.NewLabelWithoutValidation(uuid.New(), "Valid Name", false, "#FF0000", time.Now(), time.Now())

	// When setting invalid name
	lbl.SetName("")
	err := lbl.Validate()

	// Then
	assert.Equal(t, label.ErrLabelNameEmpty, err)

	// When setting valid name but invalid color
	lbl.SetName("Valid Name")
	lbl.SetColor("invalid")
	err = lbl.Validate()

	// Then
	assert.Equal(t, label.ErrLabelColorInvalid, err)

	// When setting both to valid values
	lbl.SetColor("#00FF00")
	err = lbl.Validate()

	// Then
	assert.NoError(t, err)
}

func TestLabel_EdgeCases(t *testing.T) {
	t.Run("name with exactly 100 characters", func(t *testing.T) {
		// Given
		name := string(make([]rune, 100))
		for i := range name {
			name = name[:i] + "a" + name[i+1:]
		}
		lbl := label.NewLabelWithoutValidation(uuid.New(), name, false, "#FF0000", time.Now(), time.Now())

		// When
		err := lbl.Validate()

		// Then
		assert.NoError(t, err)
	})

	t.Run("name with 101 characters", func(t *testing.T) {
		// Given
		name := string(make([]rune, 101))
		for i := range name {
			name = name[:i] + "a" + name[i+1:]
		}
		lbl := label.NewLabelWithoutValidation(uuid.New(), name, false, "#FF0000", time.Now(), time.Now())

		// When
		err := lbl.Validate()

		// Then
		assert.Equal(t, label.ErrLabelNameTooLong, err)
	})

	t.Run("zero uuid", func(t *testing.T) {
		// Given
		var zeroUUID uuid.UUID
		lbl := label.NewLabelWithoutValidation(zeroUUID, "Valid Name", false, "#FF0000", time.Now(), time.Now())

		// When
		err := lbl.Validate()

		// Then
		assert.NoError(t, err)
		assert.Equal(t, zeroUUID, lbl.Id())
	})

	t.Run("zero times", func(t *testing.T) {
		// Given
		var zeroTime time.Time
		lbl := label.NewLabelWithoutValidation(uuid.New(), "Valid Name", false, "#FF0000", zeroTime, zeroTime)

		// When
		err := lbl.Validate()

		// Then
		assert.NoError(t, err)
		assert.Equal(t, zeroTime, lbl.CreatedAt())
		assert.Equal(t, zeroTime, lbl.UpdatedAt())
	})
}
