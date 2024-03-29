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

	Config Config
}

var instance Instance

// NewInstance ...
func NewInstance(cfg Config) Instance {
	// Init redis connection
	redisConn := asynq.RedisClientOpt{
		Addr:     cfg.Redis.URL,
		Password: cfg.Redis.Password,
		DB:       0,
	}

	// Init instance
	instance.Server = initServer(redisConn, cfg)
	instance.Scheduler = initScheduler(redisConn)
	instance.Client = initClient(redisConn)
	instance.Config = cfg

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
	retryDelayFunc := cfg.RetryDelayFunc
	if retryDelayFunc == nil {
		// Default delay in 10s
		retryDelayFunc = func(n int, e error, t *asynq.Task) time.Duration {
			return 10 * time.Second
		}
	}

	// Init server
	server := asynq.NewServer(redisConn, asynq.Config{
		Concurrency: cfg.Concurrency,
		Queues: map[string]int{
			PriorityCritical: cfg.Priority.Critical,
			PriorityDefault:  cfg.Priority.Default,
			PriorityLow:      cfg.Priority.Low,
		},
		StrictPriority: cfg.Priority.StrictMode,

		RetryDelayFunc: retryDelayFunc,
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

// GetInstance ...
func GetInstance() Instance {
	return instance
}
