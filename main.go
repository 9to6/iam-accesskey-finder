package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	expireTime time.Duration
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.GET("/health-check", func(c *gin.Context) {
		c.String(http.StatusOK, "healthy")
	})

	r.GET("/expired-keys", func(c *gin.Context) {
		accessKeyInfo, err := GetExpiredAccessKeys(c.Request.Context(), expireTime)
		if err != nil {
			slog.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal server error"})
			return
		}

		slog.Info("target access keys:", "count", len(accessKeyInfo), "keys", accessKeyInfo)
		if len(accessKeyInfo) == 0 {
			c.JSON(http.StatusNoContent, gin.H{"msg": "every key is available"})
		} else {
			c.JSON(http.StatusOK, gin.H{"keys": accessKeyInfo})
		}
	})

	return r
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	val, found := os.LookupEnv("ACCESS_KEY_EXPIRE_TIME")
	if !found {
		slog.Error("need 'ACCESS_KEY_EXPIRE_TIME' as environment variable")
		return
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		slog.Error("cannot parse duration, please input golang duration style")
		return
	}
	expireTime = d
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
