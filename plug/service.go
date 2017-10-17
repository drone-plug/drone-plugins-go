package plug

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-pa/fenv"
	"github.com/joho/godotenv"
)

type Runner interface {
	SetFlags(fs *FlagSet)
	Exec(ctx context.Context, log *Logger) error
}

func Run(r Runner) {
	s := NewService()
	s.Run(r)
}

type Service struct {
	envFunc  func() map[string]string // function to provide the environment
	argsFunc func() []string          // function to provide the os.Args for parsing the flagset

	hasInit         bool // true if Service.init has been run
	fs              *flag.FlagSet
	es              *fenv.EnvSet        // fenv.Envset for fs
	pfs             *FlagSet            // FlagSet for fs
	usageErrors     map[string][]string // errors registerd by logger
	log             *Logger
	debug           bool  // plugin debug mode
	asPlugin        bool  //true when DRONE=true (environment is drone), swithces display
	continueOnError bool  // if set to true the process does not exit on usage or command error
	execErr         error // the error which can be retreived using the Err() method if  continueOnError after Run if continueOnError is enabled.

}

// ServiceOption is used to configure services with the NewService() function.
type ServiceOption func(s *Service)

// SetFlagSet is a NewService option to use a flag.FlagSet other than flag.CommandLine
func SetFlagSet(fs *flag.FlagSet) ServiceOption {
	if fs == nil {
		log.Fatal("FlagSet is nil")
	}
	return func(s *Service) {
		s.fs = fs
	}
}

// SetEnvFunc is a NewService option to use function to return the environment instead of os.Environ()
func SetEnvFunc(fn func() map[string]string) ServiceOption {
	if fn == nil {
		log.Fatal("EnvFunc is nil")
	}
	return func(s *Service) {
		s.envFunc = fn
	}
}

// SetArgsFunc is NewService option to set command line args instead of os.Args
func SetArgsFunc(fn func() []string) ServiceOption {
	if fn == nil {
		log.Fatal("ArgsFunc is nil")
	}
	return func(s *Service) {
		s.argsFunc = fn
	}
}

// SetLogger is NewService option to set a log.Logger instead of using the default logger in the log package.
func SetLogger(l *log.Logger) ServiceOption {
	if l == nil {
		log.Fatal("Logger is nil")
	}
	return func(s *Service) {
		if s.log != nil {
			s.log.logger = l
			return
		}
		s.log = &Logger{
			logger: l,
		}
	}
}

func ContinueOnError() ServiceOption {
	return func(s *Service) {
		s.continueOnError = true
	}
}

// temporary way to construct services w. special env
func NewService(opts ...ServiceOption) *Service {
	s := &Service{}
	for _, o := range opts {
		if o == nil {
			panic("NewService was given a nil ServiceOption")
		}
		o(s)
	}
	return s
}

