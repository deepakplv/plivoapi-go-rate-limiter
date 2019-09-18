# plivoapi-go-rate-limiter

This is a Golang library intented to be used with any Echo(server) based APIs to rate limit the APIs using various rate-limiting mechanisms via middlewares.

dep ensure -v -add  github.com/deepakplv/plivoapi-go-rate-limiter

And then use middlewares in router over the API(s) which need to be rate limited:
```go
import "github.com/deepakplv/plivoapi-go-rate-limiter"
var cache = redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", DB: 3})
router.GET("/some_url", controllers.SomeController, ratelimiter.FixedWindowRateLimiterMiddleware(60, 10, cache, false))
```
