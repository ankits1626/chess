// Package websocket provides real-time communication for chess games.
package websocket

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// MessageType defines the type of WebSocket message.
type MessageType string

const (
	// TypeRequest is a client request requiring a response
	TypeRequest MessageType = "request"
	// TypeResponse is a server response to a request
	TypeResponse MessageType = "response"
	// TypeEvent is a server-initiated event (no response expected)
	TypeEvent MessageType = "event"
)

// Message is the base WebSocket message structure.
type Message struct {
	ID        string                 `json:"id,omitempty"`        // For request/response matching
	Type      MessageType            `json:"type"`                // request|response|event
	Action    string                 `json:"action,omitempty"`    // Specific action (for requests)
	Event     string                 `json:"event,omitempty"`     // Event name (for events)
	Data      map[string]interface{} `json:"data,omitempty"`      // Payload
	Error     string                 `json:"error,omitempty"`     // Error message (for responses)
	Success   bool                   `json:"success,omitempty"`   // Success flag (for responses)
	Timestamp int64                  `json:"timestamp,omitempty"` // Unix timestamp
}

// NewRequest creates a request message.
func NewRequest(action string, data map[string]interface{}) *Message {
	return &Message{
		ID:        generateID(),
		Type:      TypeRequest,
		Action:    action,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// NewResponse creates a success response message.
func NewResponse(requestID string, data map[string]interface{}) *Message {
	return &Message{
		ID:        requestID,
		Type:      TypeResponse,
		Success:   true,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// NewErrorResponse creates an error response message.
func NewErrorResponse(requestID string, err string) *Message {
	return &Message{
		ID:        requestID,
		Type:      TypeResponse,
		Success:   false,
		Error:     err,
		Timestamp: time.Now().Unix(),
	}
}

// NewEvent creates an event message.
func NewEvent(event string, data map[string]interface{}) *Message {
	return &Message{
		Type:      TypeEvent,
		Event:     event,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// ToJSON converts message to JSON bytes.
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON parses JSON into message.
func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// generateID generates a unique message ID.
func generateID() string {
	return uuid.New().String()
}
