package websocket

import "encoding/json"

type MessageType string

const (
    MessageTypePing     MessageType = "ping"
    MessageTypePong     MessageType = "pong"
    MessageTypeJoin     MessageType = "join"
    MessageTypeLeave    MessageType = "leave"
    MessageTypeMove     MessageType = "move"
    MessageTypeAnalysis MessageType = "analysis"
    MessageTypeError    MessageType = "error"
)

type Message struct {
    Type    MessageType     `json:"type"`
    Payload json.RawMessage `json:"payload,omitempty"`
}

type ErrorPayload struct {
    Error string `json:"error"`
}

type MovePayload struct {
    GameID string `json:"game_id"`
    Move   string `json:"move"` // UCI or SAN notation
}
