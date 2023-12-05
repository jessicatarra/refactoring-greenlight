package helpers

import (
	"fmt"
	"github.com/jessicatarra/greenlight/internal/jsonlog"
	"sync"
)

type Resource interface {
	Background(fn func())
}

type resource struct {
	wg     *sync.WaitGroup
	logger *jsonlog.Logger
}

func NewBackgroundTask(wg *sync.WaitGroup, logger *jsonlog.Logger) Resource {
	return &resource{
		wg:     wg,
		logger: logger,
	}
}

func (r *resource) Background(fn func()) {
	r.wg.Add(1)

	go func() {
		defer r.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				r.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()

		fn()
	}()
}
