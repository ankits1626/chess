package websocket

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ============================================================================
// UUID Conversion Utilities
// ============================================================================

// UUIDToString converts pgtype.UUID to string.
// Returns empty string if UUID is not valid.
func UUIDToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}

	// pgtype.UUID stores as [16]byte
	uuidVal, err := uuid.FromBytes(u.Bytes[:])
	if err != nil {
		return ""
	}

	return uuidVal.String()
}

// StringToUUID converts string to pgtype.UUID.
// Returns invalid UUID if string is empty or malformed.
func StringToUUID(s string) pgtype.UUID {
	if s == "" {
		return pgtype.UUID{Valid: false}
	}

	uuidVal, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{Valid: false}
	}

	return pgtype.UUID{
		Bytes: uuidVal,
		Valid: true,
	}
}

// GenerateUUID creates a new random UUID as pgtype.UUID.
func GenerateUUID() pgtype.UUID {
	uuidVal := uuid.New()
	return pgtype.UUID{
		Bytes: uuidVal,
		Valid: true,
	}
}

// IsValidUUID checks if a string is a valid UUID format.
func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

// ============================================================================
// Optional Type Utilities
// ============================================================================

// StringToText converts string to pgtype.Text.
// Empty string becomes invalid Text.
func StringToText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}

// TextToString converts pgtype.Text to string.
// Returns empty string if Text is not valid.
func TextToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// IntToInt4 converts int to pgtype.Int4.
func IntToInt4(i int) pgtype.Int4 {
	return pgtype.Int4{
		Int32: int32(i),
		Valid: true,
	}
}

// Int4ToInt converts pgtype.Int4 to int.
// Returns 0 if Int4 is not valid.
func Int4ToInt(i pgtype.Int4) int {
	if !i.Valid {
		return 0
	}
	return int(i.Int32)
}

// ============================================================================
// Validation Utilities
// ============================================================================

// ValidateTimeControl checks if time control string is valid.
// Accepts formats: "5+0", "10+5", "3+2", etc.
// Returns time in seconds or error.
func ValidateTimeControl(tc string) (int, error) {
	if tc == "" {
		return 0, fmt.Errorf("time control cannot be empty")
	}

	parts := strings.Split(tc, "+")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time control format, expected 'minutes+increment'")
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes: %w", err)
	}

	increment, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid increment: %w", err)
	}

	if minutes <= 0 {
		return 0, fmt.Errorf("minutes must be positive")
	}

	if increment < 0 {
		return 0, fmt.Errorf("increment cannot be negative")
	}

	// Return total seconds for base time
	return minutes * 60, nil
}

// ValidateGameResult checks if game result is valid.
// Valid results: "1-0", "0-1", "1/2-1/2", "*"
func ValidateGameResult(result string) bool {
	switch result {
	case "1-0", "0-1", "1/2-1/2", "*":
		return true
	default:
		return false
	}
}

// ValidateSide checks if side is "white" or "black".
func ValidateSide(side string) bool {
	return side == "white" || side == "black"
}

// ValidateMoveSAN performs basic validation on SAN notation.
// More thorough validation should use a chess library.
func ValidateMoveSAN(san string) bool {
	if san == "" {
		return false
	}
	// Basic check: SAN should be 2-20 characters
	// More thorough validation in Phase 4 with chess library
	return len(san) >= 2 && len(san) <= 20
}

// ValidateMoveUCI performs basic validation on UCI notation.
// Format: e2e4, e7e5, e1g1 (castling), e7e8q (promotion)
func ValidateMoveUCI(uci string) bool {
	if len(uci) < 4 || len(uci) > 5 {
		return false
	}

	// Check format: [a-h][1-8][a-h][1-8][qrbn]?
	if uci[0] < 'a' || uci[0] > 'h' {
		return false
	}
	if uci[1] < '1' || uci[1] > '8' {
		return false
	}
	if uci[2] < 'a' || uci[2] > 'h' {
		return false
	}
	if uci[3] < '1' || uci[3] > '8' {
		return false
	}

	// Check promotion piece if present
	if len(uci) == 5 {
		promotion := uci[4]
		if promotion != 'q' && promotion != 'r' && promotion != 'b' && promotion != 'n' {
			return false
		}
	}

	return true
}

// ValidateFEN performs basic validation on FEN string.
// A complete FEN has 6 parts separated by spaces.
func ValidateFEN(fen string) bool {
	if fen == "" {
		return false
	}

	parts := strings.Split(fen, " ")
	// Complete FEN: position activeColor castling enPassant halfmove fullmove
	if len(parts) != 6 {
		return false
	}

	// Check position part has 8 ranks
	ranks := strings.Split(parts[0], "/")
	if len(ranks) != 8 {
		return false
	}

	return true
}

// ============================================================================
// Error Response Utilities
// ============================================================================

// ValidationError represents a validation error with details.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements error interface.
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error.
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// SendErrorResponse sends an error response to the client.
// Convenience function to avoid repeating error response code.
func SendErrorResponse(client *Client, requestID, errorMsg string) {
	resp := NewErrorResponse(requestID, errorMsg)
	if err := client.SendMessage(resp); err != nil {
		// Log but don't fail - client might already be disconnected
		fmt.Printf("Failed to send error response to client %s: %v\n", client.ID, err)
	}
}

// SendSuccessResponse sends a success response to the client.
// Convenience function for common success patterns.
func SendSuccessResponse(client *Client, requestID string, data map[string]interface{}) {
	resp := NewResponse(requestID, data)
	if err := client.SendMessage(resp); err != nil {
		fmt.Printf("Failed to send success response to client %s: %v\n", client.ID, err)
	}
}

// ============================================================================
// Message Data Extraction Utilities
// ============================================================================

// ExtractString safely extracts a string from message data.
// Returns empty string if key doesn't exist or type is wrong.
func ExtractString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// ExtractInt safely extracts an int from message data.
// Returns 0 if key doesn't exist or type is wrong.
func ExtractInt(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		// Handle both float64 (from JSON) and int
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		}
	}
	return 0
}

// ExtractBool safely extracts a bool from message data.
// Returns false if key doesn't exist or type is wrong.
func ExtractBool(data map[string]interface{}, key string) bool {
	if val, ok := data[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

// RequireString extracts a required string field.
// Returns error if field is missing or empty.
func RequireString(data map[string]interface{}, key string) (string, error) {
	val := ExtractString(data, key)
	if val == "" {
		return "", NewValidationError(key, "required field missing or empty")
	}
	return val, nil
}

// RequireInt extracts a required int field.
// Returns error if field is missing or zero.
func RequireInt(data map[string]interface{}, key string) (int, error) {
	val := ExtractInt(data, key)
	if val == 0 {
		return 0, NewValidationError(key, "required field missing or zero")
	}
	return val, nil
}
