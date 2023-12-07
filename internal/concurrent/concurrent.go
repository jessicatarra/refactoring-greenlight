package concurrent

import (
	"fmt"
	"github.com/jessicatarra/greenlight/internal/errors"
	"sync"
)

type Resource interface {
	BackgroundTask(fn func() error)
}

type resource struct {
}

func NewBackgroundTask() Resource {
	return &resource{}
}

func (r *resource) BackgroundTask(fn func() error) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer func() {
			err := recover()

			if err != nil {
				errors.ReportError(fmt.Errorf("%s", err))
				return
			}
		}()

		err := fn()
		if err != nil {
			errors.ReportError(fmt.Errorf("%s", err))
			return
		}
	}()
}
