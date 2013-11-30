package main

import (
	"evergreen/base"
	"evergreen/dub"
	"evergreen/io"
	"fmt"
	"path/filepath"
)

func main() {
	l := dub.CreateRegion()
	cond := dub.CreateBlock([]dub.DubOp{
		&dub.GetLocalOp{Name: "counter", Dst: 1},
		&dub.GetLocalOp{Name: "limit", Dst: 2},
		&dub.BinaryOp{
			Left:  1,
			Op:    "<",
			Right: 2,
			Dst:   3,
		},
	})
	decide := dub.CreateSwitch(3)
	body := dub.CreateBlock([]dub.DubOp{
		&dub.GetLocalOp{Name: "counter", Dst: 4},
		&dub.ConstantIntOp{Value: 1, Dst: 5},
		&dub.BinaryOp{
			Left:  4,
			Op:    "+",
			Right: 5,
			Dst:   6,
		},
		&dub.SetLocalOp{Src: 6, Name: "counter"},
	})

	l.Connect(0, cond)
	l.AttachDefaultExits(cond)

	l.Connect(0, decide)
	decide.SetExit(0, body)

	l.AttachDefaultExits(body)
	l.Connect(0, cond)
	decide.SetExit(1, l.GetExit(0))

	dot := base.RegionToDot(l)
	outfile := filepath.Join("output", "test.svg")

	result := make(chan error, 2)
	go func() {
		err := io.WriteDot(dot, outfile)
		result <- err
	}()

	fmt.Println(dub.GenerateGo())

	err := <-result
	if err != nil {
		fmt.Println(err)
	}

}
