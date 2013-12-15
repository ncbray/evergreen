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
		&dub.BinaryOp{
			Left:  0,
			Op:    "<",
			Right: 1,
			Dst:   2,
		},
	})
	decide := dub.CreateSwitch(2)
	body := dub.CreateBlock([]dub.DubOp{
		&dub.ConstantIntOp{Value: 1, Dst: 3},
		&dub.BinaryOp{
			Left:  0,
			Op:    "+",
			Right: 3,
			Dst:   0,
		},
	})

	l.Connect(0, cond)
	l.AttachDefaultExits(cond)

	l.Connect(0, decide)
	decide.SetExit(0, body)

	l.AttachDefaultExits(body)
	l.Connect(0, cond)
	decide.SetExit(1, l.GetExit(0))

	i := "integer"
	b := "boolean"

	registers := []dub.RegisterInfo{
		dub.RegisterInfo{T: i},
		dub.RegisterInfo{T: i},
		dub.RegisterInfo{T: b},
		dub.RegisterInfo{T: i},
	}

	dot := base.RegionToDot(l)
	outfile := filepath.Join("output", "test.svg")

	result := make(chan error, 2)
	go func() {
		err := io.WriteDot(dot, outfile)
		result <- err
	}()

	fmt.Println(dub.GenerateGo(l, registers))

	err := <-result
	if err != nil {
		fmt.Println(err)
	}

}
