// Tool for developing Evergreen.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var dubsrc = "dubsrc"

var outdir = "output"

type Mode struct {
	Name string
	Run  func(*Context)
	Help string
}

type Context struct {
	Modes   []*Mode
	Errored bool
}

func (ctx *Context) GetMode(requested string) *Mode {
	for _, mode := range ctx.Modes {
		if mode.Name == requested {
			return mode
		}
	}
	return nil
}

func (ctx *Context) Step(description string) {
	fmt.Println("###", description, "###")
}

func (ctx *Context) SimpleCommand(name string, args ...string) {
	full := []string{name}
	full = append(full, args...)
	fmt.Println(strings.Join(full, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ctx.CheckError(cmd.Run())
}

func (ctx *Context) EnvCommand(env []string, name string, args ...string) {
	full := []string{name}
	full = append(full, args...)
	fmt.Println(strings.Join(full, " "))
	cmd := exec.Command(name, args...)
	cmd.Env = append(env, os.Environ()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ctx.CheckError(cmd.Run())
}

func (ctx *Context) CheckError(err error) {
	if err != nil {
		fmt.Println("ERROR", err)
		ctx.Errored = true
	}
}

func Build(ctx *Context) {
	ctx.Step("Vetting sources")
	ctx.SimpleCommand("go", "vet", "./...")
	if ctx.Errored {
		return
	}
	ctx.Step("Building binary")
	ctx.SimpleCommand("go", "install", "evergreen/cmd/egc")
	if ctx.Errored {
		return
	}
	Sync(ctx)
}

func Sync(ctx *Context) {
	ctx.Step("Regenerating sources")
	ctx.SimpleCommand("bin/egc", "-indir="+filepath.Join(dubsrc, "evergreen"), "-outdir=src", "-gopackage=evergreen")
	if ctx.Errored {
		return
	}
	ctx.Step("Formatting sources")
	ctx.SimpleCommand("go", "fmt", "./...")
}

func Help(ctx *Context) {
	fmt.Printf("Usage: workflow mode [<args>]\n\n")

	fmt.Println("Available modes")
	for _, mode := range ctx.Modes {
		fmt.Printf("    %-10s %s\n", mode.Name, mode.Help)
	}
}

func Test(ctx *Context) {
	ctx.Step("Cleaning generated sources")
	ctx.CheckError(os.RemoveAll("src/generated"))
	if ctx.Errored {
		return
	}
	ctx.Step("Generating sources for testing")
	ctx.SimpleCommand("go", "run", "src/evergreen/cmd/egc/main.go", "-indir="+filepath.Join(dubsrc, "evergreen"), "-outdir=src", "-gopackage=generated", "-gentests")
	if ctx.Errored {
		return
	}
	ctx.Step("Running tests")
	ctx.SimpleCommand("go", "test", "./...")
}

func Dump(ctx *Context) {
	ctx.Step("Removing old outputs")
	ctx.CheckError(os.RemoveAll(outdir))
	if ctx.Errored {
		return
	}
	ctx.Step("Dumping graphs")
	ctx.SimpleCommand("go", "run", "src/evergreen/cmd/egc/main.go", "-indir="+filepath.Join(dubsrc, "evergreen"), "-outdir=src", "-gopackage=generated", "-dump", "-v=1")
	if ctx.Errored {
		return
	}
}

func Crossbuild(ctx *Context) {
	builddir := filepath.Join(outdir, "crossbuild")
	ctx.Step("Building binary")
	ctx.EnvCommand(
		[]string{"GOOS=linux", "GOARCH=arm", "GOARM=7"},
		"go", "build", "-o", filepath.Join(builddir, "egc"), "evergreen/cmd/egc")
	if ctx.Errored {
		return
	}
	ctx.Step("Copying sources")
	ctx.SimpleCommand("cp", "-r", filepath.Join(dubsrc, "evergreen"), builddir)
	if ctx.Errored {
		return
	}
	host_var := "CROSSBUILD_HOST"
	host := os.Getenv(host_var)
	if host != "" {
		ctx.Step("Uploading")
		ctx.SimpleCommand("scp", "-r", builddir, host+":~/crossbuild")
		ctx.Step("Profiling")
		ctx.SimpleCommand("ssh", host, "cd crossbuild;./egc -indir evergreen -outdir output -gopackage generated -cpuprofile=egc_cpu.prof")
		ctx.SimpleCommand("scp", host+":~/crossbuild/egc_cpu.prof", builddir)
	} else {
		fmt.Printf("%s not specified, skipping upload.\n", host_var)
	}
}

func main() {
	ctx := &Context{
		Modes: []*Mode{
			&Mode{
				Name: "build",
				Run:  Build,
				Help: "Regenerate the checked-in source code.",
			},
			&Mode{
				Name: "help",
				Run:  Help,
				Help: "Displays this message.",
			},
			&Mode{
				Name: "sync",
				Run:  Sync,
				Help: "Regenerate the checked-in source without trying to rebuild the program.",
			},
			&Mode{
				Name: "test",
				Run:  Test,
				Help: "Runs tests.",
			},
			&Mode{
				Name: "dump",
				Run:  Dump,
				Help: "Dumps graphs of intermediate data structures.",
			},
			&Mode{
				Name: "crossbuild",
				Run:  Crossbuild,
				Help: "Build for ARM Linux.",
			},
		},
	}

	requested := "help"
	if len(os.Args) >= 2 {
		requested = os.Args[1]
	}

	mode := ctx.GetMode(requested)
	if mode == nil {
		fmt.Printf("Unreconized mode: %s\n\n", requested)
		ctx.Errored = true
		mode = ctx.GetMode("help")
		if mode == nil {
			panic(requested)
		}
	}
	mode.Run(ctx)

	if ctx.Errored {
		os.Exit(1)
	}
}
