package limiter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-rate-limiter/config"
	"github.com/go-rate-limiter/storage"
)

type RateLimiter struct {
	storage storage.Storage
	config  *config.Config
}

func NewRateLimiter(storage storage.Storage, config *config.Config) *RateLimiter {
	return &RateLimiter{
		storage: storage,
		config:  config,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, identifier string, isToken bool) (bool, error) {
	var (
		requests      int
		window        time.Duration
		blockDuration time.Duration
	)

	if isToken {
		requests = rl.config.RateLimitTokenRequests
		window = rl.config.RateLimitTokenWindow
		blockDuration = rl.config.RateLimitTokenBlockDuration
	} else {
		requests = rl.config.RateLimitIPRequests
		window = rl.config.RateLimitIPWindow
		blockDuration = rl.config.RateLimitIPBlockDuration
	}

	blockedKey := fmt.Sprintf("blocked:%s", identifier)
	blocked, err := rl.storage.Exists(ctx, blockedKey)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	counterKey := fmt.Sprintf("counter:%s", identifier)
	count, err := rl.storage.Increment(ctx, counterKey)
	if err != nil {
		return false, err
	}

	if count == 1 {
		if err := rl.storage.Set(ctx, counterKey, 1, window); err != nil {
			return false, err
		}
	}

	if count > int64(requests) {
		if err := rl.storage.Set(ctx, blockedKey, 1, blockDuration); err != nil {
			return false, err
		}
		if err := rl.storage.Delete(ctx, counterKey); err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ip := c.ClientIP()

		token := c.GetHeader("API_KEY")
		if token != "" {
			allowed, err := rl.Allow(ctx, token, true)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"message": "you have reached the maximum number of requests or actions allowed within a certain time frame",
				})
				return
			}
			c.Next()
			return
		}

		allowed, err := rl.Allow(ctx, ip, false)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "you have reached the maximum number of requests or actions allowed within a certain time frame",
			})
			return
		}

		c.Next()
	}
}
