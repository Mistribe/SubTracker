package ginx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/mistribe/subtracker/pkg/x/exception"
)

type samplePayload struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func TestFromError_ValidationErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	v := validator.New()

	payload := samplePayload{} // Missing both fields
	err := v.Struct(payload)
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/test", nil)

	FromError(c, err)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	var problem HttpErrorResponse
	if unmarshalErr := json.Unmarshal(rec.Body.Bytes(), &problem); unmarshalErr != nil {
		t.Fatalf("failed to unmarshal problem response: %v", unmarshalErr)
	}

	if problem.Status != http.StatusBadRequest {
		t.Errorf("problem status mismatch: expected 400 got %d", problem.Status)
	}
	if problem.Title != http.StatusText(http.StatusBadRequest) {
		t.Errorf("problem title mismatch: expected %s got %s", http.StatusText(http.StatusBadRequest), problem.Title)
	}
	// Expect both field names to appear in aggregated detail
	if problem.Detail == "" || !strings.Contains(problem.Detail, "Name") || !strings.Contains(problem.Detail, "Email") {
		t.Errorf("expected detail to contain field names, got: %s", problem.Detail)
	}
}

func TestFromError_DomainException(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/test-resource", nil)

	ex := exception.NewNotFound("entity not found")
	FromError(c, ex)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}

	var problem HttpErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &problem); err != nil {
		t.Fatalf("failed to unmarshal problem: %v", err)
	}

	if problem.Detail != "entity not found" {
		t.Errorf("expected detail 'entity not found', got %s", problem.Detail)
	}
	if problem.Instance != "/test-resource" {
		// Should fall back to request URL path when no full route registered
		t.Errorf("expected instance '/test-resource', got '%s'", problem.Instance)
	}
}
