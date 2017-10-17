package plugtest

import (
	"encoding/json"
	"log"
	"strings"
)

// env contains environment variables as key/value map instead of a slice
// like the standard library does it.
type Env map[string]string

func (env Env) SetDebug() {
	env["PLUGIN_PLUGIN_DEBUG"] = "true"
}

func (env Env) SetPluginVars(o map[string]string) {
	for k, v := range o {
		env[strings.ToUpper("plugin_"+k)] = v
	}
}

func (env Env) SetVars(o map[string]string) {
	for k, v := range o {
		env[strings.ToUpper(k)] = v
	}
}

func (env Env) EnvFunc() map[string]string {
	return env
}

func JSON(v interface{}) string {
	data, err := json.Marshal(&v)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)

}
