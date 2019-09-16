func RateLimitMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            if hasLimitExceeded() {
                return c.JSON(http.StatusTooManyRequests, structs.Response{Success: false, Message: "OK"})
            }
            return next(c)
        }
    }
}

func hasLimitExceeded() bool {
    return true
}