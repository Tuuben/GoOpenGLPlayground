package main

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/example/hello-triangle/shader"
	"github.com/go-gl/example/utils"
	"github.com/go-gl/example/window"
	"github.com/go-gl/gl/v4.1-core/gl"
)

var vertices = []float32 {
  -0.5, -0.5, 0.0,
  0.5, -0.5, 0.0,
  0.5, 0.5, 0.0,
}

func main() {

  window.Create(640, 480, "Textures window", func() {
	  gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
  
    // Setup Vertext Array Object
    var VAO uint32
    gl.GenVertexArrays(1, &VAO)
    gl.BindVertexArray(VAO)

    // Setup Vertex Buffer Object
    var VBO uint32
    gl.GenBuffers(1, &VBO)
    gl.BindBuffer(gl.ARRAY_BUFFER, VBO)

    // Main triangle 
    firstVerticiePtr := unsafe.Pointer(&vertices[0])
    gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*utils.SizeOfFloat32, firstVerticiePtr, gl.STATIC_DRAW)

    // Create shader
    shader := shader.Create("./shaders/vertexShader.glsl", "./shaders/fragShader.glsl")
    shader.Use()

    gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(3*utils.SizeOfFloat32), nil)
    gl.EnableVertexAttribArray(0)

    gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
  
    gl.BindVertexArray(0)
    fmt.Println("Hello world")
  })
}
