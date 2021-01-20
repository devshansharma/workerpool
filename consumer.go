package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Consumer will consume jobs
type Consumer struct {
	jobsChan   chan int
	ingestChan chan int
}

func (c *Consumer) callbackFunc(event int) {
	c.ingestChan <- event
}

// workerFunc will perform working
func (c *Consumer) workerFunc(wg *sync.WaitGroup, i int) {
	defer wg.Done()

	fmt.Printf("Worker %d starting...\n", i)
	for eventIndex := range c.jobsChan {
		fmt.Printf("Worker %d started job %d\n", i, eventIndex)
		time.Sleep(time.Millisecond * time.Duration(1000+rand.Intn(2000)))
		fmt.Printf("Worker %d finished processing job %d\n", i, eventIndex)
	}
	fmt.Printf("Worker %d interrupted...\n", i)
}

func (c *Consumer) startConsumer(ctx context.Context) {
	for {
		select {
		case job := <-c.ingestChan:
			c.jobsChan <- job
		case <-ctx.Done():
			fmt.Println("Consumer received cancellation signal. closing jobschan!")
			close(c.jobsChan)
			fmt.Println("Consumer closed jobschan")
			return
		}
	}
}
