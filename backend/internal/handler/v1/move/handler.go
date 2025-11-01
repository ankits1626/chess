// Package move handles move-related HTTP requests for API v1.
package move

import (
	"net/http"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Handler handles move endpoints.
type Handler struct {
	repo *repository.MoveRepository
}

// NewHandler creates move handler.
func NewHandler(repo *repository.MoveRepository) *Handler {
	return &Handler{repo: repo}
}

// ListByGame retrieves all moves for a game.
// @Summary List game moves
// @Description Get all moves for a specific game, ordered by move number
// @Tags moves
// @Accept json
// @Produce json
// @Param game_id path string true "Game ID"
// @Success 200 {array} Response
// @Router /games/{game_id}/moves [get]
func (h *Handler) ListByGame(c *gin.Context) {
	gameID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	moves, err := h.repo.ListByGame(c.Request.Context(), gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch moves"})
		return
	}

	c.JSON(http.StatusOK, ToResponses(moves))
}

// Get retrieves move by ID.
// @Summary Get move
// @Description Get move by ID
// @Tags moves
// @Accept json
// @Produce json
// @Param id path string true "Move ID"
// @Success 200 {object} Response
// @Router /moves/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid move ID"})
		return
	}

	move, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "move not found"})
		return
	}

	c.JSON(http.StatusOK, ToResponse(move))
}

// Create creates new move for a game.
// @Summary Create move
// @Description Add new move to a game
// @Tags moves
// @Accept json
// @Produce json
// @Param game_id path string true "Game ID"
// @Param move body CreateRequest true "Move data"
// @Success 201 {object} Response
// @Router /games/{game_id}/moves [post]
func (h *Handler) Create(c *gin.Context) {
	gameID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build create params with type conversions
	params := database.CreateMoveParams{
		GameID:     pgtype.UUID{Bytes: gameID, Valid: true},
		MoveNumber: req.MoveNumber,
		Side:       req.Side,
		MoveSan:    req.MoveSan,
		MoveUci:    req.MoveUci,
		Fen:        req.Fen,
	}

	// Handle optional TimeTaken
	if req.TimeTaken != nil && *req.TimeTaken > 0 {
		params.TimeTaken = pgtype.Int4{Int32: *req.TimeTaken, Valid: true}
	}

	move, err := h.repo.Create(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create move"})
		return
	}

	c.JSON(http.StatusCreated, ToResponse(move))
}
