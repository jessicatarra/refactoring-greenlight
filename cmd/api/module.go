package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
)

type Module struct {
	Name     string
	Addr     int
	Handler  http.Handler
	Server   *http.Server
	stopChan chan struct{}
	wg       *sync.WaitGroup
}

func NewModule(name string, addr int, handler http.Handler, wg *sync.WaitGroup) *Module {
	return &Module{
		Name:     name,
		Addr:     addr,
		Handler:  handler,
		Server:   &http.Server{Addr: fmt.Sprintf(":%d", addr), Handler: handler},
		stopChan: make(chan struct{}),
		wg:       wg,
	}
}

func (m *Module) Run(ctx context.Context) error {
	defer m.wg.Done()

	fmt.Printf("Module %s is running on %d\n", m.Name, m.Addr)

	go func() {
		err := m.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Module %s encountered an error: %v\n", m.Name, err)
			os.Exit(1)
		}
	}()

	// Wait for stop signal
	<-m.stopChan

	fmt.Printf("Module %s received stop signal. Shutting down...\n", m.Name)

	// Shutdown the server gracefully
	<-ctx.Done()
	fmt.Printf("Module %s stopped\n", m.Name)
	err := m.Server.Shutdown(ctx)
	return err

}

func (m *Module) Stop() {
	close(m.stopChan)
}
