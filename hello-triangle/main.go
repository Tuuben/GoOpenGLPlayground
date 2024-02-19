// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 2.1.
package main // import "github.com/go-gl/example/gl21-cube"

import (
	"errors"
	"fmt"
	_ "image/png"
	"log"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const width, height = 640, 480

var roofVerts = []float32{
  -0.75, 0.0, 0.0,
  0.0, 1.0, 0.0, 
  0.75, 0.0, 0.0,
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

var vertexShaderSource = `
		#version 330 core
		layout (location = 0) in vec3 aPos;
	
		void main()
		{
			gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1);
		}
	`

var fragmentShaderSource = `
		#version 330 core
		out vec4 FragColor;
		
		void main()
		{
			FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
		} 
	`

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

func createShaderProgram(vertexSource string, fragmentSource string) uint32 {
	// Vertext shader setup
	// ================================
	glVertexSourceInt, freeFn := gl.Strs(vertexSource + "\x00")
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)

	// How to convert this in go....
	gl.ShaderSource(vertexShader, 1, glVertexSourceInt, nil)
	defer freeFn()
	gl.CompileShader(vertexShader)

  err := checkShaderCompileStatus(vertexShader)

  if err != nil {
    fmt.Println("Failed to compile Vertex Shader.")
  }

	// Fragment shader setup
	// ================================
	glFragSourceInt, freeFragFn := gl.Strs(fragmentSource+ "\x00")
	defer freeFragFn()

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)

	gl.ShaderSource(fragmentShader, 1, glFragSourceInt, nil)
	gl.CompileShader(fragmentShader)

  err = checkShaderCompileStatus(fragmentShader)

  if err != nil {
    fmt.Println("Failed to compile Fragment Shader")
  }

	// Setup shader program
	// ============================
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

  err = checkShaderCompileStatus(shaderProgram)

  if err != nil {
    fmt.Println("Failed to compile Shader Program")
  }

  return shaderProgram
}

func checkShaderCompileStatus(shader uint32) error {
	var success int32 = -1
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)

	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		fmt.Printf("Failed to compile shader: %v", log)
    return errors.New("Failed to compile")
	}

  return nil
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

	ptr := unsafe.Pointer(&triangleVerts[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(triangleVerts)*sizeOfFloat32, ptr, gl.STATIC_DRAW)


  // EBO Element Buffer Object
  // ==========
  var EBO uint32 
  gl.GenBuffers(1, &EBO)

  gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
  gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indicies) * sizeOfInt, unsafe.Pointer(&indicies[0]), gl.STATIC_DRAW)


  shaderProgram := createShaderProgram(vertexShaderSource, fragmentShaderSource)
	gl.UseProgram(shaderProgram)

//	gl.DeleteShader(vertexShader)
// gl.DeleteShader(fragmentShader)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(3*sizeOfFloat32), nil)
	gl.EnableVertexAttribArray(0)

  gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)


  // Roof
	roofPtr := unsafe.Pointer(&roofVerts[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(roofVerts)*sizeOfFloat32, roofPtr, gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(3*sizeOfFloat32), nil)
	gl.EnableVertexAttribArray(0)

  gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

  // Clear out Vertex Object Array (VAO)
  gl.BindVertexArray(0)
}
