package handlers

import (
	"fmt"
	"sync"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/websocket"
)

// GameManager manages game state (pending and active games).
// Single Responsibility: Game lifecycle management
type GameManager struct {
	db           *database.DB
	chessService GameService
	pendingGames map[string]*PendingGame
	activeGames  map[string]*ActiveGame
	mu           sync.RWMutex // Protect concurrent access
}

// PendingGame represents a game waiting for second player.
type PendingGame struct {
	GameID      string
	WhitePlayer *websocket.Client
	TimeControl string
}

// ActiveGame represents an ongoing game with both players.
type ActiveGame struct {
	GameID      string
	WhitePlayer *websocket.Client
	BlackPlayer *websocket.Client
	CurrentFEN  string
	MoveCount   int
	TimeControl string
}

// NewGameManager creates a new game manager.
func NewGameManager(db *database.DB, chessService GameService) *GameManager {
	return &GameManager{
		db:           db,
		chessService: chessService,
		pendingGames: make(map[string]*PendingGame),
		activeGames:  make(map[string]*ActiveGame),
	}
}

// CreatePendingGame adds a new pending game.
func (m *GameManager) CreatePendingGame(gameID string, creator *websocket.Client, timeControl string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.pendingGames[gameID] = &PendingGame{
		GameID:      gameID,
		WhitePlayer: creator,
		TimeControl: timeControl,
	}
}

// GetPendingGame retrieves a pending game.
func (m *GameManager) GetPendingGame(gameID string) (*PendingGame, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, exists := m.pendingGames[gameID]
	return game, exists
}

// RemovePendingGame removes a pending game.
func (m *GameManager) RemovePendingGame(gameID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.pendingGames, gameID)
}

// ActivateGame moves a game from pending to active.
func (m *GameManager) ActivateGame(gameID string, blackPlayer *websocket.Client) (*ActiveGame, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pending, exists := m.pendingGames[gameID]
	if !exists {
		return nil, fmt.Errorf("pending game not found")
	}

	// Create active game
	activeGame := &ActiveGame{
		GameID:      gameID,
		WhitePlayer: pending.WhitePlayer,
		BlackPlayer: blackPlayer,
		CurrentFEN:  m.chessService.GetStartingPosition(),
		MoveCount:   0,
		TimeControl: pending.TimeControl,
	}

	m.activeGames[gameID] = activeGame
	delete(m.pendingGames, gameID)

	return activeGame, nil
}

// GetActiveGame retrieves an active game.
func (m *GameManager) GetActiveGame(gameID string) (*ActiveGame, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, exists := m.activeGames[gameID]
	return game, exists
}

// UpdateGameState updates the FEN and move count.
func (m *GameManager) UpdateGameState(gameID string, newFEN string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	game, exists := m.activeGames[gameID]
	if !exists {
		return fmt.Errorf("active game not found")
	}

	game.CurrentFEN = newFEN
	game.MoveCount++

	return nil
}

// RemoveActiveGame removes an active game.
func (m *GameManager) RemoveActiveGame(gameID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.activeGames, gameID)
}

// GetPlayerSide returns the side (white/black) for a client in a game.
func (m *GameManager) GetPlayerSide(gameID string, client *websocket.Client) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, exists := m.activeGames[gameID]
	if !exists {
		return "", fmt.Errorf("game not found")
	}

	if game.WhitePlayer.ID == client.ID {
		return "white", nil
	}
	if game.BlackPlayer.ID == client.ID {
		return "black", nil
	}

	return "", fmt.Errorf("client not in game")
}

// GetOpponent returns the opponent client.
func (m *GameManager) GetOpponent(gameID string, client *websocket.Client) (*websocket.Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, exists := m.activeGames[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}

	if game.WhitePlayer.ID == client.ID {
		return game.BlackPlayer, nil
	}
	if game.BlackPlayer.ID == client.ID {
		return game.WhitePlayer, nil
	}

	return nil, fmt.Errorf("client not in game")
}

// GetDB returns the database connection.
func (m *GameManager) GetDB() *database.DB {
	return m.db
}

// GetChessService returns the chess service.
func (m *GameManager) GetChessService() GameService {
	return m.chessService
}
