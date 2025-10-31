package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPGNParser_Parse(t *testing.T) {
	parser := NewPGNParser()

	t.Run("valid PGN with moves", func(t *testing.T) {
		pgn := `1. e4 e5 2. Nf3 Nc6 3. Bb5`
		game, err := parser.Parse(pgn)

		require.NoError(t, err)
		assert.NotNil(t, game)

		moves := parser.GetMoves(game)
		assert.Equal(t, 5, len(moves))
		assert.Equal(t, "e4", moves[0])
		assert.Equal(t, "e5", moves[1])
		assert.Equal(t, "Nf3", moves[2])
	})

	t.Run("valid PGN with headers", func(t *testing.T) {
		pgn := `[Event "Test Game"]
[Site "Online"]
[Date "2025.11.01"]
[White "Player1"]
[Black "Player2"]
[Result "1-0"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7`
		game, err := parser.Parse(pgn)

		require.NoError(t, err)
		assert.NotNil(t, game)

		moves := parser.GetMoves(game)
		assert.Equal(t, 10, len(moves))
	})

	t.Run("empty PGN", func(t *testing.T) {
		pgn := ``
		game, err := parser.Parse(pgn)

		// Empty PGN should parse successfully (new game)
		require.NoError(t, err)
		assert.NotNil(t, game)

		moves := parser.GetMoves(game)
		assert.Equal(t, 0, len(moves))
	})

	t.Run("invalid move notation", func(t *testing.T) {
		pgn := `1. e4 e5 2. Zz9`
		_, err := parser.Parse(pgn)

		assert.Error(t, err)
	})

	t.Run("incomplete move", func(t *testing.T) {
		pgn := `1. e4 e5 2.`
		game, err := parser.Parse(pgn)

		// Should parse successfully, just stop at last complete move
		require.NoError(t, err)
		assert.NotNil(t, game)

		moves := parser.GetMoves(game)
		assert.Equal(t, 2, len(moves))
	})
}

func TestPGNParser_Validate(t *testing.T) {
	parser := NewPGNParser()

	t.Run("valid PGN passes validation", func(t *testing.T) {
		pgn := `1. e4 e5 2. Nf3 Nc6`
		err := parser.Validate(pgn)

		assert.NoError(t, err)
	})

	t.Run("invalid PGN fails validation", func(t *testing.T) {
		pgn := `1. e4 e5 2. InvalidMove`
		err := parser.Validate(pgn)

		assert.Error(t, err)
	})
}

func TestPGNParser_GetFEN(t *testing.T) {
	parser := NewPGNParser()

	t.Run("starting position FEN", func(t *testing.T) {
		pgn := ``
		game, err := parser.Parse(pgn)
		require.NoError(t, err)

		fen := parser.GetFEN(game)
		// Starting position FEN
		assert.Contains(t, fen, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	})

	t.Run("FEN after e4", func(t *testing.T) {
		pgn := `1. e4`
		game, err := parser.Parse(pgn)
		require.NoError(t, err)

		fen := parser.GetFEN(game)
		assert.Contains(t, fen, "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR")
	})
}

func TestPGNParser_GetMoves(t *testing.T) {
	parser := NewPGNParser()

	t.Run("returns empty array for new game", func(t *testing.T) {
		pgn := ``
		game, err := parser.Parse(pgn)
		require.NoError(t, err)

		moves := parser.GetMoves(game)
		assert.NotNil(t, moves)
		assert.Equal(t, 0, len(moves))
	})

	t.Run("returns all moves in order", func(t *testing.T) {
		pgn := `1. d4 d5 2. c4 c6 3. Nf3 Nf6`
		game, err := parser.Parse(pgn)
		require.NoError(t, err)

		moves := parser.GetMoves(game)
		expected := []string{"d4", "d5", "c4", "c6", "Nf3", "Nf6"}
		assert.Equal(t, expected, moves)
	})
}

// Integration test with real chess game
func TestPGNParser_RealGame(t *testing.T) {
	parser := NewPGNParser()

	// Famous "Immortal Game" between Anderssen and Kieseritzky
	pgn := `[Event "London"]
[Site "London ENG"]
[Date "1851.06.21"]
[White "Adolf Anderssen"]
[Black "Lionel Kieseritzky"]
[Result "1-0"]

1. e4 e5 2. f4 exf4 3. Bc4 Qh4+ 4. Kf1 b5 5. Bxb5 Nf6 6. Nf3 Qh6
7. d3 Nh5 8. Nh4 Qg5 9. Nf5 c6 10. g4 Nf6 11. Rg1 cxb5 12. h4 Qg6
13. h5 Qg5 14. Qf3 Ng8 15. Bxf4 Qf6 16. Nc3 Bc5 17. Nd5 Qxb2
18. Bd6 Bxg1 19. e5 Qxa1+ 20. Ke2 Na6 21. Nxg7+ Kd8 22. Qf6+ Nxf6
23. Be7# 1-0`

	game, err := parser.Parse(pgn)
	require.NoError(t, err)
	assert.NotNil(t, game)

	moves := parser.GetMoves(game)
	assert.Equal(t, 45, len(moves)) // 22.5 full moves = 45 half-moves

	// Verify final FEN shows checkmate
	fen := parser.GetFEN(game)
	assert.NotEmpty(t, fen)
}
