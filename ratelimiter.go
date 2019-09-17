// Package ratelimiter provides various rate-limiting mechanisms for echo-based APIs. 
// RateLimiters are provided as middlewares which can be applied over APIs defined in the router.
package ratelimiter

import (
    "github.com/go-redis/redis"
    "github.com/labstack/echo"
    "net/http"
)

// Response returned by the API when rate-limit exceeds
type RateLimitResponse struct {
    Message string
    APIID   string `json:"apiID"`
}

type RateLimiter interface {
    apply()                                     echo.MiddlewareFunc
    hasLimitExceeded(string)                    bool
}

// All rate-limiter implmentations need to derive from this base struct
type AbstractRateLimiter struct {
    windowSize                  int
    maxRequest                  int
    redisConnection             redis.Cmdable
    useIP                       bool
    RateLimiter
}

// Check if limit has exceeded or not, and return rate-limiting response accordingly. Specific rate-limiting 
// implmentations need to define the logic for `hasLimitExceeded`
func (rateLimiter AbstractRateLimiter) apply() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            key := c.Path()
            if rateLimiter.useIP {
                ip := c.RealIP()
                key += ":" + ip
            }
            if rateLimiter.hasLimitExceeded(key) {
                return c.JSON(http.StatusTooManyRequests, RateLimitResponse{Message: "Too many requests"})
            }
            return next(c)
        }
    }
}

// Fixed Window rate limiting implementation
type FixedWindowRateLimiter struct {
    AbstractRateLimiter
}

func (fwLimiter FixedWindowRateLimiter) hasLimitExceeded(key string) bool {
    key = "FWLimiter:" + key
    incrementScript := `
        local current
        current = redis.call("incr",KEYS[1])
        if tonumber(current) == 1 then
            redis.call("expire",KEYS[1],ARGV[1])
        end
        return current
    `
    cmd := fwLimiter.redisConnection.Eval(incrementScript, []string{key}, fwLimiter.windowSize)
    count, err := cmd.Int()
    if err != nil {
        return false
    }
    if count > fwLimiter.maxRequest {
        return true
    }
    return false
}

// Middleware for Fixed Window Rate limiter
func FixedWindowRateLimiterMiddleware(windowSize, maxRequest int, redisConnection redis.Cmdable, useIP bool) echo.MiddlewareFunc {
    fwLimiter := FixedWindowRateLimiter{AbstractRateLimiter{
        windowSize: windowSize, maxRequest: maxRequest, redisConnection: redisConnection, useIP: useIP}}
    fwLimiter.AbstractRateLimiter.RateLimiter = fwLimiter
    return fwLimiter.apply()
}
