package closer

import (
	"errors"
	"sync"
)

var globalCloser = New()

// Add ...
func Add(h handler) {
	globalCloser.Add(h)
}

// Close ...
func Close() error {
	return globalCloser.Close()
}

type handler func() error

// Closer ...
type Closer struct {
	mux      sync.Mutex
	handlers []handler
}

// New ...
func New() *Closer {
	return &Closer{}
}

// Add ...
func (c *Closer) Add(h handler) {
	c.mux.Lock()
	c.handlers = append(c.handlers, h)
	c.mux.Unlock()
}

// Close ...
func (c *Closer) Close() error {
	c.mux.Lock()
	defer c.mux.Unlock()

	var (
		allErrors []error
		mux       sync.Mutex
		wg        sync.WaitGroup
	)

	for _, h := range c.handlers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h(); err != nil {
				mux.Lock()
				allErrors = append(allErrors, err)
				mux.Unlock()
			}
		}()
	}

	wg.Wait()

	return errors.Join(allErrors...)
}
