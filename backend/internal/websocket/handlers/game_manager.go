package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/websocket"
	"github.com/ankits1626/chess-coach-backend/internal/websocket/handlers/player"
)

// GameManager manages game state (pending and active games).
// Single Responsibility: Game lifecycle management
type GameManager struct {
	db            *database.DB
	chessService  GameService
	aiService     player.AIService
	playerFactory *player.PlayerFactory
	pendingGames  map[string]*PendingGame
	activeGames   map[string]*ActiveGame
	mu            sync.RWMutex // Protect concurrent access
}

// PendingGame represents a game waiting for second player.
type PendingGame struct {
	GameID      string
	Mode        player.GameMode
	WhitePlayer player.Player
	BlackType   player.PlayerType // Type of black player (Human or Computer)
	Difficulty  string
	TimeControl string
}

// ActiveGame represents an ongoing game with both players.
type ActiveGame struct {
	GameID      string
	Mode        player.GameMode
	WhitePlayer player.Player
	BlackPlayer player.Player
	CurrentFEN  string
	MoveCount   int
	TimeControl string
}

// NewGameManager creates a new game manager.
func NewGameManager(db *database.DB, chessService GameService, aiService player.AIService) *GameManager {
	return &GameManager{
		db:            db,
		chessService:  chessService,
		aiService:     aiService,
		playerFactory: player.NewPlayerFactory(aiService),
		pendingGames:  make(map[string]*PendingGame),
		activeGames:   make(map[string]*ActiveGame),
	}
}

