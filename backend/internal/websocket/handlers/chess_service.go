package handlers

import (
	"fmt"

	"github.com/notnil/chess"
)

// ChessService implements GameService using notnil/chess library.
type ChessService struct{}

// NewChessService creates a new chess service.
func NewChessService() *ChessService {
	return &ChessService{}
}

// ValidateMove implements GameService.ValidateMove.
func (s *ChessService) ValidateMove(fen string, moveSAN string) (MoveResult, error) {
	// Parse FEN
	fenFunc, err := chess.FEN(fen)
	if err != nil {
		return MoveResult{}, fmt.Errorf("invalid FEN: %w", err)
	}

	game := chess.NewGame(fenFunc)

	// Try to parse move as UCI first
	move, err := chess.UCINotation{}.Decode(game.Position(), moveSAN)
	if err != nil {
		// Try as SAN
		move, err = chess.AlgebraicNotation{}.Decode(game.Position(), moveSAN)
		if err != nil {
			return MoveResult{IsLegal: false}, fmt.Errorf("invalid move notation: %w", err)
		}
	}

	// Check if move is legal
	validMoves := game.ValidMoves()
	isLegal := false
	for _, validMove := range validMoves {
		if validMove == move {
			isLegal = true
			break
		}
	}

	if !isLegal {
		return MoveResult{IsLegal: false}, fmt.Errorf("illegal move")
	}

	// Get UCI and SAN notation
	uci := chess.UCINotation{}.Encode(game.Position(), move)
	san := chess.AlgebraicNotation{}.Encode(game.Position(), move)

	return MoveResult{
		SAN:     san,
		UCI:     uci,
		IsLegal: true,
	}, nil
}

// ApplyMove implements GameService.ApplyMove.
func (s *ChessService) ApplyMove(fen string, moveSAN string) (GameState, error) {
	// Parse FEN
	fenFunc, err := chess.FEN(fen)
	if err != nil {
		return GameState{}, fmt.Errorf("invalid FEN: %w", err)
	}

	game := chess.NewGame(fenFunc)

	// Parse and apply move
	move, err := chess.UCINotation{}.Decode(game.Position(), moveSAN)
	if err != nil {
		move, err = chess.AlgebraicNotation{}.Decode(game.Position(), moveSAN)
		if err != nil {
			return GameState{}, fmt.Errorf("invalid move: %w", err)
		}
	}

	if err := game.Move(move); err != nil {
		return GameState{}, fmt.Errorf("failed to apply move: %w", err)
	}

	// Get new state
	position := game.Position()
	newFEN := position.String()
	moveNumber := len(game.Moves())

	// Get move notation
	uci := chess.UCINotation{}.Encode(position, move)
	san := chess.AlgebraicNotation{}.Encode(position, move)

	// Check if game is over
	isGameOver, result := s.IsGameOver(newFEN)

	return GameState{
		FEN:        newFEN,
		SAN:        san,
		UCI:        uci,
		MoveNumber: moveNumber,
		IsGameOver: isGameOver,
		Result:     result,
	}, nil
}

// IsGameOver implements GameService.IsGameOver.
func (s *ChessService) IsGameOver(fen string) (bool, GameResult) {
	fenFunc, err := chess.FEN(fen)
	if err != nil {
		return false, GameResult{}
	}

	game := chess.NewGame(fenFunc)
	outcome := game.Outcome()

	if outcome == chess.NoOutcome {
		return false, GameResult{}
	}

	result := GameResult{
		Method: game.Method().String(),
		PGN:    game.String(),
	}

	switch outcome {
	case chess.WhiteWon:
		result.Winner = "white"
	case chess.BlackWon:
		result.Winner = "black"
	case chess.Draw:
		result.Winner = "draw"
	}

	return true, result
}

// GetStartingPosition implements GameService.GetStartingPosition.
func (s *ChessService) GetStartingPosition() string {
	// Standard starting position FEN
	return "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
}
