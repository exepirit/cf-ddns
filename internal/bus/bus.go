package bus

import "sync"

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
		c.Consume(event)
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
}
