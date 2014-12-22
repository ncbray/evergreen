// Tool for regenerating checked-in sources.
package main

import (
	"evergreen/compiler"
	"evergreen/dub/flow"
	"evergreen/dub/transform"
	"evergreen/dub/transform/golang"
	"evergreen/dub/tree"
	gocore "evergreen/go/core"
	goflow "evergreen/go/flow"
	gotransform "evergreen/go/transform"
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

func dumpProgram(status compiler.PassStatus, runner *compiler.TaskRunner, program *flow.DubProgram, outputDir []string) {
	status.Begin()
	defer status.End()

	for _, dubPkg := range program.Packages {
		for _, f := range dubPkg.Funcs {
			styler := &flow.DotStyler{Decl: f, Core: program.Core}
			dot := graph.GraphToDot(f.CFG, styler)
			parts := append(outputDir, "dub_frontend")
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

func goFlowFuncName(f *goflow.LLFunc) string {
	base := "func"
	if f.Recv != goflow.NoRegister {
		ref := f.Register_Scope.Get(f.Recv)
		at := ref.T
		pt, ok := at.(*gocore.PointerType)
		if !ok {
			panic(at)
		}
		st, ok := pt.Element.(*gocore.StructType)
		if !ok {
			panic(pt.Element)
		}
		base = "meth_" + st.Name
	}
	return base + "_" + f.Name
}

func dumpFlowFuncs(flowFuncs []*goflow.LLFunc, outputDir []string) {
	for _, f := range flowFuncs {
		dot := graph.GraphToDot(f.CFG, &goflow.DotStyler{Ops: f.Ops})
		parts := append(outputDir, "dub_to_go")
		parts = append(parts, f.Package.Path...)
		parts = append(parts, fmt.Sprintf("%s.svg", goFlowFuncName(f)))
		outfile := filepath.Join(parts...)
		io.WriteDot(dot, outfile)
	}
}

func processProgram(status compiler.PassStatus, p compiler.LocationProvider, runner *compiler.TaskRunner, config *RegenerateConfig) {
	program, coreProg := tree.DubProgramFrontend(status.Pass("dub_frontend"), p, config.InputDir)
	if status.ShouldHalt() {
		return
	}
	flowProgram := transform.LowerProgram(status.Pass("lower"), program, coreProg)

	flow.TrimFlow(status.Pass("trim_flow"), flowProgram)

	if config.Dump {
		dumpProgram(status.Pass("dump"), runner, flowProgram, config.DumpDir)
	}

	analyizeProgram(flowProgram)

	flowProg, bypass := golang.GenerateGo(status.Pass("dub_to_go"), flowProgram, coreProg, config.RootPackage, config.GenerateTests)
	if config.Dump {
		dumpFlowFuncs(flowProg.Functions, config.DumpDir)
	}
	prog := gotransform.FlowToTree(status.Pass("flow_to_tree"), flowProg, bypass)

	gotree.GoProgramBackend(status.Pass("go_backend"), prog, config.OutputDir, runner)
}

func entryPoint(p compiler.LocationProvider, status compiler.PassStatus, config *RegenerateConfig) {
	status.Begin()
	defer status.End()

	runner := compiler.CreateTaskRunner(config.Jobs)

	processProgram(status, p, runner, config)
	runner.Kill()
}

func mainLoop(config *RegenerateConfig, profiling bool) {
	p := compiler.MakeProvider()
	status := compiler.MakeStatus(p)

	start := time.Now()
	for i := 0; ; i++ {
		entryPoint(p, status.Pass("regenerate"), config)
		if profiling && time.Since(start) < time.Second*10 {
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

type RegenerateConfig struct {
	Dump          bool
	InputDir      string
	OutputDir     string
	RootPackage   string
	DumpDir       []string
	GenerateTests bool
	Jobs          int
}

func main() {
	config := &RegenerateConfig{
		InputDir:    "dub",
		OutputDir:   "src",
		DumpDir:     []string{"output"},
		RootPackage: "generated",
	}

	var replace bool
	var verbosity int
	var cpuprofile string
	var memprofile string

	flag.BoolVar(&config.Dump, "dump", false, "Dump flowgraphs to disk.")
	flag.BoolVar(&replace, "replace", false, "Replace the existing implementation.")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to file")
	flag.StringVar(&memprofile, "memprofile", "", "write memory profile to this file")
	flag.IntVar(&config.Jobs, "j", runtime.NumCPU(), "Number of threads.")
	flag.IntVar(&verbosity, "v", 0, "Verbosity level.")

	flag.Parse()

	runtime.GOMAXPROCS(config.Jobs)
	compiler.Verbosity = verbosity

	if replace {
		config.RootPackage = "evergreen"
	}
	config.GenerateTests = !replace

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

	mainLoop(config, cpuprofile != "")

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
