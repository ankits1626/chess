// Package user handles user-related HTTP requests for API v1.
package user

import (
	"net/http"

	"github.com/ankits1626/chess-coach-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles user endpoints.
type Handler struct {
	repo *repository.UserRepository
}

// NewHandler creates user handler.
func NewHandler(repo *repository.UserRepository) *Handler {
	return &Handler{repo: repo}
}

// List lists users with pagination.
// @Summary List users
// @Description Get paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} Response
// @Router /users [get]
func (h *Handler) List(c *gin.Context) {
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

	users, err := h.repo.List(c.Request.Context(), params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, ToResponses(users))
}

// Get retrieves user by ID.
// @Summary Get user
// @Description Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} Response
// @Router /users/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, ToResponse(user))
}

// Create creates new user.
// @Summary Create user
// @Description Create new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateRequest true "User data"
// @Success 201 {object} Response
// @Router /users [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Rating == 0 {
		req.Rating = 1200 // Default rating
	}

	user, err := h.repo.Create(c.Request.Context(), req.Username, req.Email, req.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, ToResponse(user))
}

// Update updates existing user.
// @Summary Update user
// @Description Update user details
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body UpdateRequest true "User data"
// @Success 200 {object} Response
// @Router /users/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.Update(c.Request.Context(), id, req.Username, req.Email, req.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, ToResponse(user))
}

// Delete deletes user.
// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204
// @Router /users/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}
