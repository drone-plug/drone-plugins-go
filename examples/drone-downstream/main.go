package main

import (
	"context"

	"github.com/drone-plug/drone-plugins-go/plug"
)

var build = "0" // build number set at compile-time

// Plugin defines the Downstream plugin parameters.
type Plugin struct {
	Repos         []string
	Server        string
	Token         string
	Fork          bool
	Build         plug.Build
	AnotherOption string
}

func (p *Plugin) SetFlags(fs *plug.FlagSet) {
	fs.BoolVar(&p.Fork, "fork", false, "Trigger a new build for a repository")
	fs.StringSliceVar(&p.Repos, "repositories", "List of repositories to trigger")
	fs.StringVar(&p.Server, "server", "", "Trigger a drone build on a custom server")
	fs.Env(&p.Server, "", "plugin_server2", "downstream_server", "downstream_server2") // emtpy string means default PLUGIN_...
	fs.StringVar(&p.Token, "token", "", "Drone API token from your user settings")
	fs.Env(&p.Token, "downstream_token", "")
	fs.BuildVar(&p.Build)
	fs.EnvFiles()
	fs.StringVar(&p.AnotherOption, "another-option", "", "option without PLUGIN_ name")
	fs.Env(&p.AnotherOption, "another_option")

}

// Exec runs the plugin
func (p *Plugin) Exec(ctx context.Context, log *plug.Logger) error {
	isValid := true
	if len(p.Token) == 0 {
		log.Usageln(&p.Token, "you must provide your Drone access token.")
		log.Debugln("not valid")
		isValid = false
	}
	if len(p.Server) == 0 {
		log.Usageln(&p.Server, "you must provide your Drone server.")
		log.Debugln("not valid")
		isValid = false
	}
	if !isValid {
		return plug.ErrUsageError
	}
	//client := drone.NewClientToken(p.Server, p.Token)
	// ...
	log.Println("success!")
	return nil
}

func main() {
	var p Plugin
	plug.Run(&p)
}
