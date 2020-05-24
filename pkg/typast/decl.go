package typast

const (
	// Function type
	Function DeclType = iota

	// Interface type
	Interface

	// Struct type
	Struct

	// Generic type
	Generic
)

type (
	// Decl stand of declaration
	Decl struct {
		Name string   `json:"name"`
		Path string   `json:"path"`
		Pkg  string   `json:"pkg"`
		Type DeclType `json:"type"`
	}

	// DeclType is declaration type
	DeclType int
)

func (d DeclType) String() string {
	return [...]string{"Function", "Interface", "Struct", "Generic"}[d]
}
