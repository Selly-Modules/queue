package queue

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
)

// RunTask ...
func RunTask(client *asynq.Client, typename string, data map[string]interface{}, retryTimes int) (*asynq.TaskInfo, error) {
	// Convert to []byte
	payload, err := json.Marshal(data)
	if err != nil {
		msg := fmt.Sprintf("task type: %s - error when create new task: %s", typename, err.Error())
		return nil, errors.New(msg)
	}

	// Create task and options
	task := asynq.NewTask(typename, payload)
	options := make([]asynq.Option, 0)

	// Retry times
	if retryTimes < 0 {
		retryTimes = 0
	}
	options = append(options, asynq.MaxRetry(retryTimes))

	// Enqueue task
	return client.Enqueue(task, options...)
}

// ScheduledTask create new task and run at specific time
// cronSpec follow cron expression
// https://www.freeformatter.com/cron-expression-generator-quartz.html
func ScheduledTask(scheduler *asynq.Scheduler, typename string, data map[string]interface{}, cronSpec string) (string, error) {
	// Convert to []byte
	payload, err := json.Marshal(data)
	if err != nil {
		msg := fmt.Sprintf("task type: %s - error when create new task: %s", typename, err.Error())
		return "", errors.New(msg)
	}

	// Create task and options
	task := asynq.NewTask(typename, payload)

	// TODO: Support options later
	// options := make([]asynq.Option, 0)

	return scheduler.Register(cronSpec, task)
}
