package queue

import "github.com/hibiken/asynq"

// NewClient ...
func NewClient(cfg Config) *asynq.Client {
	// Init redis connection
	redisConn := asynq.RedisClientOpt{
		Addr:     cfg.Redis.URL,
		Password: cfg.Redis.Password,
		DB:       0,
	}

	// Init client
	if client := asynq.NewClient(redisConn); client == nil {
		panic("error when initializing queue CLIENT")
	} else {
		return client
	}
}
