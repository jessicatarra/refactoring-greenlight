package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ModularMonolith struct {
	Modules []*Module
	wg      sync.WaitGroup
}

func NewModularMonolith() *ModularMonolith {
	return &ModularMonolith{
		Modules: []*Module{},
		wg:      sync.WaitGroup{},
	}
}

func (mm *ModularMonolith) AddModule(name string, addr int, handler http.Handler) {
	module := NewModule(name, addr, handler)
	mm.Modules = append(mm.Modules, module)
}

func (mm *ModularMonolith) Run() {
	for _, module := range mm.Modules {
		mm.wg.Add(1)
		go module.Run(&mm.wg)
	}

	// Wait for termination signal
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Wait for termination signal or all modules to finish
	select {
	case <-signalCh:
		fmt.Println("Termination signal received. Stopping modules...")
		mm.StopAllModules()
	case <-mm.waitAll():
		fmt.Println("All modules finished. Exiting...")
	}
}

func (mm *ModularMonolith) StopAllModules() {
	for _, module := range mm.Modules {
		module.Stop()
	}
}

func (mm *ModularMonolith) waitAll() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		mm.wg.Wait()
		close(ch)
	}()
	return ch
}
