package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ModularMonolith struct {
	Modules []*Module
	wg      *sync.WaitGroup
	ctx     context.Context
}

func NewModularMonolith(ctx context.Context, wg *sync.WaitGroup) *ModularMonolith {
	return &ModularMonolith{
		Modules: []*Module{},
		wg:      wg,
		ctx:     ctx,
	}
}

func (mm *ModularMonolith) AddModule(name string, addr int, handler http.Handler) {
	module := NewModule(name, addr, handler, mm.wg)
	mm.Modules = append(mm.Modules, module)
}

func (mm *ModularMonolith) Run() error {
	for _, module := range mm.Modules {
		mm.wg.Add(1)
		go module.Run(mm.ctx)
	}

	// Wait for termination signal
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal or all modules to finish
	select {
	case <-signalCh:
		fmt.Println("Termination signal received. Stopping modules...")
		mm.StopAllModules()
	case <-mm.waitAll():
		fmt.Println("All modules finished. Exiting...")
	}

	fmt.Println("All modules stopped. Exiting...")
	mm.waitAll()
	return nil
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
