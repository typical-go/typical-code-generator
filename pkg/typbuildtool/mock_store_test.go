package typbuildtool_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
)

func TestMockStore(t *testing.T) {
	store := typbuildtool.NewMockStore()
	store.Put(&typbuildtool.MockTarget{MockDir: "pkg1", SrcName: "target1"})
	store.Put(&typbuildtool.MockTarget{MockDir: "pkg1", SrcName: "target2"})
	store.Put(&typbuildtool.MockTarget{MockDir: "pkg2", SrcName: "target3"})
	store.Put(&typbuildtool.MockTarget{MockDir: "pkg1", SrcName: "target4"})
	store.Put(&typbuildtool.MockTarget{MockDir: "pkg1", SrcName: "target5"})
	store.Put(&typbuildtool.MockTarget{MockDir: "pkg2", SrcName: "target6"})

	require.Equal(t, map[string][]*typbuildtool.MockTarget{
		"pkg1": []*typbuildtool.MockTarget{
			{MockDir: "pkg1", SrcName: "target1"},
			{MockDir: "pkg1", SrcName: "target2"},
			{MockDir: "pkg1", SrcName: "target4"},
			{MockDir: "pkg1", SrcName: "target5"},
		},
		"pkg2": []*typbuildtool.MockTarget{
			{MockDir: "pkg2", SrcName: "target3"},
			{MockDir: "pkg2", SrcName: "target6"},
		},
	}, store.Map())
}
