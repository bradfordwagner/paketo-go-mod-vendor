package gomodvendor

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

//go:generate faux --interface Executable --output fakes/executable.go
type Executable interface {
	Execute(pexec.Execution) error
}

type ModVendor struct {
	executable Executable
	logs       scribe.Emitter
	clock      chronos.Clock
}

func NewModVendor(executable Executable, logs scribe.Emitter, clock chronos.Clock) ModVendor {
	return ModVendor{
		executable: executable,
		logs:       logs,
		clock:      clock,
	}
}

func (m ModVendor) ShouldRun(workingDir string) (bool, string, error) {
	ok, err := fs.Exists(filepath.Join(workingDir, "vendor"))
	if err != nil {
		return false, "", err
	}
	if ok {
		return false, "modules are already vendored", nil
	}

	return true, "", nil
}

func (m ModVendor) Execute(buildConfiguration BuildConfiguration, path, workingDir string) error {
	args := []string{"mod", "vendor"}

	m.logs.Process("Executing build process")

	env := m.configEnv(buildConfiguration, workingDir)
	m.logs.Subprocess("Env='%s'", strings.Join(env, ","))
	m.logs.Subprocess("Running 'go %s'", strings.Join(args, " "))

	duration, err := m.clock.Measure(func() error {
		return m.executable.Execute(pexec.Execution{
			Args:   args,
			Env:    env,
			Dir:    workingDir,
			Stdout: m.logs.ActionWriter,
			Stderr: m.logs.ActionWriter,
		})
	})
	if err != nil {
		m.logs.Action("Failed after %s", duration.Round(time.Millisecond))
		return err
	}

	m.logs.Action("Completed in %s", duration.Round(time.Millisecond))
	m.logs.Break()

	return nil
}

func (m ModVendor) configEnv(buildConfiguration BuildConfiguration, path string) (env []string) {
	env = append(env, fmt.Sprintf("GOMODCACHE=%s", path))
	if buildConfiguration.GoProxy != "" {
		env = append(env, fmt.Sprintf("GOPROXY=%s", buildConfiguration.GoProxy))
		m.logs.Subprocess(fmt.Sprintf("GOPROXY=%s", buildConfiguration.GoProxy))
	}
	if buildConfiguration.GoNoProxy != "" {
		env = append(env, fmt.Sprintf("GONOPROXY=%s", buildConfiguration.GoNoProxy))
		m.logs.Subprocess(fmt.Sprintf("GONOPROXY=%s", buildConfiguration.GoNoProxy))
	}
	if buildConfiguration.GoSumDB != "" {
		env = append(env, fmt.Sprintf("GOSUMDB=%s", buildConfiguration.GoSumDB))
		m.logs.Subprocess(fmt.Sprintf("GOSUMDB=%s", buildConfiguration.GoSumDB))
	}
	if buildConfiguration.GoNoSumDB != "" {
		env = append(env, fmt.Sprintf("GONOSUMDB=%s", buildConfiguration.GoNoSumDB))
		m.logs.Subprocess(fmt.Sprintf("GONOSUMDB=%s", buildConfiguration.GoNoSumDB))
	}
	if buildConfiguration.GoPrivate != "" {
		env = append(env, fmt.Sprintf("GOPRIVATE=%s", buildConfiguration.GoPrivate))
		m.logs.Subprocess(fmt.Sprintf("GOPRIVATE=%s", buildConfiguration.GoPrivate))
	}
	return
}
