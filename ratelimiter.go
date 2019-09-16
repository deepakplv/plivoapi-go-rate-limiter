package ratelimiter

import (
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "net/http"
)

type Response struct {
    Success bool
    Message string
    Body    interface{}
    APIID   string `json:"apiID"`
}

func RateLimitMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if hasLimitExceeded() {
                return c.JSON(http.StatusTooManyRequests, Response{Success: false, Message: "OK"})
            }
            return next(c)
        }
    }
}

func hasLimitExceeded() bool {
    return true
}
