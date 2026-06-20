package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	storage "github.com/Niranjan0524/taskforge/server/internals/Storage"
	"github.com/Niranjan0524/taskforge/server/internals/handlers/webSockets"
	"github.com/redis/go-redis/v9"
)

type WorkerPool struct {
	store      storage.Storage
	workerSize int
	redis      *redis.Client
}

func NewWorkerPool(store storage.Storage, workerSize int, redisClient *redis.Client) *WorkerPool {
	return &WorkerPool{
		store:      store,
		workerSize: workerSize,
		redis:      redisClient,
	}
}

func (p *WorkerPool) Start(ctx context.Context) error {

	go p.runRecoveryWorker(ctx)
	if err := webSockets.PublishWorkerStatus(ctx, p.redis, "started"); err != nil {
		log.Println("error publishing worker started status:", err)
	}

	for i := 1; i <= p.workerSize; i++ {
		go p.runWorker(ctx, i)
	}

	<-ctx.Done()

	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := webSockets.PublishWorkerStatus(stopCtx, p.redis, "stopped"); err != nil {
		log.Println("error publishing worker stopped status:", err)
	}

	return ctx.Err()
}

func (p *WorkerPool) runWorker(ctx context.Context, workerID int) {
	log.Println("worker started:", workerID)

	for {
		select {
		case <-ctx.Done():
			log.Println("worker stopped:", workerID)
			return
		default:

			task, err := p.store.PopTask(ctx)

			if err != nil {
				if ctx.Err() != nil {
					return
				}

				log.Println("error popping task:", err)
				continue
			}

			if task.ID == "" {
				fmt.Println("No more tasks available")
				continue
			}

			if err := p.store.MarkTaskRunning(ctx, task.ID); err != nil {
				log.Println("error marking task running:", err)
				continue
			}

			if err := ExecuteTask(p.store, ctx, task); err != nil {
				log.Println("error executing task:", err)

				// Check if task was cancelled
				if err.Error() == CancelledTaskError {
					log.Println("task was cancelled:", task.ID)
					// Task is already marked as cancelled by the cancel API
					continue
				}

				// Regular execution error - mark as failed
				markErr := p.store.MarkTaskFailed(ctx, task.ID)

				if markErr != nil {
					log.Println("error marking task failed:", markErr)
				}

				toRequeue, reqErr := p.store.CheckAndRetryTask(ctx, task.ID)
				if toRequeue == false {
					fmt.Println("Error in updating the Retry", reqErr)
					continue
				}
				// move back to queue
				err := p.store.Requeue(ctx, task.ID)
				if err != nil {
					fmt.Println(err)
				}
				continue
			}

			if err := p.store.MarkTaskCompleted(ctx, task.ID); err != nil {
				log.Println("error marking task completed:", err)
			}
		}
	}
}

func (p *WorkerPool) runRecoveryWorker(ctx context.Context) {
	fmt.Println("Running recovery Routine")
	ticker := time.NewTicker(time.Minute)

	for range ticker.C {
		staleTasks, err := p.store.GetStaleTasks(ctx)

		if err != nil {
			fmt.Println("Error", err)
			return
		}

		for _, taskID := range staleTasks {
			toRequeue, reqErr := p.store.CheckAndRetryTask(ctx, taskID)
			if toRequeue == false {
				fmt.Println("Error in updating the Retry", reqErr)
				continue
			}
			// move back to queue
			err := p.store.Requeue(ctx, taskID)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
