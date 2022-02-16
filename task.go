package queue

import (
	"github.com/hibiken/asynq"
	"time"
)

const (
	NoTimeout time.Duration = 0
)

// RunTask ...
func (i Instance) RunTask(typename string, payload []byte, priority string, retryTimes int) (*asynq.TaskInfo, error) {
	// Create task and options
	task := asynq.NewTask(typename, payload)
	options := make([]asynq.Option, 0)

	// Priority
	if priority != PriorityCritical && priority != PriorityDefault && priority != PriorityLow {
		priority = PriorityDefault
	}
	options = append(options, asynq.Queue(priority))

	// Retry times
	if retryTimes < 0 {
		retryTimes = 0
	}
	options = append(options, asynq.MaxRetry(retryTimes))

	// Task timeout
	if i.Config.TaskTimeout != 0 {
		options = append(options, asynq.Timeout(i.Config.TaskTimeout))
	}

	// Enqueue task
	return i.Client.Enqueue(task, options...)
}

// ScheduledTask create new task and run at specific time
// cronSpec follow cron expression
// https://www.freeformatter.com/cron-expression-generator-quartz.html
func (i Instance) ScheduledTask(typename string, payload []byte, cronSpec string) (string, error) {
	// Create task and options
	task := asynq.NewTask(typename, payload)

	// TODO: Support options later
	// options := make([]asynq.Option, 0)

	return i.Scheduler.Register(cronSpec, task)
}
