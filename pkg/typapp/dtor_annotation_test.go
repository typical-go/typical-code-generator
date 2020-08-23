package typapp_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typgo"
)

func TestDtorAnnotation_Annotate(t *testing.T) {
	os.MkdirAll("folder3/pkg3", 0777)
	defer os.RemoveAll("folder3")

	unpatch := execkit.Patch([]*execkit.RunExpectation{})
	defer unpatch(t)

	var out strings.Builder
	typapp.Stdout = &out
	defer func() { typapp.Stdout = os.Stdout }()

	dtorAnnot := &typapp.DtorAnnotation{}
	ctx := &typast.Context{
		Destination: "folder3/pkg3",
		Context: &typgo.Context{
			BuildSys: &typgo.BuildSys{
				Descriptor: &typgo.Descriptor{ProjectName: "some-project"},
			},
		},
		Summary: &typast.Summary{
			Annots: []*typast.Annot{
				{
					TagName: "@dtor",
					Decl: &typast.Decl{
						File: typast.File{Package: "pkg"},
						Type: &typast.FuncDecl{Name: "Clean"},
					},
				},
			},
		},
	}

	require.NoError(t, dtorAnnot.Annotate(ctx))

	b, _ := ioutil.ReadFile("folder3/pkg3/dtor_annotated.go")
	require.Equal(t, `package pkg3

/* Autogenerated by Typical-Go. DO NOT EDIT.

TagName:
	@dtor

Help:
	https://pkg.go.dev/github.com/typical-go/typical-go/pkg/typapp?tab=doc#DtorAnnotation
*/

import (
	"github.com/typical-go/typical-go/pkg/typapp"
)

func init() { 
	typapp.AppendDtor(
		&typapp.Destructor{Fn: pkg.Clean},
	)
}`, string(b))

	require.Equal(t, "Generate @dtor to folder3/pkg3/dtor_annotated.go\n", out.String())

}

func TestDtorAnnotation_Annotate_Predefined(t *testing.T) {
	os.MkdirAll("folder4/pkg4", 0777)
	defer os.RemoveAll("folder4")

	unpatch := execkit.Patch([]*execkit.RunExpectation{})
	defer unpatch(t)

	var out strings.Builder
	typapp.Stdout = &out
	defer func() { typapp.Stdout = os.Stdout }()

	dtorAnnot := &typapp.DtorAnnotation{
		Target:   "some-target",
		TagName:  "@some-tag",
		Template: "some-template",
	}
	ctx := &typast.Context{
		Destination: "folder4/pkg4",
		Context: &typgo.Context{
			BuildSys: &typgo.BuildSys{
				Descriptor: &typgo.Descriptor{ProjectName: "some-project"},
			},
		},
		Summary: &typast.Summary{
			Annots: []*typast.Annot{
				{
					TagName: "@some-tag",
					Decl: &typast.Decl{
						File: typast.File{Package: "pkg"},
						Type: &typast.FuncDecl{Name: "Clean"},
					},
				},
			},
		},
	}

	require.NoError(t, dtorAnnot.Annotate(ctx))

	b, _ := ioutil.ReadFile("folder4/pkg4/some-target")
	require.Equal(t, `some-template`, string(b))

	require.Equal(t, "Generate @dtor to folder4/pkg4/some-target\n", out.String())
}

func TestDtorAnnotation_Annotate_RemoveTargetWhenNoAnnotation(t *testing.T) {
	os.MkdirAll("folder4/pkg4", 0777)
	defer os.RemoveAll("folder4")

	var out strings.Builder
	typapp.Stdout = &out
	defer func() { typapp.Stdout = os.Stdout }()

	dtorAnnot := &typapp.DtorAnnotation{Target: "some-target"}
	ctx := &typast.Context{
		Destination: "folder4/pkg4",
		Context:     &typgo.Context{},
		Summary:     &typast.Summary{},
	}

	ioutil.WriteFile("folder4/pkg4/some-target", []byte("some-content"), 0777)
	require.NoError(t, dtorAnnot.Annotate(ctx))
	_, err := os.Stat("folder4/pkg4/some-target")
	require.True(t, os.IsNotExist(err))

	require.Equal(t, "", out.String())
}

func TestDtorAnnotation_IsDtor(t *testing.T) {
	testcases := []struct {
		TestName string
		*typapp.DtorAnnotation
		Annot    *typast.Annot
		Expected bool
	}{
		{
			TestName:       "private function",
			DtorAnnotation: &typapp.DtorAnnotation{TagName: "@dtor"},
			Annot: &typast.Annot{
				TagName: "@dtor",
				Decl: &typast.Decl{
					Type: &typast.FuncDecl{Name: "someFunction"},
				},
			},
			Expected: false,
		},
		{
			TestName:       "public function",
			DtorAnnotation: &typapp.DtorAnnotation{TagName: "@dtor"},
			Annot: &typast.Annot{
				TagName: "@dtor",
				Decl: &typast.Decl{
					Type: &typast.FuncDecl{Name: "SomeFunction"},
				},
			},
			Expected: true,
		},
		{
			TestName:       "not function",
			DtorAnnotation: &typapp.DtorAnnotation{TagName: "@dtor"},
			Annot: &typast.Annot{
				TagName: "@dtor",
				Decl: &typast.Decl{
					Type: &typast.InterfaceDecl{
						TypeDecl: typast.TypeDecl{Name: "SomeInterface"},
					},
				},
			},
			Expected: false,
		},
		{
			TestName:       "method function",
			DtorAnnotation: &typapp.DtorAnnotation{TagName: "@dtor"},
			Annot: &typast.Annot{
				TagName: "@dtor",
				Decl: &typast.Decl{
					Type: &typast.FuncDecl{Name: "SomeMethod", Recv: &typast.FieldList{}},
				},
			},
			Expected: false,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			require.Equal(t, tt.Expected, tt.IsDtor(tt.Annot))
		})
	}
}
