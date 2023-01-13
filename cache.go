package sqlxcache

import (
	"sync"
	"time"
)

type data struct {
	value any
	expAt time.Time
}

type Cache struct {
	store      map[string]data
	mux        sync.RWMutex
	wg         sync.WaitGroup
	loopCancel chan struct{}
}

var defaultCleanUpInterval = 5 * time.Second

func NewCache(cleanupInterval *time.Duration) *Cache {
	c := Cache{
		store:      make(map[string]data),
		mux:        sync.RWMutex{},
		wg:         sync.WaitGroup{},
		loopCancel: make(chan struct{}),
	}

	if cleanupInterval == nil {
		cleanupInterval = &defaultCleanUpInterval
	}

	c.StartCleanUp(*cleanupInterval)
	return &c
}

func (c *Cache) Put(k string, v any, exp time.Time) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.store[k] = data{
		value: v,
		expAt: exp,
	}
}

func (c *Cache) Get(k string) any {
	c.mux.RLock()
	defer c.mux.RUnlock()

	v, ok := c.store[k]
	if !ok {
		return nil
	}
	return v

}

func (c *Cache) StopCleanUp() {
	close(c.loopCancel)
	c.wg.Done()
}

func (c *Cache) StartCleanUp(cleanupInterval time.Duration) {
	c.wg.Add(1)
	go func(ci time.Duration) {
		defer c.wg.Done()
		c.cleanupLoop(ci)
	}(cleanupInterval)
}

func (c *Cache) cleanupLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-c.loopCancel:
			break
		case <-t.C:
			for k, v := range c.store {
				now := time.Now()
				if v.expAt.Before(now) {
					delete(c.store, k)
				}
			}
		}
	}
}
