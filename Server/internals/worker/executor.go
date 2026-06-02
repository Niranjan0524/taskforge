package worker

import (
	"context"
	"fmt"

	storage "github.com/Niranjan0524/taskforge/server/internals/Storage"
)

func ExecuteTask(ctx context.Context, task storage.Task) error {
	switch task.Type {
	case "send_email":
		fmt.Println("sending email:", task.Payload)
	case "generate_report":
		fmt.Println("generating report:", task.Payload)
	default:
		fmt.Println("unknown task type:", task.Type)
	}

	return nil
}
