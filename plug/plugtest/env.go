package plugtest

import "strings"

func (t *PT) SetDebug() {
	t.T.Helper()
	t.before()
	t.env["PLUGIN_PLUGIN_DEBUG"] = "true"
}

func (t *PT) SetPluginVars(o map[string]string) {
	t.T.Helper()
	t.before()
	for k, v := range o {
		t.env[strings.ToUpper("plugin_"+k)] = v
	}
}

func (t *PT) SetVars(o map[string]string) {
	t.T.Helper()
	t.before()
	for k, v := range o {
		t.env[strings.ToUpper(k)] = v
	}
}

func (t *PT) envFunc() map[string]string {
	return t.env
}
