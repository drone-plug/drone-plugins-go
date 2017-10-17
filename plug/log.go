package plug

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/go-pa/fenv"
)

// Logger works like log.Log with additional features for drone plugin specific usage
type Logger struct {
	logger *log.Logger
	s      *Service
	// mu     sync.Mutex
}

// Output forwards logging output to the configured logger or the stdlib log.Output if no custom logger is defined.
func (l *Logger) Output(calldepth int, s string) error {
	if l.logger == nil {
		return log.Output(calldepth+1, s)
	}
	return l.logger.Output(calldepth+1, s)
}

// These functions write to the standard logger.

// Debug calls Output to print to the standard logger if plugins debug is enabled.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Debug(v ...interface{}) {
	if !l.s.debug {
		return
	}
	_ = l.Output(2, fmt.Sprint(v...))
}

// Debugf calls Output to print to the standard logger if plugins debug is enabled.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if !l.s.debug {
		return
	}

	_ = l.Output(2, fmt.Sprintf(format, v...))

}

// Debugln calls Output to print to the standard logger if plugins debug is enabled.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Debugln(v ...interface{}) {
	if !l.s.debug {
		return
	}
	_ = l.Output(2, fmt.Sprintln(v...))
}

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	_ = l.Output(2, fmt.Sprint(v...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	_ = l.Output(2, fmt.Sprintf(format, v...))
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Println(v ...interface{}) {
	_ = l.Output(2, fmt.Sprintln(v...))
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func (l *Logger) Fatal(v ...interface{}) {
	_ = l.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	_ = l.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func (l *Logger) Fatalln(v ...interface{}) {
	_ = l.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// for internal use, triggers if the API is in an bad state
func (l *Logger) programmingFatalf(format string, v ...interface{}) {
	_ = l.Output(2, "programming error: "+fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Usage TODO
func (l *Logger) Usage(ref interface{}, v ...interface{}) {
	if l.s.debug {
		_ = l.Output(2, l.fmtFlagUsage(ref, fmt.Sprint(v...)))
	}
	flg, err := l.findEnvFlag(ref)
	if err != nil {
		_ = l.Output(2, err.Error())
		os.Exit(1)
	}
	errs := l.s.usageErrors[flg.Flag.Name]
	errs = append(errs, fmt.Sprint(v...))
	l.s.log.s.usageErrors[flg.Flag.Name] = errs

}

// Usagef TODO
func (l *Logger) Usagef(ref interface{}, format string, v ...interface{}) {
	if l.s.debug {
		_ = l.Output(2, l.fmtFlagUsage(ref, fmt.Sprintf(format, v...)))
	}
	flg, err := l.findEnvFlag(ref)
	if err != nil {
		_ = l.Output(2, err.Error())
		os.Exit(1)
	}
	errs := l.s.usageErrors[flg.Flag.Name]
	errs = append(errs, fmt.Sprintf(format, v...))
	l.s.usageErrors[flg.Flag.Name] = errs

}

// Usageln TODO
func (l *Logger) Usageln(ref interface{}, v ...interface{}) {
	if l.s.debug {
		_ = l.Output(2, l.fmtFlagUsage(ref, fmt.Sprintln(v...)))
	}
	flg, err := l.findEnvFlag(ref)
	if err != nil {
		_ = l.Output(2, err.Error())
		os.Exit(1)
	}
	errs := l.s.usageErrors[flg.Flag.Name]
	errs = append(errs, fmt.Sprintln(v...))
	l.s.usageErrors[flg.Flag.Name] = errs
}

func (l *Logger) fmtFlagUsage(ref interface{}, rest string) string {
	flg, err := l.findEnvFlag(ref)
	if err != nil {
		l.Debugln(err)
	} else {
		l.Debugf("adding usage error for '%s': %s", flg.Flag.Name, rest)
		l.s.usageErrors[flg.Flag.Name] = append(l.s.usageErrors[flg.Flag.Name], rest)
	}
	fmtname := flg.Flag.Name
	if flg != nil {
		if flg.Name != "" {
			fmtname = fmtDroneYMLName(flg.Name)
		}
	}
	return fmt.Sprintf("plugin option '%s' error: %s", fmtname, rest)
}

func (l *Logger) findEnvFlag(value interface{}) (*fenv.EnvFlag, error) {
	// returns the flag.Flag instace bound to ref or nil if not found
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("not a pointer: %v", value)
	}
	vp := rv.Pointer()
	var flg *fenv.EnvFlag
	l.s.es.VisitAll(func(f fenv.EnvFlag) {
		p := reflect.ValueOf(f.Flag.Value).Pointer()
		if vp == p {
			flg = &f
		}
	})
	return flg, nil

}
