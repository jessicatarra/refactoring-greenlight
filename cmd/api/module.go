package main

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type Module struct {
	Name     string
	Addr     int
	Handler  http.Handler
	Server   *http.Server
	stopChan chan struct{}
}

func NewModule(name string, addr int, handler http.Handler) *Module {
	return &Module{
		Name:     name,
		Addr:     addr,
		Handler:  handler,
		Server:   &http.Server{Addr: fmt.Sprintf(":%d", addr), Handler: handler},
		stopChan: make(chan struct{}),
	}
}

func (m *Module) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Module %s is running on %d\n", m.Name, m.Addr)

	go func() {
		err := m.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Module %s encountered an error: %v\n", m.Name, err)
		}
	}()

	// Wait for stop signal
	<-m.stopChan

	fmt.Printf("Module %s received stop signal. Shutting down...\n", m.Name)

	// Shutdown the server gracefully
	if err := m.Server.Shutdown(nil); err != nil {
		fmt.Printf("Module %s shutdown error: %v\n", m.Name, err)
	} else {
		fmt.Printf("Module %s stopped\n", m.Name)
	}
}

func (m *Module) Stop() {
	close(m.stopChan)
}
