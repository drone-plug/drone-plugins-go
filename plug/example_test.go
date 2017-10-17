package plug_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/drone-plug/drone-plugins-go/plug"
)

// Plugin .
type Plugin struct {
	Name        string
	AccessKey   string
	Flowers     []string // []string flag type
	MoreFlowers []string
	Dogs        map[string]string // Drone specific map[string]string
	Cats        map[string]string // for the JSON encoded input.

	Commit      plug.Commit // This plugin need all of drones Commit metadata
	BuildNumber int64       // and the build number
	Event       string      // + event.
}

// NewPlugin returns a Plugin with default values set
func NewPlugin() *Plugin {
	return &Plugin{
		// Default values for string slice types.
		MoreFlowers: []string{"default", "value"},
		// Default values for string map types.
		Cats: map[string]string{"default": "value"},
		// Default values for string map types.
		Dogs: map[string]string{"default": "value"},
	}
}

// SetFlags implments plug.Runner
func (p *Plugin) SetFlags(fs *plug.FlagSet) {
	// Use the regular flag package.
	fs.StringVar(&p.Name, "name", "default value", "your name")
	// PLUGIN_ACCESS_KEY env var is automatic.
	fs.StringVar(&p.AccessKey, "access_key", "", "")
	// Alternative names can be specified. emtpy string is converted to the auto generated name.
	fs.Env(&p.AccessKey, "", "aws_access_key")
	fs.StringSliceVar(&p.Flowers, "flowers", "list of flowers")
	fs.StringSliceVar(&p.MoreFlowers, "more.flowers", "")
	fs.StringMapVar(&p.Cats, "cats", "")
	fs.StringMapVar(&p.Dogs, "dogs", "")
	// Register all fields for Repo, Build or Commit
	fs.CommitVar(&p.Commit)
	// or single properties...
	fs.BuildEventVar(&p.Event)
	fs.BuildNumberVar(&p.BuildNumber)
}

// Exec implments plug.Runner
func (p *Plugin) Exec(ctx context.Context, log *plug.Logger) error {
	if p.Name == "My name" {
		log.Usagef(&p.Name, "that is not your name")
		// prints: invalid value for 'name': that is not your name
		// return plug.ErrUsageError
	}
	log.Debugln("this is only printed when the plugin is run in debug mode")
	log.Println("normal plugin output")

	after, _ := json.MarshalIndent(&p, "", " ")
	fmt.Println("build number:", p.BuildNumber)
	fmt.Println("commit message:", p.Commit.Message)
	fmt.Println("after:", string(after))
	return nil

}

func Example_plugin() {
	p := NewPlugin()

	// set some values in the env for testing the flags API below.
	for k, v := range map[string]string{
		"PLUGIN_NAME":          "My name",
		"PLUGIN_MORE_FLOWERS":  "one,two,three",
		"AWS_ACCESS_KEY":       "such access key",
		"PLUGIN_DOGS":          `{"key":"value"}`,
		"DRONE_COMMIT_MESSAGE": "Added text",
		"DRONE_BUILD_NUMBER":   "12",
	} {
		_ = os.Setenv(k, v)
	}

	plug.Run(p)
	// Output:
	// build number: 12
	// commit message: Added text
	// after: {
	//  "Name": "My name",
	//  "AccessKey": "such access key",
	//  "Flowers": null,
	//  "MoreFlowers": [
	//   "one",
	//   "two",
	//   "three"
	//  ],
	//  "Dogs": {
	//   "key": "value"
	//  },
	//  "Cats": {
	//   "default": "value"
	//  },
	//  "Commit": {
	//   "Sha": "",
	//   "Ref": "",
	//   "Link": "",
	//   "Branch": "",
	//   "Message": "Added text",
	//   "Author": {
	//    "Name": "",
	//    "Email": "",
	//    "Avatar": ""
	//   }
	//  },
	//  "BuildNumber": 12,
	//  "Event": ""
	// }

}
