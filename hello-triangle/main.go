// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 2.1.
package main // import "github.com/go-gl/example/gl21-cube"

import (
	_ "image/png"
	"log"
	"math"
	"runtime"
	"unsafe"

	"github.com/go-gl/example/hello-triangle/shader"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const width, height = 640, 480

var roofVerts = []float32{
  // positions      //colors
  -0.75, 0.0, 0.0,  1.0, 0.0, 0.0,
  0.0, 1.0, 0.0,    0.0, 1.0, 0.0,
  0.75, 0.0, 0.0,   0.0, 0.0, 1.0,
}

var triangleVerts = []float32{
  0.5,  0.0, 0.0,
  0.5, -1.0, 0.0, 
  -0.5, -1.0, 0.0,
  -0.5,  0.0, 0.0, 
}

var indicies = []uint32 {
  0, 1, 3,
  1, 2, 3,
}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
  glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Hello Triangwleh", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	for !window.ShouldClose() {
		setupScene()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func setupScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	//var sizeOfFloat64 int = int(unsafe.Sizeof(float64(0)))
	var sizeOfFloat32 int = int(unsafe.Sizeof(float32(0)))
  var sizeOfInt int =  int(unsafe.Sizeof(uint32(0)))

	// VAO Vertext Array Object
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	// VBO BUFFER, stores a bunch of verticies at once, sending verts to the GPU from
	// CPU is kinda slow so we want to store as much as possible at once.
	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)

  // Main house triangle verts being set
	ptr := unsafe.Pointer(&triangleVerts[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(triangleVerts)*sizeOfFloat32, ptr, gl.STATIC_DRAW)

  // EBO Element Buffer Object
  // ==========
  var EBO uint32 
  gl.GenBuffers(1, &EBO)
  gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)

  gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indicies) * sizeOfInt, unsafe.Pointer(&indicies[0]), gl.STATIC_DRAW)

  //shaderProgram := createShaderProgram(vertexShaderSource, fragmentShaderSource)
  houseShader := shader.Create("./shader/vertexShader.glsl", "./shader/fragShader.glsl")
  houseShader.Use()
	//gl.UseProgram(shaderProgram)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(3*sizeOfFloat32), nil)
	gl.EnableVertexAttribArray(0)

  gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

  // Roof
  roofShader := shader.Create("./shader/vertRoof.glsl", "./shader/fragRoof.glsl")

  timeValue := glfw.GetTime()
  greenValue := (math.Sin(timeValue) / 2.0) + 0.5

  roofShader.Use()
  roofShader.SetUniformVec4("ourColor", 0.0, float32(greenValue), 0.0, 1.0)

	roofPtr := unsafe.Pointer(&roofVerts[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(roofVerts)*sizeOfFloat32, roofPtr, gl.STATIC_DRAW)

  // Roof verts attrib pointer
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(6*sizeOfFloat32), nil)
	gl.EnableVertexAttribArray(0)

  // Roof color atrib pointer
  gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(6*sizeOfFloat32), unsafe.Pointer(uintptr(3 * sizeOfFloat32)))
  gl.EnableVertexAttribArray(1)

  gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

  gl.BindVertexArray(0)
}
