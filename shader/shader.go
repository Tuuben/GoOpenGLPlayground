package shader

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

type Shader struct {
  ProgramId uint32
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

func Create(vertexPath string, fragmentPath string) Shader { 
  // Load vertex file
  vertexFileContent, err := os.ReadFile(vertexPath)

  if err != nil {
    panic("Failed to load vertex shader from file.")
  }

  vertexSource := string(vertexFileContent)

  // Load fragment file
  fragmentFileContent, err := os.ReadFile(fragmentPath)

  if err != nil {
    panic("Failed to load fragment from file.")
  }

  fragmentSource := string(fragmentFileContent)


	// Vertext shader setup
	// ================================
	glVertexSourceInt, freeVertexFn := gl.Strs(vertexSource + "\x00")
	defer freeVertexFn()
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)

	gl.ShaderSource(vertexShader, 1, glVertexSourceInt, nil)
	gl.CompileShader(vertexShader)
  defer gl.DeleteShader(vertexShader)

  err = checkShaderCompileStatus(vertexShader)

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
  defer gl.DeleteShader(fragmentShader)

  err = checkShaderCompileStatus(fragmentShader)

  if err != nil {
    fmt.Println("Failed to compile Fragment Shader")
  }

	// Setup shader program
	// ============================
	programId := gl.CreateProgram()
	gl.AttachShader(programId, vertexShader)
	gl.AttachShader(programId, fragmentShader)
	gl.LinkProgram(programId)

  err = checkShaderCompileStatus(programId)

  if err != nil {
    fmt.Println("Failed to compile Shader Program")
  }

  return Shader{
    ProgramId: programId,
  } 
}

func (shader Shader) Use() {
  gl.UseProgram(shader.ProgramId)
}

func (shader Shader) SetUniformBool(name string, value bool) {
  nameCStr := gl.Str(name + "\x00")
  boolToIntMap := map[bool]int32{ true: 1, false: 0 }
  gl.Uniform1i(gl.GetUniformLocation(shader.ProgramId, nameCStr), boolToIntMap[value])
}

func (shader Shader) SetUniformInt(name string, value int32) {
  nameCStr := gl.Str(name + "\x00")
  gl.Uniform1i(gl.GetUniformLocation(shader.ProgramId, nameCStr), value)
}

func (shader Shader) SetUniformFloat(name string, value float32) {
  nameCStr := gl.Str(name + "\x00")
  gl.Uniform1f(gl.GetUniformLocation(shader.ProgramId, nameCStr), value)
}

func (shader Shader) SetUniformVec4(name string, v0 float32, v1 float32, v2 float32, v3 float32) {
  nameCStr := gl.Str(name + "\x00")
  gl.Uniform4f(gl.GetUniformLocation(shader.ProgramId, nameCStr), v0, v1, v2, v3)
}

