package plugtest

import (
	"bytes"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/drone-plug/drone-plugins-go/plug"
)

// T .
type PT struct {
	env    map[string]string
	T      *testing.T
	R      plug.Runner
	logbuf *bytes.Buffer
	Err    error // error from service.Run
	hasRun bool
}

func New(t *testing.T, r plug.Runner) *PT {
	if t == nil {
		log.Fatal("t can not be nil")
	}
	if r == nil {
		t.Fatal("Runner can not be nil")
	}
	pt := &PT{
		env:    make(map[string]string),
		T:      t,
		R:      r,
		logbuf: &bytes.Buffer{},
	}
	pt.SetVars(map[string]string{
		"drone": "true",
	})
	return pt
}

func (t *PT) Run() error {
	// t.T.Helper()
	if t.hasRun {
		t.T.Fatal("this test has already run")
	}
	t.hasRun = true
	log := log.New(t.logbuf, "", 0)
	s := plug.NewService(
		plug.SetFlagSet(flag.NewFlagSet("-", flag.ContinueOnError)),
		plug.SetEnvFunc(t.envFunc),
		plug.SetArgsFunc(func() []string { return []string{os.Args[0]} }),
		plug.SetLogger(log),
		plug.ContinueOnError(),
	)
	s.Run(t.R)
	t.Err = s.Err()
	return t.Err

}

func (t *PT) Output() string {
	// todo: document that buffer is emtpied
	t.after()
	return t.logbuf.String()
}

// after ensures that Run has been called.
func (t *PT) after() {
	t.T.Helper()
	if !t.hasRun {
		_ = t.Run()
	}
}

// before ensures that run as not been called. stops the test
func (t *PT) before() {
	t.T.Helper()
	if t.hasRun {
		t.T.Fatal("has already run")
	}
}
