package queue

// Config ...
type Config struct {
	// For message queue
	Redis struct {
		URL      string
		Password string
	}
	// Priority to process task, eg: Critical 6, Default 3, Low 1
	// Using for server only
	// https://github.com/hibiken/asynq/wiki/Queue-Priority
	Concurrency int
	Priority    struct {
		Critical   int
		Default    int
		Low        int
		StrictMode bool
	}
}
