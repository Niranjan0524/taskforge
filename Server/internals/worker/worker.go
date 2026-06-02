package worker

import (
	"context"
	"log"

	storage "github.com/Niranjan0524/taskforge/server/internals/Storage"
)

type WorkerPool struct {
	store      storage.Storage
	workerSize int
}

func NewWorkerPool(store storage.Storage, workerSize int) *WorkerPool {
	return &WorkerPool{
		store:      store,
		workerSize: workerSize,
	}
}

func (p *WorkerPool) Start(ctx context.Context) error {

	for i := 1; i <= p.workerSize; i++ {
		go p.runWorker(ctx, i)
	}
	<-ctx.Done()
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

			if err := p.store.MarkTaskRunning(ctx, task.ID); err != nil {
				log.Println("error marking task running:", err)
				continue
			}

			if err := ExecuteTask(ctx, task); err != nil {
				log.Println("error executing task:", err)
				if markErr := p.store.MarkTaskFailed(ctx, task.ID); markErr != nil {
					log.Println("error marking task failed:", markErr)
				}
				continue
			}

			if err := p.store.MarkTaskCompleted(ctx, task.ID); err != nil {
				log.Println("error marking task completed:", err)
			}
		}
	}
}
