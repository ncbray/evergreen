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
	"runtime"
	"runtime/pprof"
)

func dumpProgram(status framework.PassStatus, runner *framework.TaskRunner, program []*flow.DubPackage) {
	status.Begin()
	defer status.End()

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

func GenerateGo(status framework.PassStatus, program []*flow.DubPackage, runner *framework.TaskRunner) {
	status.Begin()
	defer status.End()

	root := "generated"
	if replace {
		root = "evergreen"
	}
	prog := golang.GenerateGo(status.Pass("generate_go"), program, root, !replace)

	// Compact simple expressions back into tree form.
	gotree.Consolidate(status.Pass("consolidate"), prog)

	// Give everything names: variables, etc.
	gotree.Nameify(status.Pass("nameify"), prog)

	// Generate the sources.
	gotree.OutputProgram(status.Pass("output"), prog, "src", runner)
}

func processProgram(status framework.PassStatus, p framework.LocationProvider, runner *framework.TaskRunner, root string) {
	program, funcs := tree.DubProgramFrontend(status.Pass("dub_frontend"), p, root)
	if status.ShouldHalt() {
		return
	}
	flowProgram := transform.LowerProgram(status.Pass("lower"), program, funcs)

	if dump {
		dumpProgram(status.Pass("dump"), runner, flowProgram)
	}

	analyizeProgram(flowProgram)
	GenerateGo(status.Pass("go_backend"), flowProgram, runner)
}

func entryPoint(p framework.LocationProvider, status framework.PassStatus) {
	status.Begin()
	defer status.End()

	runner := framework.CreateTaskRunner(jobs)

	root_dir := "dub"
	processProgram(status, p, runner, root_dir)
	runner.Kill()
}

var dump bool
var replace bool
var cpuprofile string
var memprofile string
var jobs int
var verbosity int

func main() {
	flag.BoolVar(&dump, "dump", false, "Dump flowgraphs to disk.")
	flag.BoolVar(&replace, "replace", false, "Replace the existing implementation.")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to file")
	flag.StringVar(&memprofile, "memprofile", "", "write memory profile to this file")
	flag.IntVar(&jobs, "j", runtime.NumCPU(), "Number of threads.")
	flag.IntVar(&verbosity, "v", 0, "Verbosity level.")

	flag.Parse()

	runtime.GOMAXPROCS(jobs)
	framework.Verbosity = verbosity

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

	entryPoint(p, status.Pass("regenerate"))

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
