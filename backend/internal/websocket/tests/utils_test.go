package tests

import (
	"testing"

	"github.com/ankits1626/chess-coach-backend/internal/websocket"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ============================================================================
// UUID Conversion Tests
// ============================================================================

func TestUUIDToString(t *testing.T) {
	// Test valid UUID
	uuidVal := uuid.New()
	pgUUID := pgtype.UUID{
		Bytes: uuidVal,
		Valid: true,
	}

	str := websocket.UUIDToString(pgUUID)
	if str == "" {
		t.Error("Expected non-empty string for valid UUID")
	}

	if str != uuidVal.String() {
		t.Errorf("Expected %s, got %s", uuidVal.String(), str)
	}

	// Test invalid UUID
	invalidUUID := pgtype.UUID{Valid: false}
	str = websocket.UUIDToString(invalidUUID)
	if str != "" {
		t.Error("Expected empty string for invalid UUID")
	}
}

func TestStringToUUID(t *testing.T) {
	// Test valid UUID string
	uuidStr := "123e4567-e89b-12d3-a456-426614174000"
	pgUUID := websocket.StringToUUID(uuidStr)

	if !pgUUID.Valid {
		t.Error("Expected valid UUID")
	}

	// Convert back and compare
	result := websocket.UUIDToString(pgUUID)
	if result != uuidStr {
		t.Errorf("Expected %s, got %s", uuidStr, result)
	}

	// Test empty string
	pgUUID = websocket.StringToUUID("")
	if pgUUID.Valid {
		t.Error("Expected invalid UUID for empty string")
	}

	// Test invalid string
	pgUUID = websocket.StringToUUID("not-a-uuid")
	if pgUUID.Valid {
		t.Error("Expected invalid UUID for malformed string")
	}
}

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123e4567-e89b-12d3-a456-426614174000", true},
		{"", false},
		{"not-a-uuid", false},
		{"123e4567-e89b-12d3", false},
	}

	for _, tt := range tests {
		result := websocket.IsValidUUID(tt.input)
		if result != tt.expected {
			t.Errorf("IsValidUUID(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

// ============================================================================
// Validation Tests
// ============================================================================

func TestValidateTimeControl(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
		expected    int
	}{
		{"5+0", false, 300},
		{"10+5", false, 600},
		{"3+2", false, 180},
		{"", true, 0},
		{"5", true, 0},
		{"5+", true, 0},
		{"+5", true, 0},
		{"abc+def", true, 0},
		{"0+0", true, 0},
		{"5+-1", true, 0},
	}

	for _, tt := range tests {
		result, err := websocket.ValidateTimeControl(tt.input)

		if tt.expectError {
			if err == nil {
				t.Errorf("ValidateTimeControl(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ValidateTimeControl(%q) unexpected error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ValidateTimeControl(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		}
	}
}

func TestValidateGameResult(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1-0", true},
		{"0-1", true},
		{"1/2-1/2", true},
		{"*", true},
		{"2-0", false},
		{"", false},
		{"draw", false},
	}

	for _, tt := range tests {
		result := websocket.ValidateGameResult(tt.input)
		if result != tt.expected {
			t.Errorf("ValidateGameResult(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestValidateSide(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"white", true},
		{"black", true},
		{"White", false},
		{"BLACK", false},
		{"", false},
		{"red", false},
	}

	for _, tt := range tests {
		result := websocket.ValidateSide(tt.input)
		if result != tt.expected {
			t.Errorf("ValidateSide(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestValidateMoveUCI(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"e2e4", true},
		{"e7e5", true},
		{"e1g1", true},
		{"e7e8q", true},
		{"a1h8r", true},
		{"e2", false},
		{"e2e", false},
		{"e2e4e", false},
		{"e9e4", false},
		{"i2e4", false},
		{"e2e4x", false},
		{"", false},
	}

	for _, tt := range tests {
		result := websocket.ValidateMoveUCI(tt.input)
		if result != tt.expected {
			t.Errorf("ValidateMoveUCI(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestValidateFEN(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", true},
		{"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", true},
		{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR", false}, // Missing parts
		{"", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		result := websocket.ValidateFEN(tt.input)
		if result != tt.expected {
			t.Errorf("ValidateFEN(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

// ============================================================================
// Data Extraction Tests
// ============================================================================

func TestExtractString(t *testing.T) {
	data := map[string]interface{}{
		"name": "Alice",
		"age":  25,
	}

	// Existing string key
	result := websocket.ExtractString(data, "name")
	if result != "Alice" {
		t.Errorf("Expected Alice, got %s", result)
	}

	// Non-existent key
	result = websocket.ExtractString(data, "email")
	if result != "" {
		t.Errorf("Expected empty string, got %s", result)
	}

	// Wrong type
	result = websocket.ExtractString(data, "age")
	if result != "" {
		t.Errorf("Expected empty string for wrong type, got %s", result)
	}
}

func TestRequireString(t *testing.T) {
	data := map[string]interface{}{
		"name":  "Alice",
		"empty": "",
	}

	// Valid string
	result, err := websocket.RequireString(data, "name")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "Alice" {
		t.Errorf("Expected Alice, got %s", result)
	}

	// Missing key
	_, err = websocket.RequireString(data, "missing")
	if err == nil {
		t.Error("Expected error for missing key")
	}

	// Empty string
	_, err = websocket.RequireString(data, "empty")
	if err == nil {
		t.Error("Expected error for empty string")
	}
}
