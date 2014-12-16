// Tool for regenerating checked-in sources.
package main

import (
	"evergreen/compiler"
	"evergreen/dub/core"
	"evergreen/dub/flow"
	"evergreen/dub/transform"
	"evergreen/dub/transform/golang"
	"evergreen/dub/tree"
	gotree "evergreen/go/tree"
	"evergreen/graph"
	"evergreen/io"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"
)

func dumpProgram(status compiler.PassStatus, runner *compiler.TaskRunner, program *flow.DubProgram) {
	status.Begin()
	defer status.End()

	for _, dubPkg := range program.Packages {
		for _, f := range dubPkg.Funcs {
			styler := &flow.DotStyler{Decl: f, Core: program.Core}
			dot := graph.GraphToDot(f.CFG, styler)
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

func analyizeProgram(program *flow.DubProgram) {
	for _, dubPkg := range program.Packages {
		for _, s := range dubPkg.Structs {
			if s.Implements != nil {
				s.Implements.IsParent = true
			}
		}
	}
}

func GenerateGo(status compiler.PassStatus, program *flow.DubProgram, coreProg *core.CoreProgram, runner *compiler.TaskRunner) {
	status.Begin()
	defer status.End()

	root := "generated"
	if replace {
		root = "evergreen"
	}
	prog := golang.GenerateGo(status.Pass("generate_go"), program, coreProg, root, !replace, dump)

	// Compact simple expressions back into tree form.
	gotree.Consolidate(status.Pass("consolidate"), prog)

	// Give everything names: variables, etc.
	gotree.Nameify(status.Pass("nameify"), prog)

	// Generate the sources.
	gotree.OutputProgram(status.Pass("output"), prog, "src", runner)
}

func processProgram(status compiler.PassStatus, p compiler.LocationProvider, runner *compiler.TaskRunner, root string) {
	program, coreProg := tree.DubProgramFrontend(status.Pass("dub_frontend"), p, root)
	if status.ShouldHalt() {
		return
	}
	flowProgram := transform.LowerProgram(status.Pass("lower"), program, coreProg)

	flow.TrimFlow(status.Pass("trim_flow"), flowProgram)

	if dump {
		dumpProgram(status.Pass("dump"), runner, flowProgram)
	}

	analyizeProgram(flowProgram)
	GenerateGo(status.Pass("go_backend"), flowProgram, coreProg, runner)
}

func entryPoint(p compiler.LocationProvider, status compiler.PassStatus) {
	status.Begin()
	defer status.End()

	runner := compiler.CreateTaskRunner(jobs)

	root_dir := "dub"
	processProgram(status, p, runner, root_dir)
	runner.Kill()
}

func mainLoop() {
	p := compiler.MakeProvider()
	status := compiler.MakeStatus(p)

	start := time.Now()
	for i := 0; ; i++ {
		entryPoint(p, status.Pass("regenerate"))
		if cpuprofile != "" && time.Since(start) < time.Second*10 {
			fmt.Println("Re-running to improve profiling data", i)
		} else {
			break
		}
	}

	if status.ShouldHalt() {
		fmt.Printf("%d errors\n", status.ErrorCount())
		os.Exit(1)
	}
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
	compiler.Verbosity = verbosity

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			fmt.Println(err.Error())
			return
		} else {
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
	}

	mainLoop()

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			fmt.Println(err.Error())
			return
		} else {
			pprof.WriteHeapProfile(f)
			f.Close()
		}
	}
}
