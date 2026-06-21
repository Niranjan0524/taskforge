package redisStore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	// pipe.LPush(ctx, "queue:tasks", task.ID)

	priorityScore := float64(task.Priority)*1e13 - float64(time.Now().UnixMilli())
	pipe.ZAdd(ctx, "queue:priority", redis.Z{Score: priorityScore, Member: task.ID})

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
	result, err := r.Client.ZPopMax(ctx, "queue:priority").Result()

	if err != nil {
		fmt.Println("Error in popping task", err)
		return storage.Task{}, err
	}

	if len(result) == 0 {
		fmt.Println("No more tasks available")
		time.Sleep(2 * time.Second)
		return storage.Task{}, nil
	}
	taskId := result[0].Member.(string)
	priorityScore := result[0].Score
	fmt.Println("popped task id:", taskId)
	fmt.Println("Popped task priority:", priorityScore)

	taskErr, task := r.GetTask(ctx, taskId)

	// fmt.Println("popped task: ", task)
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

	pipe.ZRem(ctx, "tasks:processing", rawTaskId)
	pipe.SRem(ctx, "tasks:completed", rawTaskId)
	pipe.SRem(ctx, "tasks:failed", rawTaskId)
	pipe.ZRem(ctx, "queue:priority", rawTaskId)

	if status == "completed" || status == "pending" || status == "cancelled" || status == "failed" {
		pipe.SAdd(ctx, "tasks:"+status, rawTaskId)
	} else if status == "dead" {
		pipe.ZAdd(ctx, "queue:dead",
			redis.Z{
				Score:  float64(time.Now().Unix()),
				Member: taskId,
			},
		)
	} else if status == "running" {
		pipe.ZAdd(ctx, "tasks:processing",
			redis.Z{
				Score:  float64(time.Now().Unix()),
				Member: taskId,
			},
		)
	}

	_, err = pipe.Exec(ctx)

	if err == nil {
		data, marshalErr := webSockets.MarshalTaskStatus(rawTaskId, task)
		if marshalErr != nil {
			return marshalErr
		}

		if publishErr := r.Client.Publish(ctx, webSockets.TaskStatusChannel, string(data)).Err(); publishErr != nil {
			fmt.Println("Error publishing task status", publishErr)
		}

		webSockets.BroadcastTaskStatus(rawTaskId, task)
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
		fmt.Println("Failed to store in processing queue")
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

func (r *redisStruct) GetStaleTasks(ctx context.Context) ([]string, error) {

	fiveMinutesAgo := time.Now().Add(-5 * time.Minute).Unix()

	tasks, err := r.Client.ZRangeArgs(
		ctx,
		redis.ZRangeArgs{
			Key:     "tasks:running",
			Start:   "-inf",
			Stop:    strconv.FormatInt(fiveMinutesAgo, 10),
			ByScore: true,
		},
	).Result()

	if err != nil {
		fmt.Println("Error in getting tasks from processing queue", err)
		return []string{}, err
	}

	return tasks, nil
}
func (r *redisStruct) MoveTaskToDeadQueue(ctx context.Context, taskId string) error {

	err := r.UpdateTaskStatus(ctx, taskId, "dead")

	if err != nil {
		fmt.Println("Failed to move task to dead queue", err)
		return err
	}

	return nil
}
func (r *redisStruct) CheckAndRetryTask(ctx context.Context, taskId string) (bool, error) {
	err, task := r.GetTask(ctx, taskId)

	if task.RetryCount >= task.MaxRetries {
		movErr := r.MoveTaskToDeadQueue(ctx, task.ID)
		if movErr != nil {
			fmt.Println("Failed to move task to dead queue", taskId)
			return false, nil
		}
		return false, nil
	}

	task.RetryCount++

	updatedJSON, err := json.Marshal(task)
	if err != nil {
		return false, err
	}

	err = r.Client.Set(
		ctx,
		redisTaskKey(task.ID),
		updatedJSON,
		0,
	).Err()

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *redisStruct) Requeue(ctx context.Context, taskId string) error {

	err, task := r.GetTask(ctx, taskId)

	if err != nil {
		fmt.Println("Error in getting tasks from processing queue", err)
		return err
	}
	statusErr := r.UpdateTaskStatus(ctx, task.ID, "pending")
	if statusErr != nil {
		fmt.Println("Task status update failed", statusErr)
		return statusErr
	}
	priorityScore := float64(task.Priority)*1e13 - float64(time.Now().UnixMilli())
	_, reqErr := r.Client.ZAdd(ctx, "queue:priority", redis.Z{Score: priorityScore, Member: task.ID}).Result()

	if reqErr != nil {
		fmt.Println("Requeue Error", reqErr)
		return reqErr
	}

	return nil
}

func (r *redisStruct) IsCancelled(ctx context.Context, task storage.Task) (bool, error) {

	err, task := r.GetTask(ctx, task.ID)

	if err != nil {
		fmt.Println("Error in checking status", err)
		return false, err
	}

	if task.Status == "cancelled" {
		return true, nil
	}

	return false, nil

}

func (r *redisStruct) CancelTask(ctx context.Context, taskId string) error {
	err, task := r.GetTask(ctx, taskId)

	if err != nil {
		fmt.Println("Error in cancelling task", err)
		return err
	}

	switch task.Status {

	case "completed":
		return errors.New("task already completed")

	case "failed":
		return errors.New("task already failed")

	default:
		err := r.UpdateTaskStatus(ctx, taskId, "cancelled")
		if err != nil {
			fmt.Println("error", err)
			return err
		}
	}

	return nil
}
