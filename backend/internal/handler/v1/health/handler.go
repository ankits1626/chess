// Package health provides health check endpoints.
package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents health check response.
type Response struct {
	Status  string `json:"status" example:"ok"`
	Service string `json:"service" example:"chess-coach-api"`
}

// Check handles health check endpoint.
// @Summary Health check
// @Description Check if service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /health [get]
func Check(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Service: "chess-coach-api",
	})
}
