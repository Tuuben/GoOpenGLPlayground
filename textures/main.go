package main

import (
	"fmt"

	"github.com/go-gl/example/window"
)


func main() {
  window.Create(640, 480, "Textures window", func() {
    fmt.Println("Hello world")
  })
}
