package middleware

import (
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/cors"
)

// AllowedOrigins returns the list of allowed origins for CORS and WebSocket
func AllowedOrigins() []string {
    return []string{
        "http://localhost:5173", // Vite dev server
        "http://localhost:3000", // Alternative dev port
    }
}

func CORS() fiber.Handler {
    return cors.New(cors.Config{
        AllowOrigins:     AllowedOrigins(),
        AllowMethods: []string{
            fiber.MethodGet,
            fiber.MethodPost,
            fiber.MethodPut,
            fiber.MethodDelete,
            fiber.MethodOptions,
        },
        AllowHeaders: []string{
            "Origin",
            "Content-Type",
            "Accept",
            "Authorization",
        },
        AllowCredentials: true,
        MaxAge:           300, // 5 minutes
    })
}
