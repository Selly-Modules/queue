package queue

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// NewScheduler ...
func NewScheduler(cfg Config) *asynq.Scheduler {
	// Init redis connection
	redisConn := asynq.RedisClientOpt{
		Addr:     cfg.Redis.URL,
		Password: cfg.Redis.Password,
		DB:       0,
	}

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

	// Return
	return scheduler
}
