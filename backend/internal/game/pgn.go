package game

import (
    "fmt"
    "strings"

    "github.com/notnil/chess"
)

type PGNParser struct{}

func NewPGNParser() *PGNParser {
    return &PGNParser{}
}

func (p *PGNParser) Parse(pgnString string) (*chess.Game, error) {
    reader := strings.NewReader(pgnString)
    pgn, err := chess.PGN(reader)
    if err != nil {
        return nil, fmt.Errorf("invalid PGN: %w", err)
    }

    game := chess.NewGame(pgn)
    return game, nil
}

func (p *PGNParser) Validate(pgnString string) error {
    _, err := p.Parse(pgnString)
    return err
}

func (p *PGNParser) GetMoves(game *chess.Game) []string {
    moves := []string{}
    notation := chess.AlgebraicNotation{}

    // Replay game from start to get moves in SAN format
    tempGame := chess.NewGame()
    for _, move := range game.Moves() {
        // Encode move with current position (before move is played)
        moveStr := notation.Encode(tempGame.Position(), move)
        moves = append(moves, moveStr)

        // Play the move to advance position
        tempGame.Move(move)
    }

    return moves
}

func (p *PGNParser) GetFEN(game *chess.Game) string {
    return game.Position().String()
}
