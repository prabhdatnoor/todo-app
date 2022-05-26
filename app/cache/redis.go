package cache

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"os"
	"strconv"
)

var (
	Redis *redis.Storage
	Store *session.Store
)

func ConnectStore() (*session.Store, error) {
	port, err := strconv.ParseInt(os.Getenv("REDIS_PORT"), 10, 32)

	host := os.Getenv("REDIS_HOST")

	if err != nil {
		Redis = redis.New()
		fmt.Print("ERROR: no redis port provided, reverting to default 6379")
	} else {
		Redis = redis.New(redis.Config{
			Host: host,
			Port: int(port)})
	}

	return session.New(session.Config{
		Storage: Redis,
	}), nil
}
