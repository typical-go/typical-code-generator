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

func TestCtorAnnotation_Annotate(t *testing.T) {
	os.MkdirAll("folder1/dest1", 0777)
	defer os.RemoveAll("folder1")

	unpatch := execkit.Patch([]*execkit.RunExpectation{})
	defer unpatch(t)

	var out strings.Builder
	typapp.Stdout = &out
	defer func() { typapp.Stdout = os.Stdout }()

	ctorAnnot := &typapp.CtorAnnotation{}
	ctx := &typast.Context{
		Destination: "folder1/dest1",
		Context: &typgo.Context{
			BuildSys: &typgo.BuildSys{
				Descriptor: &typgo.Descriptor{ProjectName: "some-project"},
			},
		},
		Summary: &typast.Summary{
			Annots: []*typast.Annot{
				{
					TagName: "@ctor",
					Decl: &typast.Decl{
						Type: &typast.FuncDecl{Name: "NewObject"},
						File: typast.File{Package: "pkg"},
					},
				},
				{
					TagName:  "@ctor",
					TagParam: `name:"obj2"`,
					Decl: &typast.Decl{
						File: typast.File{Package: "pkg2"},
						Type: &typast.FuncDecl{Name: "NewObject2"},
					},
				},
			},
		},
	}

	require.NoError(t, ctorAnnot.Annotate(ctx))

	b, _ := ioutil.ReadFile("folder1/dest1/ctor_annotated.go")
	require.Equal(t, `package dest1

/* Autogenerated by Typical-Go. DO NOT EDIT.

TagName:
	@ctor

Help:
	https://pkg.go.dev/github.com/typical-go/typical-go/pkg/typapp?tab=doc#CtorAnnotation
*/

import (
	"github.com/typical-go/typical-go/pkg/typapp"
)

func init() { 
	typapp.AppendCtor(
		&typapp.Constructor{Name: "", Fn: pkg.NewObject},
		&typapp.Constructor{Name: "obj2", Fn: pkg2.NewObject2},
	)
}`, string(b))

	require.Equal(t, "Generate @ctor to folder1/dest1/ctor_annotated.go\n", out.String())

}

func TestCtorAnnotation_Annotate_Predefined(t *testing.T) {
	os.MkdirAll("folder2/dest2", 0777)
	defer os.RemoveAll("folder2")

	unpatch := execkit.Patch([]*execkit.RunExpectation{})
	defer unpatch(t)

	var out strings.Builder
	typapp.Stdout = &out
	defer func() { typapp.Stdout = os.Stdout }()

	ctorAnnot := &typapp.CtorAnnotation{
		TagName:  "@some-tag",
		Target:   "some-target",
		Template: "some-template",
	}
	ctx := &typast.Context{
		Destination: "folder2/dest2",
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
						Type: &typast.FuncDecl{Name: "NewObject"},
					},
				},
			},
		},
	}

	require.NoError(t, ctorAnnot.Annotate(ctx))

	b, _ := ioutil.ReadFile("folder2/dest2/some-target")
	require.Equal(t, `some-template`, string(b))

	require.Equal(t, "Generate @ctor to folder2/dest2/some-target\n", out.String())
}

func TestCtorAnnotation_Annotate_RemoveTargetWhenNoAnnotation(t *testing.T) {
	os.MkdirAll("folder4/pkg4", 0777)
	defer os.RemoveAll("folder4")

	ctorAnnot := &typapp.CtorAnnotation{Target: "some-target"}
	ctx := &typast.Context{
		Destination: "folder4/pkg4",
		Context:     &typgo.Context{},
		Summary:     &typast.Summary{},
	}

	ioutil.WriteFile("folder4/pkg4/some-target", []byte("some-content"), 0777)
	require.NoError(t, ctorAnnot.Annotate(ctx))
	_, err := os.Stat("folder4/pkg4/some-target")
	require.True(t, os.IsNotExist(err))
}

func TestCtorAnnotation_IsCtor(t *testing.T) {
	testcases := []struct {
		TestName string
		*typapp.CtorAnnotation
		Annot    *typast.Annot
		Expected bool
	}{
		{
			TestName:       "private function",
			CtorAnnotation: &typapp.CtorAnnotation{TagName: "@ctor"},
			Annot: &typast.Annot{
				TagName: "@ctor",
				Decl: &typast.Decl{
					Type: &typast.FuncDecl{Name: "someFunction"},
				},
			},
			Expected: false,
		},
		{
			TestName:       "public function",
			CtorAnnotation: &typapp.CtorAnnotation{TagName: "@ctor"},
			Annot: &typast.Annot{
				TagName: "@ctor",
				Decl: &typast.Decl{
					Type: &typast.FuncDecl{Name: "SomeFunction"},
				},
			},
			Expected: true,
		},
		{
			TestName:       "not function",
			CtorAnnotation: &typapp.CtorAnnotation{TagName: "@ctor"},
			Annot: &typast.Annot{
				TagName: "@ctor",
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
			CtorAnnotation: &typapp.CtorAnnotation{TagName: "@ctor"},
			Annot: &typast.Annot{
				TagName: "@ctor",
				Decl: &typast.Decl{
					Type: &typast.FuncDecl{Name: "SomeMethod", Recv: &typast.FieldList{}},
				},
			},
			Expected: false,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			require.Equal(t, tt.Expected, tt.IsCtor(tt.Annot))
		})
	}
}
