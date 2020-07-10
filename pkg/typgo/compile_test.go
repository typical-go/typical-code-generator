package typgo_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typgo"
)

func TestCompiler(t *testing.T) {
	testcases := []struct {
		testName string
		typgo.Compiler
		context     *typgo.Context
		expectedErr string
	}{
		{
			Compiler:    typgo.NewCompiler(func(*typgo.Context) error { return errors.New("some-error") }),
			expectedErr: "some-error",
		},
		{
			Compiler: typgo.NewCompiler(func(*typgo.Context) error { return nil }),
		},
		{
			Compiler: typgo.Compilers{
				typgo.NewCompiler(func(*typgo.Context) error { return nil }),
				typgo.NewCompiler(func(*typgo.Context) error { return errors.New("some-error") }),
			},
			expectedErr: "some-error",
		},
		{
			Compiler: typgo.Compilers{
				typgo.NewCompiler(func(*typgo.Context) error { return errors.New("some-error") }),
				typgo.NewCompiler(func(*typgo.Context) error { return nil }),
			},
			expectedErr: "some-error",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			err := tt.Compile(tt.context)
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
