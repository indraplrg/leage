package middleware

import (
	"net/http"
	"share-notes-app/internal/dtos"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)



func RateLimit(redisClient *redis.Client) gin.HandlerFunc {
		store := ratelimit.RedisStore(&ratelimit.RedisOptions{
			RedisClient: redisClient,
			Rate: 10 * time.Minute,
			Limit: 200,
		})
		mw := ratelimit.RateLimiter(store, &ratelimit.Options{
			ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, dtos.BaseResponse{
				Success: false,
				Message: "too many request, try again later",
			})
			},
			KeyFunc: func(c *gin.Context) string{
				logrus.Info("hit api dari", c.ClientIP())
				return c.ClientIP()
			},
		})
		return mw
}