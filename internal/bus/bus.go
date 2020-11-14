package bus

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type ApplicationBus struct {
	consumers []Consumer
	lock      sync.RWMutex
}

func New() *ApplicationBus {
	return &ApplicationBus{
		consumers: make([]Consumer, 0),
		lock:      sync.RWMutex{},
	}
}

func (a *ApplicationBus) Publish(event interface{}) {
	wg := sync.WaitGroup{}
	consume := func(c Consumer) {
		wg.Add(1)
		done := make(chan struct{})
		go func() {
			c.Consume(event)
			done <- struct{}{}
		}()
		t := time.NewTimer(time.Second)
		select {
		case <-t.C:
			log.Printf("failed consume %T event in consumer %T. possible deadlock", event, c)
		case <-done:
		}
		wg.Done()
	}

	a.lock.RLock()
	defer a.lock.RUnlock()
	for _, consumer := range a.consumers {
		go consume(consumer)
	}

	wg.Wait()
}

func (a *ApplicationBus) Subscribe(c Consumer) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.consumers = append(a.consumers, c)
	fmt.Printf("Consumer %T registered\n", c)
}
