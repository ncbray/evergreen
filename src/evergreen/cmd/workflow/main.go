// Tool for developing Evergreen.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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
	ctx.SimpleCommand("go", "install", "evergreen/cmd/regenerate")
	if ctx.Errored {
		return
	}
	Sync(ctx)
}

func Sync(ctx *Context) {
	ctx.Step("Regenerating sources")
	ctx.SimpleCommand("bin/regenerate", "-replace")
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
	ctx.SimpleCommand("go", "run", "src/evergreen/cmd/regenerate/main.go")
	if ctx.Errored {
		return
	}
	ctx.Step("Running tests")
	ctx.SimpleCommand("go", "test", "./...")
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
