package plugtest

import (
	"bytes"
	"flag"
	"log"
	"testing"

	"github.com/drone-plug/drone-plugins-go/plug"
)

// T .
type PT struct {
	Env
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

	env := make(Env)
	env.SetVars(map[string]string{
		"drone": "true",
	})
	return &PT{
		Env:    env,
		T:      t,
		R:      r,
		logbuf: &bytes.Buffer{},
	}
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
		plug.SetEnvFunc(t.Env.EnvFunc),
		plug.SetArgsFunc(func() []string { return []string{"command"} }),
		plug.SetLogger(log),
		plug.ContinueOnError(),
	)
	s.Run(t.R)
	t.Err = s.Err()
	return t.Err

}

func (t *PT) ensure() {
	if !t.hasRun {
		t.Run()
	}
}
func (t *PT) AssertSuccess() {
	t.ensure()
	if t.Err != nil {
		t.T.Log(t.logbuf.String())
		t.T.Fatal("should have succeeded", t.Err)
	}
}

func (t *PT) AssertFail() {
	t.ensure()
	if t.Err == nil {
		t.T.Log(t.logbuf.String())
		t.T.Fatal("should have failed")
	}

}

func (t *PT) Output() string {
	t.ensure()
	// todo: document that buffer is emtpied
	return t.logbuf.String()
}

func (t *PT) AssertOutput(text string) {
	// todo: document that buffer is emtpied
	out := t.Output()
	if out != text {
		t.T.Fatalf("output not as expected!\n got:\n%s\n expected:\n%s", out, text)
	}
}
