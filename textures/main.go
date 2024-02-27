package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
  _ "image/png"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/example/hello-triangle/shader"
	"github.com/go-gl/example/utils"
	"github.com/go-gl/example/window"
	"github.com/go-gl/gl/v2.1/gl"
)


var vertices = []float32{
  // positions          // colors           // texture coords
  0.5,  0.5, 0.0,   1.0, 0.0, 0.0,    1.0, 1.0, // top right
  0.5, -0.5, 0.0,   0.0, 1.0, 0.0,    1.0, 0.0, // bottom right
  -0.5, -0.5, 0.0,   0.0, 0.0, 1.0,   0.0, 0.0, // bottom left
  -0.5,  0.5, 0.0,   1.0, 1.0, 0.0,   0.0, 1.0,   // top left
}

var indicies = []uint32 {
  0, 1, 3,
  1, 2, 3,
}

var width, height, nrChannels int;
var VAO, VBO, EBO uint32;
var shaderProgram shader.Shader;
var texture uint32;

func main() {
	runtime.LockOSThread()
  window.Create(640, 480, "Textures window", onWindowStart, onWindowUpdate)
}

func onWindowStart() {
  fmt.Println("Start: ")
  // LOAD IMAGE
  // ==============
  imgFile, err := os.Open("./images/container.jpg")
  if err != nil {
    log.Println("Failed to open image.")
  }
  defer imgFile.Close()

  img, _, err := image.Decode(imgFile)
  if err != nil {
    log.Println("Failed to decode image.")
    log.Println(err)
  }

  rgba := image.NewRGBA(img.Bounds())
  if rgba.Stride != rgba.Rect.Size().X*4 {
    log.Println("Failed to use stride.")
  }

  draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

  // Setup GL draw
  // ================
  shaderProgram = shader.Create("./shaders/vertexShader.glsl", "./shaders/fragShader.glsl")

  gl.GenVertexArrays(1, &VAO)
  gl.GenBuffers(1, &VBO)
  gl.GenBuffers(1, &EBO)

  gl.BindVertexArray(VAO)

  // Bind verticies array to VBO 
  gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
  gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*utils.SizeOfFloat32, gl.Ptr(&vertices[0]), gl.STATIC_DRAW)

  // Bind indicies to EBO buffer
  gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER , EBO)
  gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indicies) * utils.SizeOfUint32, gl.Ptr(&indicies[0]), gl.STATIC_DRAW)

  // Position attribute
  gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(8*utils.SizeOfFloat32), nil)
  gl.EnableVertexAttribArray(0) // Enable the array for verts
  // Color attribute
  gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(8*utils.SizeOfFloat32), gl.Ptr(uintptr(3 * utils.SizeOfFloat32)))
  gl.EnableVertexAttribArray(1) // Enable the array for colors 
  // Texture attribute
  gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(8*utils.SizeOfFloat32), gl.Ptr(uintptr(6 * utils.SizeOfFloat32)))
  gl.EnableVertexAttribArray(2) // enable the array for tex coords

  // Render 
  shaderProgram.Use()
  shaderProgram.SetUniformInt("texture1", 0)

  // Load texture
  loadedTexture, err := newTexture("./images/gravel.jpeg");

  if err != nil {
    fmt.Println("Failed to load texture");
  }

  texture = loadedTexture;
}

func onWindowUpdate() {
	  gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.ClearColor(0.2, 0.3, 0.3, 1.0)
  
    gl.ActiveTexture(gl.TEXTURE0)
	  gl.BindTexture(gl.TEXTURE_2D, texture)

    shaderProgram.Use()
    gl.BindVertexArray(VAO)
    gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil) 

    //gl.BindVertexArray(0)
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
  gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
  gl.GenerateMipmap(gl.TEXTURE_2D)

	return texture, nil
}
