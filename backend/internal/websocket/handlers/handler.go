package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/ankits1626/chess-coach-backend/internal/websocket"
)

// MessageHandler defines the interface for handling WebSocket messages.
type MessageHandler interface {
	HandleMessage(ctx context.Context, client *websocket.Client, msg *websocket.Message)
}

// GameService defines chess game operations (Dependency Inversion).
type GameService interface {
	// ValidateMove checks if a move is legal
	ValidateMove(fen string, moveSAN string) (MoveResult, error)

	// ApplyMove applies a move and returns new state
	ApplyMove(fen string, moveSAN string) (GameState, error)

	// IsGameOver checks if game is finished
	IsGameOver(fen string) (bool, GameResult)

	// GetStartingPosition returns initial FEN
	GetStartingPosition() string
}

// MoveResult contains validated move information.
type MoveResult struct {
	SAN     string
	UCI     string
	IsLegal bool
}

// GameState contains game state after a move.
type GameState struct {
	FEN        string
	SAN        string
	UCI        string
	MoveNumber int
	IsGameOver bool
	Result     GameResult
}

// GameResult represents game outcome.
type GameResult struct {
	Winner string // "white", "black", "draw", ""
	Method string // "checkmate", "resignation", "stalemate", etc.
	PGN    string
}

// HandlerRouter routes messages to specific handlers.
type HandlerRouter struct {
	gameManager *GameManager
	handlers    map[string]ActionHandler
}

// ActionHandler processes a specific action.
type ActionHandler interface {
	Handle(ctx context.Context, client *websocket.Client, msg *websocket.Message) error
}

// NewHandlerRouter creates a new message router.
func NewHandlerRouter(gameManager *GameManager) *HandlerRouter {
	router := &HandlerRouter{
		gameManager: gameManager,
		handlers:    make(map[string]ActionHandler),
	}

	// Register all handlers
	router.RegisterHandler("createGame", &CreateGameHandler{manager: gameManager})
	router.RegisterHandler("joinGame", &JoinGameHandler{manager: gameManager})
	router.RegisterHandler("makeMove", &MakeMoveHandler{manager: gameManager})
	router.RegisterHandler("resign", &ResignHandler{manager: gameManager})
	router.RegisterHandler("leaveGame", &LeaveGameHandler{manager: gameManager})

	return router
}

// RegisterHandler registers a new action handler.
func (r *HandlerRouter) RegisterHandler(action string, handler ActionHandler) {
	r.handlers[action] = handler
}

// HandleMessage implements MessageHandler interface.
func (r *HandlerRouter) HandleMessage(ctx context.Context, client *websocket.Client, msg *websocket.Message) {
	log.Printf("HandlerRouter: Processing message type=%s action=%s from client=%s",
		msg.Type, msg.Action, client.ID)

	// Only handle request messages
	if msg.Type != websocket.TypeRequest {
		log.Printf("HandlerRouter: Ignoring non-request message type=%s", msg.Type)
		return
	}

	// Find handler for action
	handler, exists := r.handlers[msg.Action]
	if !exists {
		websocket.SendErrorResponse(client, msg.ID, fmt.Sprintf("Unknown action: %s", msg.Action))
		return
	}

	// Execute handler
	if err := handler.Handle(ctx, client, msg); err != nil {
		log.Printf("HandlerRouter: Handler error for action=%s: %v", msg.Action, err)
		websocket.SendErrorResponse(client, msg.ID, err.Error())
	}
}
