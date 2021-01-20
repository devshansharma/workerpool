package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const workerPoolSize = 4

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// create a consumer
	consumer := &Consumer{
		jobsChan:   make(chan int, workerPoolSize),
		ingestChan: make(chan int, 1),
	}

	producer := &Producer{callbackFunc: consumer.callbackFunc}

	go consumer.startConsumer(ctx)

	wg.Add(workerPoolSize)
	for i := 0; i < workerPoolSize; i++ {
		go consumer.workerFunc(wg, i)
	}

	producer.start()

	// block until receives event
	termChannel := make(chan os.Signal)
	signal.Notify(termChannel, syscall.SIGINT, syscall.SIGTERM)
	<-termChannel

	// cancel context and wait for all processes to go down
	fmt.Println("Shutdown Signal Received.")
	cancelFunc()
	wg.Wait()
	fmt.Println("All workers done, shutting down!")
}

// when we receive interrupt signal,
// cancelFunc will be called, instructing workers to go down.
// since it will take workers some time to process on going requests before going down
// waitgroup will block access.

// consumer is taking input from outside world with ingestChannel and
// also taking from context.Done()
// then if there is data in ingestChannel it will give it to jobsChannel to perform work
// otherwise close the jobsChannel if the context is Done.
