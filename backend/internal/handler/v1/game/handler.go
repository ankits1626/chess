// Package game handles game-related HTTP requests for API v1.
package game

import (
	"net/http"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Handler handles game endpoints.
type Handler struct {
	repo *repository.GameRepository
}

// NewHandler creates game handler.
func NewHandler(repo *repository.GameRepository) *Handler {
	return &Handler{repo: repo}
}

// List lists games with pagination.
// @Summary List games
// @Description Get paginated list of games
// @Tags games
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} Response
// @Router /games [get]
func (h *Handler) List(c *gin.Context) {
	var params struct {
		Limit  int32 `form:"limit" binding:"omitempty,min=1,max=100"`
		Offset int32 `form:"offset" binding:"omitempty,min=0"`
	}

	// Default limit
	if params.Limit == 0 {
		params.Limit = 10
	}

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	games, err := h.repo.List(c.Request.Context(), params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch games"})
		return
	}

	c.JSON(http.StatusOK, ToResponses(games))
}

// Get retrieves game by ID.
// @Summary Get game
// @Description Get game by ID
// @Tags games
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} Response
// @Router /games/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	game, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	c.JSON(http.StatusOK, ToResponse(game))
}

// ListByUser retrieves user's games.
// @Summary List user's games
// @Description Get paginated list of games for a specific user
// @Tags games
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} Response
// @Router /users/{user_id}/games [get]
func (h *Handler) ListByUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var params struct {
		Limit  int32 `form:"limit" binding:"omitempty,min=1,max=100"`
		Offset int32 `form:"offset" binding:"omitempty,min=0"`
	}

	if params.Limit == 0 {
		params.Limit = 10
	}

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	games, err := h.repo.ListByUser(c.Request.Context(), userID, params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user games"})
		return
	}

	c.JSON(http.StatusOK, ToResponses(games))
}

// Create creates new game.
// @Summary Create game
// @Description Create new chess game
// @Tags games
// @Accept json
// @Produce json
// @Param game body CreateRequest true "Game data"
// @Success 201 {object} Response
// @Router /games [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build create params with type conversions
	params := database.CreateGameParams{
		WhitePlayerID: pgtype.UUID{Bytes: req.WhitePlayerID, Valid: true},
		BlackPlayerID: pgtype.UUID{Bytes: req.BlackPlayerID, Valid: true},
		Pgn:           req.Pgn,
	}

	// Handle optional Result
	if req.Result != "" {
		params.Result = pgtype.Text{String: req.Result, Valid: true}
	}

	// Handle optional TimeControl
	if req.TimeControl > 0 {
		params.TimeControl = pgtype.Int4{Int32: req.TimeControl, Valid: true}
	}

	game, err := h.repo.Create(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create game"})
		return
	}

	c.JSON(http.StatusCreated, ToResponse(game))
}
