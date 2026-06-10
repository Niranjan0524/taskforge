package worker

import (
	"context"
	"errors"
	"fmt"
	"time"

	storage "github.com/Niranjan0524/taskforge/server/internals/Storage"
)

func ExecuteTask(store storage.Storage, ctx context.Context, task storage.Task) error {
	switch task.Type {
	case "send_email":
		fmt.Println("sending email:", task.Payload)
		// time.Sleep(10 * time.Second)
		check, err := checkIfCancelled(store, ctx, task)
		if check == true {
			return err
		}

	case "generate_report":
		fmt.Println("generating report:", task.Payload)
		check, err := checkIfCancelled(store, ctx, task)
		if check == true {
			return err
		}

	case "resize_image":
		fmt.Println("Resizing Image", task.Payload)
		check, err := checkIfCancelled(store, ctx, task)
		if check == true {
			return err
		}
	default:
		fmt.Println("unknown task type:", task.Type)
		check, err := checkIfCancelled(store, ctx, task)
		if check == true {
			return err
		}
	}

	return nil
}

func checkIfCancelled(store storage.Storage, ctx context.Context, task storage.Task) (bool, error) {

	for i := 0; i < 10; i++ {
		check, err := store.IsCancelled(ctx, task)
		if err != nil {
			continue
		}
		if check == true {
			fmt.Println("Task is cancelled")
			return false, errors.New("Task is Cancelled")
		}
		time.Sleep(time.Second)
	}
	fmt.Println("taskExecuted")
	return false, nil
}
