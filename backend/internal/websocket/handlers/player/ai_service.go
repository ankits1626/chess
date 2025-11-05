package player

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

// AIService provides AI move calculation using chess engines
type AIService interface {
	GetBestMove(ctx context.Context, game *chess.Game, difficulty string) (*chess.Move, error)
	Close() error
}

// StockfishService implements AIService using Stockfish chess engine
type StockfishService struct {
	engine     *uci.Engine
	enginePath string
}

// NewStockfishService creates a new Stockfish AI service
func NewStockfishService() (*StockfishService, error) {
	// Get Stockfish path from environment
	stockfishPath := os.Getenv("STOCKFISH_PATH")

	// Fallback to common paths if not set
	if stockfishPath == "" {
		paths := []string{
			"/usr/local/bin/stockfish",    // Our Docker build location
			"/usr/games/stockfish",        // Alpine package location
			"/opt/homebrew/bin/stockfish", // macOS Apple Silicon
			"/usr/local/bin/stockfish",    // macOS Intel
			"/usr/bin/stockfish",          // Linux
		}

		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				stockfishPath = path
				break
			}
		}

		if stockfishPath == "" {
			return nil, fmt.Errorf("stockfish not found, set STOCKFISH_PATH environment variable")
		}
	}

	// Create UCI engine
	eng, err := uci.New(stockfishPath)
	if err != nil {
		return nil, fmt.Errorf("failed to start stockfish at %s: %w", stockfishPath, err)
	}

	// Initialize engine
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		return nil, fmt.Errorf("failed to initialize stockfish: %w", err)
	}

	return &StockfishService{
		engine:     eng,
		enginePath: stockfishPath,
	}, nil
}

// GetBestMove calculates the best move for the current position
func (s *StockfishService) GetBestMove(ctx context.Context, game *chess.Game, difficulty string) (*chess.Move, error) {
	// Set position
	cmdPos := uci.CmdPosition{Position: game.Position()}
	cmdGo := uci.CmdGo{MoveTime: s.getMoveTime(difficulty)}

	// Run calculation
	if err := s.engine.Run(cmdPos, cmdGo); err != nil {
		return nil, fmt.Errorf("stockfish calculation failed: %w", err)
	}

	// Get best move
	searchResults := s.engine.SearchResults()
	if searchResults.BestMove == nil {
		return nil, fmt.Errorf("stockfish returned no move")
	}

	// Return the move directly - it's already a *chess.Move
	return searchResults.BestMove, nil
}

// getMoveTime returns the thinking time based on difficulty
func (s *StockfishService) getMoveTime(difficulty string) time.Duration {
	switch difficulty {
	case "easy":
		return 100 * time.Millisecond // Quick, weaker moves
	case "hard":
		return 2 * time.Second // Stronger, more calculated
	default: // "medium"
		return 500 * time.Millisecond // Balanced
	}
}

// Close stops the Stockfish engine
func (s *StockfishService) Close() error {
	if s.engine != nil {
		return s.engine.Close()
	}
	return nil
}
