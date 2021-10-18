package queue

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// NewWorker ...
func NewWorker(cfg Config) *asynq.ServeMux {
	// Init redis connection
	redisConn := asynq.RedisClientOpt{
		Addr:     cfg.Redis.URL,
		Password: cfg.Redis.Password,
		DB:       0,
	}

	// Set default for concurrency
	if cfg.Concurrency == 0 {
		cfg.Concurrency = 100
	}

	// Set default for priority
	if cfg.Priority.Critical == 0 || cfg.Priority.Default == 0 || cfg.Priority.Low == 0 {
		cfg.Priority.Critical = 6
		cfg.Priority.Default = 3
		cfg.Priority.Low = 1
		cfg.Priority.StrictMode = false
	}

	// Init worker
	worker := asynq.NewServer(redisConn, asynq.Config{
		Concurrency: cfg.Concurrency,
		Queues: map[string]int{
			priorityCritical: cfg.Priority.Critical,
			priorityDefault:  cfg.Priority.Default,
			priorityLow:      cfg.Priority.Low,
		},
		StrictPriority: cfg.Priority.StrictMode,

		// TODO:
		// This is default option, retry after 10s, will add to config later
		RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
			return 10 * time.Second
		},
	})

	// Init mux server
	mux := asynq.NewServeMux()

	// Run server
	go func() {
		if err := worker.Run(mux); err != nil {
			msg := fmt.Sprintf("error when initializing queue WORKER: %s", err.Error())
			panic(msg)
		}
	}()

	// Return
	return mux
}
