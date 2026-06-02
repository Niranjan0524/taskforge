package worker

import (
	"context"
	"fmt"
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
			// later: pop task from Redis queue and execute it
			fmt.Println("executing Tasks")
		}
	}
}
