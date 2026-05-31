package redisStore

import (
	"context"
	"encoding/json"

	storage "github.com/Niranjan0524/taskforge/server/internals/Storage"
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

func (r *redisStruct) CreateTask(ctx context.Context, task storage.Task) error {

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	taskKey := "task:" + task.ID

	//this pipe helps to run multiple redis commands together
	pipe := r.Client.TxPipeline()

	pipe.Set(ctx, taskKey, taskJSON, 0)
	pipe.LPush(ctx, "queue:tasks", task.ID)
	pipe.SAdd(ctx, "tasks:pending", task.ID)

	_, err = pipe.Exec(ctx)
	return err
}
