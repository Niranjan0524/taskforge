package redisStore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	storage "github.com/Niranjan0524/taskforge/server/internals/Storage"
	"github.com/Niranjan0524/taskforge/server/internals/handlers/webSockets"
	"github.com/redis/go-redis/v9"
)

type redisStruct struct {
	Client *redis.Client
}

func NewRedisStore(client *redis.Client) *redisStruct {
	return &redisStruct{
		Client: client,
	}
}

func redisTaskKey(taskId string) string {
	if strings.HasPrefix(taskId, "task:") {
		return taskId
	}

	return "task:" + taskId
}

func (r *redisStruct) CreateTask(ctx context.Context, task storage.Task) error {

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	taskKey := redisTaskKey(task.ID)

	//this pipe helps to run multiple redis commands together
	pipe := r.Client.TxPipeline()

	pipe.Set(ctx, taskKey, taskJSON, 0)
	pipe.LPush(ctx, "queue:tasks", task.ID)
	pipe.SAdd(ctx, "tasks:pending", task.ID)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *redisStruct) GetTask(ctx context.Context, taskId string) (error, storage.Task) {

	if strings.TrimSpace(taskId) == "" {
		return errors.New("No userId found"), storage.Task{}
	}

	taskJSON, taskErr := r.Client.Get(ctx, redisTaskKey(taskId)).Result()
	if taskErr != nil {
		return taskErr, storage.Task{}
	}

	var task storage.Task

	err := json.Unmarshal([]byte(taskJSON), &task)

	if err != nil {
		return err, storage.Task{}
	}

	return nil, task
}

func (r *redisStruct) GetTaskStatus(ctx context.Context, taskID string) (string, error) {
	if strings.TrimSpace(taskID) == "" {
		return "", errors.New("no taskId found")
	}

	taskKey := redisTaskKey(taskID)
	taskJSON, err := r.Client.Get(ctx, taskKey).Result()
	if err != nil {
		return "", err
	}

	var task storage.Task
	if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
		return "", err
	}

	return task.Status, nil
}

func (r *redisStruct) GetAllTasks(ctx context.Context) ([]storage.Task, error) {
	var cursor uint64
	var tasks []storage.Task

	for {
		keys, nextCursor, err := r.Client.Scan(ctx, cursor, "task:*", 100).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			taskJSON, err := r.Client.Get(ctx, key).Result()
			if err != nil {
				return nil, err
			}

			var task storage.Task
			if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
				return nil, err
			}

			tasks = append(tasks, task)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return tasks, nil
}

func (r *redisStruct) DeleteTask(ctx context.Context, taskID string) error {

	if taskID == "" {
		return errors.New("invalid task key")
	}
	taskKey := redisTaskKey(taskID)
	pipe := r.Client.TxPipeline()

	pipe.Del(ctx, taskKey)
	pipe.SRem(ctx, "tasks:pending", taskID)
	pipe.LRem(ctx, "queue:tasks", 0, taskID)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *redisStruct) PopTask(ctx context.Context) (storage.Task, error) {

	fmt.Println("waiting for task from queue:tasks")
	result, err := r.Client.BRPop(ctx, 0, "queue:tasks").Result()

	if err != nil {
		fmt.Println("Error in popping task", err)
		return storage.Task{}, err
	}

	taskId := result[1]
	fmt.Println("popped task id:", taskId)
	taskErr, task := r.GetTask(ctx, taskId)

	fmt.Println("popped task: ", task)
	if taskErr != nil {
		fmt.Println("Error in popTask", taskErr)
		return storage.Task{}, taskErr
	}

	return task, nil
}

func (r *redisStruct) UpdateTaskStatus(ctx context.Context, taskId string, status string) error {

	if strings.TrimSpace(taskId) == "" || strings.TrimSpace(status) == "" {
		fmt.Println("Insufficient Details", taskId)
		return errors.New("Insufficient Details")
	}

	taskKey := redisTaskKey(taskId)
	rawTaskId := strings.TrimPrefix(taskId, "task:")

	taskJson, err := r.Client.Get(ctx, taskKey).Result()

	if err != nil {
		fmt.Println("Error in getting the task", err)
		return err
	}

	var task storage.Task
	if err := json.Unmarshal([]byte(taskJson), &task); err != nil {
		return err
	}
	task.Status = status

	updatedTask, err := json.Marshal(task)

	if err != nil {
		return err
	}

	pipe := r.Client.TxPipeline()

	pipe.Set(ctx, taskKey, updatedTask, 0)

	pipe.SRem(ctx, "tasks:pending", rawTaskId)
	pipe.SRem(ctx, "tasks:running", rawTaskId)
	pipe.SRem(ctx, "tasks:completed", rawTaskId)
	pipe.SRem(ctx, "tasks:failed", rawTaskId)

	pipe.SAdd(ctx, "tasks:"+status, rawTaskId)

	_, err = pipe.Exec(ctx)

	if err == nil {
		data, marshalErr := webSockets.MarshalTaskStatus(rawTaskId, status)
		if marshalErr != nil {
			return marshalErr
		}

		if publishErr := r.Client.Publish(ctx, webSockets.TaskStatusChannel, string(data)).Err(); publishErr != nil {
			fmt.Println("Error publishing task status", publishErr)
		}

		webSockets.BroadcastTaskStatus(rawTaskId, status)
	}
	return err
}

func (r *redisStruct) MarkTaskRunning(ctx context.Context, taskId string) error {

	if strings.TrimSpace(taskId) == "" {
		fmt.Println("Empty tasksId", taskId)
		return errors.New("No TaskId Found")
	}

	err := r.UpdateTaskStatus(ctx, taskId, "running")

	if err != nil {
		return err
	}
	return nil
}
func (r *redisStruct) MarkTaskFailed(ctx context.Context, taskId string) error {

	if strings.TrimSpace(taskId) == "" {
		fmt.Println("Empty tasksId", taskId)
		return errors.New("No TaskId Found")
	}

	err := r.UpdateTaskStatus(ctx, taskId, "failed")

	if err != nil {
		return err
	}
	return nil
}
func (r *redisStruct) MarkTaskCompleted(ctx context.Context, taskId string) error {

	if strings.TrimSpace(taskId) == "" {
		fmt.Println("Empty tasksId", taskId)
		return errors.New("No TaskId Found")
	}

	err := r.UpdateTaskStatus(ctx, taskId, "completed")

	if err != nil {
		return err
	}
	return nil
}
