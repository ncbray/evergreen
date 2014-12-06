package main

import (
	"evergreen/base"
	"evergreen/dub/flow"
	"evergreen/dub/transform"
	"evergreen/dub/transform/golang"
	"evergreen/dub/tree"
	"evergreen/framework"
	gotree "evergreen/go/tree"
	"evergreen/io"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
)

func dumpProgram(runner *framework.TaskRunner, program []*flow.DubPackage) {
	for _, dubPkg := range program {
		for _, f := range dubPkg.Funcs {
			styler := &flow.DotStyler{Decl: f}
			dot := base.GraphToDot(f.CFG, styler)
			parts := []string{"output"}
			parts = append(parts, dubPkg.Path...)
			parts = append(parts, fmt.Sprintf("%s.svg", f.Name))
			outfile := filepath.Join(parts...)

			runner.Run(func() {
				// Dump flowgraph
				io.WriteDot(dot, outfile)
			})
		}
	}
}

func analyizeProgram(program []*flow.DubPackage) {
	for _, dubPkg := range program {
		for _, s := range dubPkg.Structs {
			if s.Implements != nil {
				s.Implements.IsParent = true
			}
		}
	}
}

func GenerateGo(program []*flow.DubPackage, runner *framework.TaskRunner) {
	root := "generated"
	if replace {
		root = "evergreen"
	}
	prog := golang.GenerateGo(program, root, !replace)

	// Compact simple expressions back into tree form.
	gotree.Consolidate(prog)

	// Give everything names: variables, etc.
	gotree.Nameify(prog)

	// Generate the sources.
	gotree.OutputProgram(prog, "src", runner)
}

func processProgram(status framework.Status, p framework.LocationProvider, runner *framework.TaskRunner, root string) {
	program, funcs := tree.ParseProgram(status.CreateChild(), p, root)
	if status.ShouldHalt() {
		return
	}
	flowProgram := transform.LowerProgram(program, funcs)

	if dump {
		dumpProgram(runner, flowProgram)
	}

	analyizeProgram(flowProgram)
	GenerateGo(flowProgram, runner)
}

func entryPoint(p framework.LocationProvider, status framework.Status) {
	runner := framework.CreateTaskRunner(8)

	root_dir := "dub"
	processProgram(status, p, runner, root_dir)
	runner.Kill()
}

var dump bool
var replace bool
var cpuprofile string
var memprofile string

func main() {
	flag.BoolVar(&dump, "dump", false, "Dump flowgraphs to disk.")
	flag.BoolVar(&replace, "replace", false, "Replace the existing implementation.")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to file")
	flag.StringVar(&memprofile, "memprofile", "", "write memory profile to this file")
	flag.Parse()

	p := framework.MakeProvider()
	status := framework.MakeStatus(p)

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			status.Error(err.Error())
			return
		} else {
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
	}

	entryPoint(p, status)

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			status.Error(err.Error())
			return
		} else {
			pprof.WriteHeapProfile(f)
			f.Close()
		}
	}

	if status.ShouldHalt() {
		fmt.Printf("%d errors\n", status.ErrorCount())
		os.Exit(1)
	}
}
