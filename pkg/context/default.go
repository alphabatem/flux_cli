package context

// DefaultService should be extended for each service. Handles the internal context routing.
type DefaultService struct {
	ctx *Context
}

// Configure is the base that will be called for EVERY service extending DefaultService.
func (ds *DefaultService) Configure(ctx *Context) error {
	ds.ctx = ctx
	return nil
}

// Start the service.
func (ds *DefaultService) Start() error {
	return nil
}

// Shutdown performs any shutdown procedures.
func (ds *DefaultService) Shutdown() {}

// Service returns the inner context service by ID.
func (ds *DefaultService) Service(id string) Service {
	return ds.ctx.Service(id)
}

// Services returns the inner context list of service keys.
func (ds *DefaultService) Services() []string {
	return ds.ctx.Services()
}
