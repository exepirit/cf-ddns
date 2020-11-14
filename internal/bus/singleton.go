package bus

import "sync"

var mainBus *ApplicationBus
var busObjectLock sync.Mutex

func Get() Publisher {
	busObjectLock.Lock()
	defer busObjectLock.Unlock()
	if mainBus == nil {
		mainBus = New()
	}
	return mainBus
}
