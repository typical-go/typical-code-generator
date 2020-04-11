package typapp

import "go.uber.org/dig"

var (
	global []*Constructor
)

var (
	_ Provider = (*Constructor)(nil)
)

// Constructor details
type Constructor struct {
	fn interface{}
}

// NewConstructor return new instance of constructor
func NewConstructor(fn interface{}) *Constructor {
	return &Constructor{
		fn: fn,
	}
}

// Provide the constructor to dig container
func (c *Constructor) Provide(di *dig.Container) (err error) {
	return di.Provide(c.fn)
}

// Fn is function of constructor
func (c *Constructor) Fn() interface{} {
	return c.fn
}

// Constructors is list of constructor
func (c *Constructor) Constructors() []*Constructor {
	return []*Constructor{c}
}

// AppendConstructor to append constructor
func AppendConstructor(cons ...*Constructor) {
	global = append(global, cons...)
}
