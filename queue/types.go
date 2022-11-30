package queue

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
)

type TaskErrorHandler struct {
}

func (t *TaskErrorHandler) HandleError(ctx context.Context, task *asynq.Task, err error) {
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	if retried >= maxRetry {
		//重试次数用完了 archived
		fmt.Println(task.Type())
		fmt.Println(err)
	}
}
