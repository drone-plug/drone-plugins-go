package plug

import (
	"bytes"
	"strings"

	"github.com/go-pa/fenv"
	"github.com/olekukonko/tablewriter"
)

func (s *Service) envUsage() {
	s.es.VisitAll(func(f fenv.EnvFlag) {
		{
			s.log.Println(fmtDroneYMLName(f.Name), f.Flag.Value.String())
		}
	})
}

func fmtDroneYMLName(envname string) string {
	stripped := strings.ToLower(envname)
	stripped = strings.TrimPrefix(stripped, "plugin_")
	stripped = strings.TrimPrefix(stripped, "drone_")
	// todo: handle plugin_ and drone_ in some way
	return stripped

}

func (s *Service) usageFuncYml() {
	var b bytes.Buffer
	w := tablewriter.NewWriter(&b)
	// w.Init(&b, 30, 8, 3, ' ', 0)

	w.SetBorder(false)
	w.SetColumnSeparator("")
	w.SetColWidth(120)
	w.SetAutoMergeCells(true)
	w.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})

	sep := func() {
		w.Append([]string{""})
	}
	add := func(value ...string) {
		args := []string{"", value[0] + ":"}
		args = append(args, value[1:]...)

		w.Append(args)
	}

	writeUsage := func(e fenv.EnvFlag) {

		var pluginNames, rawNames []string
	nameLoop:
		for _, v := range e.Names {
			if strings.HasPrefix(v, "DRONE_") {
				continue nameLoop
			}
			if strings.HasPrefix(v, "PLUGIN_") {
				pluginNames = append(pluginNames, fmtDroneYMLName(v))
				continue nameLoop
			}
			rawNames = append(rawNames, fmtDroneYMLName(v))
		}
		if len(pluginNames) == 0 && len(rawNames) == 0 {
			return
		}
		setName := fmtDroneYMLName(e.Name)

		sep()

		if len(pluginNames) > 0 {
			w.Append([]string{strings.Join(pluginNames, ", "), "", e.Flag.Usage})
			// add("option name", strings.Join(pluginNames, ", "))
		}
		if len(rawNames) > 0 {
			add("envvar name", strings.Join(rawNames, ", "))
			// w.Append([]string{strings.Join(pluginNames, ", "), "", ""})
		}

		// w.Append([]string{"", "", e.Flag.Usage})
		if e.IsSelfSet {
			add("set by", setName)
		} else if e.IsSet {
			add("set by flag", e.Flag.Name)
		}

		if s.debug {
			if e.Value != "" {
				add("env value", e.Value)
			}
		}
		if e.Flag.Value.String() != "" {
			add("value", e.Flag.Value.String())
		}
		if e.Err != nil {
			add("**ERROR**", e.Err.Error())
		}
		if errs, ok := s.usageErrors[e.Flag.Name]; ok {
			add("**USAGE ERROR**", strings.Join(errs, "\n"))
		}
	}
	var (
		unsetFlags []fenv.EnvFlag // flags which are not set
		defFlags   []fenv.EnvFlag // flags which are set to their default value
		setFlags   []fenv.EnvFlag // flags which are set
		errFlags   []fenv.EnvFlag // flags which failed to set due to parsing or validation errors
	)

	s.es.VisitAll(func(e fenv.EnvFlag) {
		if e.Err != nil || len(s.usageErrors[e.Flag.Name]) > 0 {
			errFlags = append(errFlags, e)
			return
		}
		if !e.IsSet {
			if e.Flag.Value.String() != "" {
				defFlags = append(defFlags, e)
				return
			}
			unsetFlags = append(unsetFlags, e)
			return
		}
		setFlags = append(setFlags, e)
	})

	debugHeader := func(v string) {
		if s.debug {
			w.Append([]string{v, "----------"})
		}
	}
	if len(unsetFlags) > 0 {
		debugHeader("UNSET")
		for _, f := range unsetFlags {
			writeUsage(f)
		}
	}

	if len(defFlags) > 0 {
		debugHeader("DEFAULT")
		for _, f := range defFlags {
			writeUsage(f)
		}
	}
	if len(setFlags) > 0 {
		debugHeader("SET")
		for _, f := range setFlags {
			writeUsage(f)
		}
	}
	if len(errFlags) > 0 {
		debugHeader("ERRORS")
		for _, f := range errFlags {
			writeUsage(f)
		}
	}

	s.log.Println("plugin usage:")

	w.Render()
	s.log.Println("\n" + b.String())

}