// CreatePendingGame adds a new pending game.
func (m *GameManager) CreatePendingGame(
	gameID string,
	mode player.GameMode,
	creator *websocket.Client,
	difficulty string,
	playerColor string,
	timeControl string,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Determine player types based on mode and chosen color
	whiteType, blackType := mode.GetPlayerTypes()

	// If playing as black in human vs computer, swap the types
	if mode == player.GameModeHumanVsComputer && playerColor == "black" {
		whiteType = player.PlayerTypeComputer
		blackType = player.PlayerTypeHuman
	}

	// Create white player
	var whitePlayer player.Player
	var err error

	if whiteType == player.PlayerTypeHuman {
		// Human is white
		config := player.HumanPlayerConfig{Client: creator}
		whitePlayer, err = m.playerFactory.CreatePlayer(config)
		if err != nil {
			return fmt.Errorf("failed to create white player: %w", err)
		}
	} else if whiteType == player.PlayerTypeComputer {
		// Computer is white
		config := player.ComputerPlayerConfig{Difficulty: difficulty}
		whitePlayer, err = m.playerFactory.CreatePlayer(config)
		if err != nil {
			return fmt.Errorf("failed to create white player: %w", err)
		}
	}

	m.pendingGames[gameID] = &PendingGame{
		GameID:      gameID,
		Mode:        mode,
		WhitePlayer: whitePlayer,
		BlackType:   blackType,
		Difficulty:  difficulty,
		TimeControl: timeControl,
	}

	return nil
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
func (m *GameManager) ActivateGame(gameID string, blackClient *websocket.Client) (*ActiveGame, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pending, exists := m.pendingGames[gameID]
	if !exists {
		return nil, fmt.Errorf("pending game not found")
	}

	// Use the black player type from pending game (already determined by CreatePendingGame)
	blackType := pending.BlackType

	// Create black player based on type
	var blackPlayer player.Player
	var err error

	if blackType == player.PlayerTypeHuman {
		if blackClient == nil {
			return nil, fmt.Errorf("black client required for human player")
		}
		config := player.HumanPlayerConfig{Client: blackClient}
		blackPlayer, err = m.playerFactory.CreatePlayer(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create black player: %w", err)
		}
	} else if blackType == player.PlayerTypeComputer {
		config := player.ComputerPlayerConfig{Difficulty: pending.Difficulty}
		blackPlayer, err = m.playerFactory.CreatePlayer(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create black player: %w", err)
		}
	}

	activeGame := &ActiveGame{
		GameID:      gameID,
		Mode:        pending.Mode,
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

	if game.WhitePlayer.GetID() == client.ID {
		return "white", nil
	}
	if game.BlackPlayer.GetID() == client.ID {
		return "black", nil
	}

	return "", fmt.Errorf("client not in game")
}

// GetOpponent returns the opponent client.
func (m *GameManager) GetOpponent(gameID string, client *websocket.Client) (player.Player, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, exists := m.activeGames[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}

	if game.WhitePlayer.GetID() == client.ID {
		return game.BlackPlayer, nil
	}
	if game.BlackPlayer.GetID() == client.ID {
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

// GetPlayer returns the player for a given side
func (m *GameManager) GetPlayer(gameID string, side string) (player.Player, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, exists := m.activeGames[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}

	if side == "white" {
		return game.WhitePlayer, nil
	} else if side == "black" {
		return game.BlackPlayer, nil
	}

	return nil, fmt.Errorf("invalid side: %s", side)
}

// GetGameMode returns the mode for a game
func (m *GameManager) GetGameMode(gameID string) (player.GameMode, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, exists := m.activeGames[gameID]
	if !exists {
		return "", fmt.Errorf("game not found")
	}

	return game.Mode, nil
}

// IsComputerTurn checks if it's a computer player's turn
func (m *GameManager) IsComputerTurn(gameID string, fen string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, exists := m.activeGames[gameID]
	if !exists {
		return false
	}

	// Parse FEN to determine whose turn it is
	// FEN format: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	// The field after the position is the active color: 'w' or 'b'
	fields := strings.Fields(fen)
	if len(fields) < 2 {
		return false
	}

	activeColor := fields[1]
	if activeColor == "w" && game.WhitePlayer.GetType() == player.PlayerTypeComputer {
		return true
	}
	if activeColor == "b" && game.BlackPlayer.GetType() == player.PlayerTypeComputer {
		return true
	}

	return false
}

// HandleComputerMove triggers the computer to make a move if it's their turn
func (m *GameManager) HandleComputerMove(ctx context.Context, gameID string) error {
	m.mu.RLock()
	game, exists := m.activeGames[gameID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("game not found")
	}

	// Check if it's computer's turn
	if !m.IsComputerTurn(gameID, game.CurrentFEN) {
		return nil // Not computer's turn, nothing to do
	}

	// Determine which player is the computer
	var computerPlayer player.Player
	var side string

	if game.WhitePlayer.GetType() == player.PlayerTypeComputer {
		computerPlayer = game.WhitePlayer
		side = "white"
	} else {
		computerPlayer = game.BlackPlayer
		side = "black"
	}

	// Request move from computer
	moveUCI, err := computerPlayer.RequestMove(ctx, game.CurrentFEN)
	if err != nil {
		log.Printf("Computer move failed for game %s: %v", gameID, err)
		return fmt.Errorf("computer move failed: %w", err)
	}

	log.Printf("Computer (%s) in game %s plays: %s", side, gameID, moveUCI)

	// Apply the move
	gameState, err := m.chessService.ApplyMove(game.CurrentFEN, moveUCI)
	if err != nil {
		log.Printf("Invalid computer move %s for game %s: %v", moveUCI, gameID, err)
		return fmt.Errorf("invalid computer move: %w", err)
	}

	// Update game state
	if err := m.UpdateGameState(gameID, gameState.FEN); err != nil {
		return fmt.Errorf("failed to update game state: %w", err)
	}

	// Check if game is over
	isGameOver, gameResult := m.chessService.IsGameOver(gameState.FEN)
	if isGameOver {
		log.Printf("Game %s ended: %s", gameID, gameResult.Method)
	}

	// Convert to player.GameResult
	var playerResult *player.GameResult
	if isGameOver {
		playerResult = &player.GameResult{
			Winner: gameResult.Winner,
			Method: gameResult.Method,
			PGN:    gameResult.PGN,
		}
	}

	// Notify both players
	moveNotification := player.MoveNotification{
		GameID:     gameID,
		MoveSAN:    gameState.SAN,
		MoveUCI:    gameState.UCI,
		FEN:        gameState.FEN,
		MoveNumber: game.MoveCount,
		IsGameOver: isGameOver,
		Result:     playerResult,
	}

	// Send to white player
	if err := game.WhitePlayer.SendMove(moveNotification); err != nil {
		log.Printf("Failed to send move to white player: %v", err)
	}

	// Send to black player
	if err := game.BlackPlayer.SendMove(moveNotification); err != nil {
		log.Printf("Failed to send move to black player: %v", err)
	}

	// If game is over, clean up
	if isGameOver {
		m.RemoveActiveGame(gameID)
	}

	return nil
}
