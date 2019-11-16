package typmod_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typmod"
)

func TestIsProvider(t *testing.T) {
	testCases := []struct {
		obj        interface{}
		isProvider bool
	}{
		{dummyObj{}, true},
		{struct{}{}, false},
	}
	for i, tt := range testCases {
		require.Equal(t, tt.isProvider, typmod.IsProvider(tt.obj), i)
	}
}

func TestIsPreparer(t *testing.T) {
	testCases := []struct {
		obj        interface{}
		isPreparer bool
	}{
		{dummyObj{}, true},
		{struct{}{}, false},
	}
	for i, tt := range testCases {
		require.Equal(t, tt.isPreparer, typmod.IsPreparer(tt.obj), i)
	}
}

func TestIsDestroyer(t *testing.T) {
	testCases := []struct {
		obj         interface{}
		isDestroyer bool
	}{
		{dummyObj{}, true},
		{struct{}{}, false},
	}
	for i, tt := range testCases {
		require.Equal(t, tt.isDestroyer, typmod.IsDestroyer(tt.obj), i)
	}
}

func TestConfigurer(t *testing.T) {
	testCases := []struct {
		obj          interface{}
		isConfigurer bool
	}{
		{dummyObj{}, true},
		{struct{}{}, false},
	}
	for i, tt := range testCases {
		require.Equal(t, tt.isConfigurer, typmod.IsConfigurer(tt.obj), i)
	}
}

type dummyObj struct{}

func (dummyObj) Run() interface{}                    { return nil }
func (dummyObj) Prepare() []interface{}              { return nil }
func (dummyObj) Provide() []interface{}              { return nil }
func (dummyObj) Destroy() []interface{}              { return nil }
func (dummyObj) Configure() typmod.Configuration { return typmod.Configuration{} }