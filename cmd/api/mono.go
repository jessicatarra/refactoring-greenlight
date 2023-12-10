package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Module interface {
	Start(wg *sync.WaitGroup)
	Shutdown(ctx context.Context, cancel func())
}

type ModularMonolith struct {
	Modules []Module
	wg      *sync.WaitGroup
}

func NewModularMonolith(wg *sync.WaitGroup) *ModularMonolith {
	return &ModularMonolith{
		Modules: []Module{},
		wg:      wg,
	}
}

func (mm *ModularMonolith) AddModule(module Module) {
	//module := NewModule(name, addr, handler, mm.wg)
	mm.Modules = append(mm.Modules, module)
}

func (mm *ModularMonolith) Run() error {
	for _, module := range mm.Modules {
		mm.wg.Add(1)
		go func(m Module) {
			defer mm.wg.Done()
			m.Start(mm.wg)
		}(module)
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, module := range mm.Modules {
		module.Shutdown(ctx, cancel)
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