/// Run runs the service
func (s *Service) Run(r Runner) {
	s.init()
	env := s.envFunc()
	pfs := &FlagSet{FlagSet: s.fs, es: s.es}
	r.SetFlags(pfs)

	s.asPlugin = env["DRONE"] == "true"
	s.debug = env["PLUGIN_PLUGIN_DEBUG"] != ""
	{
		logFlags := 0
		if s.debug {
			logFlags = log.Ltime | log.Lshortfile
		}
		if s.log.logger == nil {
			log.SetFlags(logFlags)
		} else {
			s.log.logger.SetFlags(logFlags)
		}
	}
	if s.debug {
		s.log.Debugln("drone plugins debug mode is active!")
	}
	if s.debug {
		for k, v := range env {
			if strings.HasPrefix(k, "PLUGIN_") || strings.HasPrefix(k, "DRONE_") {
				s.log.Debugf("[env] %s=%s", k, v)
			}
		}
		s.es.VisitAll(func(e fenv.EnvFlag) {
			s.log.Debugf("[assign] flag '%s' for env vars: %s",
				e.Flag.Name, strings.Join(e.Names, ", "))
		})
	}

	if pfs.envFilesActive {
		s.readEnvfiles(env, pfs.envFiles)
	}

	if s.asPlugin {
		s.fs.Usage = s.usageFuncYml
		s.fs.Init(s.args()[0], flag.ContinueOnError)
	}

	if err := s.es.ParseEnv(env); err != nil {
		s.execErr = err
		s.fs.Usage()
		if !s.continueOnError {
			os.Exit(1)
		}
		return

	}

	if err := s.fs.Parse(s.args()[1:]); err != nil {
		s.execErr = err
		s.log.Println(err)
		if !s.continueOnError {
			os.Exit(1)
		}
		return

	}
	if s.debug {
		s.es.VisitAll(func(e fenv.EnvFlag) {
			if !e.IsSelfSet && e.IsSet {
				s.log.Debugf("[flag] '%s' set: %v", e.Flag.Name, e.Flag.Value)
			}
		})
		s.es.VisitAll(func(e fenv.EnvFlag) {
			if e.IsSelfSet {
				s.log.Debugf("[envflag] '%s' set by env var '%s': %v", e.Flag.Name, e.Name, e.Flag.Value)
			}
		})
		s.usageFuncYml()
	}
	ctx := context.Background()
	s.log.Debugln("------ executing plugin func  -----")
	err := r.Exec(ctx, s.log)
	s.log.Debugln("------ plugin func done  -----")
	s.execErr = err
	var hasErrors bool
	if err != nil {
		s.log.Debugln("ErrUsageError returned")
		if s.debug {
			_ = s.log.Output(2, fmt.Sprintf("plugin runner error: %v", err))
		}
		if err == ErrUsageError {
			hasErrors = true
		}
	} else {
		s.es.VisitAll(func(e fenv.EnvFlag) {
			if e.Err != nil {
				hasErrors = true
			}
		})
	}
	if hasErrors {
		s.fs.Usage()
		if !s.continueOnError {
			os.Exit(1)
		}
		return
	}
}

// init various internal variables and sets default values
func (s *Service) init() {
	if s.hasInit {
		return
	}
	s.hasInit = true
	s.usageErrors = make(map[string][]string)
	if s.fs == nil {
		s.fs = flag.CommandLine
	}
	s.es = fenv.NewEnvSet(s.fs, fenv.Prefix("plugin_"))
	if s.envFunc == nil {
		s.envFunc = fenv.OSEnv
	}
	if s.log == nil {
		s.log = &Logger{}
	}
	s.log.s = s

}

func (s *Service) readEnvfiles(env map[string]string, defaultValue []string) {
	// Special handling of the -env-file flag to load it before the main
	// flagset/envsets are parsed.
	var envfiles []string
	envfiles = append(envfiles, defaultValue...)
	s.log.Debugln("[envfile] read env files")
	fs := flag.NewFlagSet("envfile", flag.ContinueOnError)

	es := fenv.NewEnvSet(fs, fenv.Prefix("plugin_"), fenv.ContinueOnError())
	envfiles = append(envfiles)
	fs.Var((*stringSliceFlag)(&envfiles), envfileFlagName, "source env file")
	if err := es.ParseEnv(env); err != nil {
		s.log.Fatal(err)
	}

	{
		var args []string
	loop:
		for _, arg := range s.args()[1:] {
			if len(args) == 0 && arg == "-"+envfileFlagName {
				args = append(args, arg)
				continue loop
			}

			if len(args) == 1 {
				args = append(args, arg)
				break loop
			}

		}
		if err := fs.Parse(args); err != nil {
			s.log.Fatal(err)
		}
	}

	if len(envfiles) > 0 {
	files:
		for _, filename := range envfiles {
			s.log.Debugf("[envfile] loading env file: %v", filename)
			e, err := godotenv.Read(filename)
			if err != nil {
				s.log.Debugf("[envfile] error loading env file: %v", err)
				continue files
			}
			for k, v := range e {
				if _, ok := env[k]; !ok {
					s.log.Debugf("[envfile] setting %s=%s", k, v)
					env[k] = v
				} else {
					s.log.Debugf("[envfile] skipping already defined var: %s=%s", k, v)
				}
			}
		}
	}

}

func (s *Service) args() []string {
	if s.argsFunc != nil {
		return s.argsFunc()
	}
	return os.Args
}

func (s *Service) parse() error {
	return nil
}

// Err returns an *ExitError struct consisting of various error information
func (s *Service) Err() error {
	if s.execErr == nil && len(s.usageErrors) == 0 {
		return nil
	}
	return &ExecError{
		Err:         s.execErr,
		UsageErrors: s.usageErrors,
	}
}
