package main

import (
	"testing"

	"github.com/drone-plug/drone-plugins-go/plug/plugtest"
)

func TestExecFail(t *testing.T) {
	p := &Plugin{}
	pt := plugtest.New(t, p)
	pt.AssertFail()
}

func TestExecOK(t *testing.T) {
	p := &Plugin{}
	pt := plugtest.New(t, p)

	pt.SetPluginVars(map[string]string{
		"token": "asdf",
	})
	pt.SetVars(map[string]string{
		"downstream_server": "asdf",
	})
	pt.AssertSuccess()
	pt.AssertOutput(`success!
`)
}

func TestOutputs(t *testing.T) {
	// this test is mostly here for testing usage formatting output, set silece to false to always print output
	// const silence = false
	const silence = true

	type testCase struct {
		plugvars map[string]string
		envvars  map[string]string
		fail     bool
	}

	tests := []testCase{
		{
			fail: true,
		},
		{
			plugvars: map[string]string{
				"version": "asd",
			},
			fail: true,
		},
		{
			plugvars: map[string]string{
				"plugin_debug": "true",
				"version":      "1.0",
			},
			fail: true,
		},
		{
			plugvars: map[string]string{
				"token":  "tokenvalue",
				"server": "servervalue",
			},
			// fail: true,
		},
		{
			plugvars: map[string]string{
				"token":  "tokenvalue",
				"server": "servervalue",
				"fork":   "koo",
			},
			fail: true,
		},
	}

ts:
	for _, tc := range tests {
		p := &Plugin{}
		pt := plugtest.New(t, p)
		pt.SetPluginVars(tc.plugvars)
		pt.SetVars(tc.envvars)
		_ = pt.Run()

		if tc.fail {
			pt.AssertFail()
		} else {
			pt.AssertSuccess()
		}
		if silence {
			continue ts
		}
		t.Logf(`

*************************************
  args: none
  plugin vars: %v
  env vars: %v
  err: %v
  output:

%s`,
			tc.plugvars,
			tc.envvars,
			pt.Err,
			pt.Output())

	}

}
