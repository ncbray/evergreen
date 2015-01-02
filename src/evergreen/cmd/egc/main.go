// Tool for compiling dub sources.
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
	"strings"
	"time"
)

func dumpProgram(status compiler.PassStatus, runner *compiler.TaskRunner, program *flow.DubProgram, outputDir []string) {
	status.Begin()
	defer status.End()

	for _, dubPkg := range program.Packages {
		for _, f := range dubPkg.Funcs {
			// Don't output empty functions.
			if f.CFG.NumNodes() <= 2 {
				continue
			}
			styler := &flow.DotStyler{Decl: f, Core: program.Core}
			dot := graph.GraphToDot(f.CFG, styler)
			parts := append(outputDir, "dub_frontend")
			parts = append(parts, dubPkg.Path...)
			parts = append(parts, fmt.Sprintf("%s.svg", f.Name))
			outfile := filepath.Join(parts...)

			runner.Run(func() {
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

func goFlowFuncName(cf *gocore.Function, f *goflow.FlowFunc) string {
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
	return base + "_" + cf.Name
}

func dumpFlowFuncs(status compiler.PassStatus, runner *compiler.TaskRunner, goFlowProgram *goflow.FlowProgram, goCoreProg *gocore.CoreProgram, outputDir []string) {
	status.Begin()
	defer status.End()

	iter := goFlowProgram.FlowFunc_Scope.Iter()
	for iter.Next() {
		fIndex, f := iter.Value()
		// Don't output empty functions.
		if f.CFG.NumNodes() <= 2 {
			continue
		}
		cf := goCoreProg.Function_Scope.Get(gocore.Function_Ref(fIndex))
		dot := graph.GraphToDot(f.CFG, &goflow.DotStyler{Func: f, Core: goCoreProg})
		parts := append(outputDir, "dub_to_go")
		p := goCoreProg.Package_Scope.Get(cf.Package)
		parts = append(parts, p.Path...)
		parts = append(parts, fmt.Sprintf("%s.svg", goFlowFuncName(cf, f)))
		outfile := filepath.Join(parts...)

		runner.Run(func() {
			io.WriteDot(dot, outfile)
		})
	}
}

func processProgram(status compiler.PassStatus, p compiler.LocationProvider, runner *compiler.TaskRunner, config *EGCConfig) {
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

	goFlowProg, goCoreProg, bypass := golang.GenerateGo(status.Pass("dub_to_go"), flowProgram, coreProg, config.RootPackage, config.GenerateTests)
	if config.Dump {
		dumpFlowFuncs(status.Pass("dump_go"), runner, goFlowProg, goCoreProg, config.DumpDir)
	}
	goTreeProg := gotransform.FlowToTree(status.Pass("flow_to_tree"), goFlowProg, goCoreProg, bypass)

	gotree.GoProgramBackend(status.Pass("go_backend"), goTreeProg, goCoreProg, config.OutputDir, runner)
}

func entryPoint(p compiler.LocationProvider, status compiler.PassStatus, config *EGCConfig) {
	status.Begin()
	defer status.End()

	runner := compiler.CreateTaskRunner(config.Jobs)

	processProgram(status, p, runner, config)
	runner.Kill()
}

func mainLoop(config *EGCConfig, profiling bool) {
	p := compiler.MakeProvider()
	status := compiler.MakeStatus(p)

	start := time.Now()
	for i := 0; ; i++ {
		entryPoint(p, status.Pass("egc"), config)
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

type EGCConfig struct {
	Dump          bool
	InputDir      string
	OutputDir     string
	RootPackage   []string
	DumpDir       []string
	GenerateTests bool
	Jobs          int
}

func flagError(message string) {
	fmt.Println(message)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	config := &EGCConfig{
		DumpDir: []string{"output"},
	}

	var rootPackage string
	var verbosity int
	var cpuprofile string
	var memprofile string

	flag.StringVar(&config.InputDir, "indir", "", "Directory containing input files.")
	flag.StringVar(&config.OutputDir, "outdir", "", "Directory to output generated files.")
	flag.StringVar(&rootPackage, "gopackage", "", "Root package for generate sources.")

	flag.BoolVar(&config.Dump, "dump", false, "Dump flowgraphs to output/... (requires graphviz).")
	flag.BoolVar(&config.GenerateTests, "gentests", false, "Generate dub tests.")

	flag.StringVar(&cpuprofile, "cpuprofile", "", "Write cpu profile to file.")
	flag.StringVar(&memprofile, "memprofile", "", "Write memory profile to this file.")
	flag.IntVar(&config.Jobs, "j", runtime.NumCPU(), "Number of threads.")
	flag.IntVar(&verbosity, "v", 0, "Verbosity level.")

	flag.Parse()

	if config.InputDir == "" {
		flagError("-indir is required")
	}
	if config.OutputDir == "" {
		flagError("-outdir is required")
	}
	if rootPackage == "" {
		flagError("-gopackage is required")
	}
	config.RootPackage = strings.Split(rootPackage, "/")

	runtime.GOMAXPROCS(config.Jobs)
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
