package queue

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// Instance ...
type Instance struct {
	Client    *asynq.Client
	Server    *asynq.ServeMux
	Scheduler *asynq.Scheduler
}

// NewInstance ...
func NewInstance(cfg Config) Instance {
	// Init redis connection
	redisConn := asynq.RedisClientOpt{
		Addr:     cfg.Redis.URL,
		Password: cfg.Redis.Password,
		DB:       0,
	}

	// Init instance
	instance := Instance{}
	instance.Server = initServer(redisConn, cfg)
	instance.Scheduler = initScheduler(redisConn)
	instance.Client = initClient(redisConn)

	// Return instance
	return instance
}

func initServer(redisConn asynq.RedisClientOpt, cfg Config) *asynq.ServeMux {
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

	// Init server
	server := asynq.NewServer(redisConn, asynq.Config{
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
		if err := server.Run(mux); err != nil {
			msg := fmt.Sprintf("error when initializing queue SERVER: %s", err.Error())
			panic(msg)
		}
	}()

	return mux
}

func initScheduler(redisConn asynq.RedisClientOpt) *asynq.Scheduler {
	// Always run at HCM timezone
	l, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	// Init scheduler
	scheduler := asynq.NewScheduler(redisConn, &asynq.SchedulerOpts{
		Location: l,
	})

	// Run scheduler
	go func() {
		if err := scheduler.Run(); err != nil {
			msg := fmt.Sprintf("error when initializing queue SCHEDULER: %s", err.Error())
			panic(msg)
		}
	}()

	return scheduler
}

func initClient(redisConn asynq.RedisClientOpt) *asynq.Client {
	client := asynq.NewClient(redisConn)
	if client == nil {
		panic("error when initializing queue CLIENT")
	}
	return client
}
