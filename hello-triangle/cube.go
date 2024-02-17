// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 2.1.
package main // import "github.com/go-gl/example/gl21-cube"

import (
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

var triangleVerts = []float32{
	-0.5, -0.5, 0.0,
	0.5, -0.5, 0.0,
	0.0, 0.5, 0.0,
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
	var sizeOfFloat64 int = int(unsafe.Sizeof(float64(0)))
	var sizeOfFloat32 int = int(unsafe.Sizeof(float32(0)))

	gl.BufferData(gl.ARRAY_BUFFER, len(triangleVerts)*sizeOfFloat64, ptr, gl.STATIC_DRAW)

	// Vertext shader setup
	// ================================
	vertexShaderSource := `
		#version 330 core
		layout (location = 0) in vec3 aPos;
	
		void main()
		{
			gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1);
		}
	`
	glVertexSourceInt, freeFn := gl.Strs(vertexShaderSource + "\x00")

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)

	// How to convert this in go....
	gl.ShaderSource(vertexShader, 1, glVertexSourceInt, nil)
	defer freeFn()
	gl.CompileShader(vertexShader)

	// Check for compile error
	var success int32 = -1
	var infoLog *uint8
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &success)

	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(vertexShader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(vertexShader, logLength, nil, gl.Str(log))

		fmt.Printf("Failed to compile vertex shader: %v", log)
		return
	}

	// Fragment shader setup
	// ================================
	fragmentShaderSource := `
		#version 330 core
		out vec4 FragColor;
		
		void main()
		{
			FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
		} 
	`

	glFragSourceInt, freeFragFn := gl.Strs(fragmentShaderSource + "\x00")
	defer freeFragFn()

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)

	gl.ShaderSource(fragmentShader, 1, glFragSourceInt, nil)
	gl.CompileShader(fragmentShader)

	// Check for compile error
	success = -1
	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &success)

	if success != 1 {
		gl.GetShaderInfoLog(vertexShader, 512, nil, infoLog)
		fmt.Println("Error in fragment shader\n", &infoLog)
	}

	// Setup shader program
	// ============================
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	success = -1
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)

	if success != 1 {
		gl.GetProgramInfoLog(shaderProgram, 512, nil, infoLog)
		fmt.Println("Error in shader program: \n", &infoLog)
	}

	gl.UseProgram(shaderProgram)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(3*sizeOfFloat32), nil)
	gl.EnableVertexAttribArray(0)

	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}
