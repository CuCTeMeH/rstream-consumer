package task_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cuctemeh/rstream-consumer/internal/task"
)

func TestRunner_Run(t *testing.T) {
	runner := task.NewTaskRunner()

	// Define a task that will return an error
	task := func(ctx context.Context) error {
		return errors.New("task error")
	}

	// Run the task
	runner.Run(task)

	// Wait for the task to finish
	err := runner.Wait(context.Background())

	// Check that the error from the task was returned
	if err == nil || err.Error() != "task error" {
		t.Errorf("Wait() = %v, want 'task error'", err)
	}
}
