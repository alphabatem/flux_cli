package context

// Service interface defines what each service needs to expose to be included
// within the service context.
type Service interface {
	Id() string
	Configure(ctx *Context) error
	Start() error
	Shutdown()
}
