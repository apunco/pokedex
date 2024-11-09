package pokecache

import (
	"sync"
	"time"
)

type cache struct {
	cache  map[string]cacheEntry
	mu     *sync.Mutex
	exitC  chan (bool)
	closed bool
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *cache {
	cache := cache{
		cache:  make(map[string]cacheEntry),
		mu:     &sync.Mutex{},
		exitC:  make(chan bool, 1),
		closed: false,
	}

	go cache.reapLoop(interval)
	return &cache
}

func (c *cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.cache[key]
	return elem.val, ok
}

func (c *cache) remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
}

func (c *cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.exitC:
			c.CleanUp(ticker)
			return
		case <-ticker.C:
			c.mu.Lock()
			for key := range c.cache {

				if c.cache[key].createdAt.Add(interval).Before(time.Now()) {
					delete(c.cache, key)
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *cache) CleanUp(t *time.Ticker) {
	t.Stop()
	c.mu.Lock()
	c.cache = make(map[string]cacheEntry)
	c.mu.Unlock()
}

func (c *cache) Close() {
	c.mu.Lock()
	if !c.closed {
		c.closed = true
		c.mu.Unlock()
		c.exitC <- true
		return
	}
	c.mu.Unlock()
}
