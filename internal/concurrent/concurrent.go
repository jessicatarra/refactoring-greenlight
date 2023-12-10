package concurrent

import (
	"context"
	"fmt"
	"github.com/jessicatarra/greenlight/internal/errors"
	"sync"
)

type Resource interface {
	BackgroundTask(fn func() error)
}

type resource struct {
	wg  *sync.WaitGroup
	ctx context.Context
}

func NewBackgroundTask(ctx context.Context, wg *sync.WaitGroup) Resource {
	return &resource{ctx: ctx, wg: wg}
}

func (r *resource) BackgroundTask(fn func() error) {
	r.wg.Add(1)

	go func() {
		defer r.wg.Done()
		defer func() {
			err := recover()

			if err != nil {
				errors.ReportError(fmt.Errorf("%s", err))
				return
			}
		}()

		select {
		case <-r.ctx.Done():
			return
		default:
			err := fn()
			if err != nil {
				errors.ReportError(fmt.Errorf("%s", err))
				return
			}
		}
	}()
}
