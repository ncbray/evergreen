package main

import (
  "fmt"
  "muir/dub"
)

type ExactChar struct {
	exact rune
}

func main() {
  fmt.Println(dub.GenerateGo())
}
