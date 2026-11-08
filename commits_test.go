package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCommitInfo(t *testing.T) {
	testCases := []struct {
		desc  string
		msg   string
		typ   string
		scope string
		err   bool
	}{
		{
			desc:  "type only",
			msg:   "feat: Add something",
			typ:   "feat",
			scope: "",
			err:   false,
		}, {
			desc:  "with scope",
			msg:   "refactor(backend): Restructure services",
			typ:   "refactor",
			scope: "backend",
			err:   false,
		},
		{
			desc:  "multi line",
			msg:   "chore(deps): Bump actions/setup-go\n\n Latest release!",
			typ:   "chore",
			scope: "deps",
			err:   false,
		},
		{
			desc: "without colon",
			msg:  "feat; Add something",
			err:  true,
		},
		{
			desc: "mismatching parenthes",
			msg:  "tests(frontend0: Remove unneccesary tests",
			err:  true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			info, err := ParseCommitInfo(tC.msg)
			if tC.err {
				require.NotNil(t, err)
				return
			}
			assert.Equal(t, tC.typ, info.typ)
			assert.Equal(t, tC.scope, info.scope)
		})
	}
}
