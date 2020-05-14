package typgo

const (
	PrecondPhase Phase = iota
	TestPhase
	CompilePhase
	RunPhase
	ReleasePhase
	PublishPhase
	CleanPhase
)

var _ Build = (Builds)(nil)

type (
	// Phase of build process
	Phase int

	// Build responsible to execute build process
	Build interface {
		Execute(*Context, Phase) (bool, error)
	}

	// Builds is array of build
	Builds []Build
)

func (d Phase) String() string {
	return [...]string{"Function", "Interface", "Struct", "Generic"}[d]
}

// Execute build
func (b Builds) Execute(ctx *Context, phase Phase) (bool, error) {
	var ok bool
	for _, build := range b {
		ok1, err := build.Execute(ctx, phase)
		if err != nil {
			return ok1, err
		}
		ok = ok || ok1
	}
	return ok, nil
}