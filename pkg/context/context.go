package context

import (
	context2 "context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Context is a small service wrapper that handles the startup/shutdown of services.
type Context struct {
	startOrder map[int]string
	serviceMap map[string]Service
}

// NewCtx creates a new context containing the given services.
func NewCtx(svcs ...Service) (*Context, error) {
	ctx := Context{
		startOrder: make(map[int]string, len(svcs)),
		serviceMap: make(map[string]Service, len(svcs)),
	}

	for _, s := range svcs {
		if err := ctx.Register(s); err != nil {
			return nil, err
		}
	}

	return &ctx, nil
}

// Register a new service into the context and preserve the order passed.
func (ctx *Context) Register(service Service) error {
	if _, ok := ctx.serviceMap[service.Id()]; ok {
		return fmt.Errorf("service %s already registered", service.Id())
	}

	currLen := len(ctx.serviceMap)

	ctx.startOrder[currLen] = service.Id()
	ctx.serviceMap[service.Id()] = service

	return nil
}

// Service returns the pointer to the given service.
func (ctx *Context) Service(id string) Service {
	return ctx.serviceMap[id]
}

// Run starts the context. Each service is configured first, then started.
func (ctx *Context) Run() error {
	_, cancel := context2.WithCancel(context2.Background())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v. Shutting down...", sig)

		for i := 0; i < len(ctx.startOrder); i++ {
			svcId := ctx.startOrder[i]
			log.Printf("Shutting down %s...", svcId)
			ctx.serviceMap[svcId].Shutdown()
		}
		cancel()
	}()

	for i := 0; i < len(ctx.startOrder); i++ {
		svcId := ctx.startOrder[i]

		if err := ctx.Configure(ctx.serviceMap[svcId]); err != nil {
			log.Fatalf("Context Configure Error: %s - %s", svcId, err)
			return err
		}
	}

	for i := 0; i < len(ctx.startOrder); i++ {
		svcId := ctx.startOrder[i]

		if err := ctx.Start(ctx.serviceMap[svcId]); err != nil {
			log.Fatalf("Context Start Error: %s - %s", svcId, err)
			return err
		}
	}

	return nil
}

// Configure the given service.
func (ctx *Context) Configure(svc Service) error {
	log.Printf("Context Configure: %s", svc.Id())

	if err := svc.Configure(ctx); err != nil {
		return err
	}

	return nil
}

// Start the given service.
func (ctx *Context) Start(svc Service) error {
	log.Printf("Context Start: %s", svc.Id())

	if err := svc.Start(); err != nil {
		return err
	}

	return nil
}

func (ctx *Context) Services() []string {
	var keys []string
	for k := range ctx.serviceMap {
		keys = append(keys, k)
	}

	return keys
}
