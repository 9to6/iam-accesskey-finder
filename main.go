package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	expireTime int
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/health-check", func(c *gin.Context) {
		c.String(http.StatusOK, "healthy")
	})

	// Get user value
	r.GET("/expired-keys", func(c *gin.Context) {
		accessKeyInfo, err := GetExpiredAccessKeys(expireTime)
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
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	val, found := os.LookupEnv("ACCESS_KEY_EXPIRE_TIME")
	if !found {
		slog.Error("need 'ACCESS_KEY_EXPIRE_TIME' as environment variable")
		return
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		slog.Error("expire time have to be integer")
		return
	}
	expireTime = i
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
