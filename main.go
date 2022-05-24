package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/schleising/GoOpenGL/shaders"
	"github.com/schleising/GoOpenGL/textures"
)

const (
	width  = 800
	height = 600

	positionSize      = 3
	colourSize        = 3
	texCoordSize      = 2
	vertexSize        = positionSize + colourSize + texCoordSize
	sizeOfFloat32     = 4
	sizeOfUint32      = 4
	numAttributes     = 3
	pointsPerTriangle = 3
	numTriangles      = 2

	vertexShaderFile   = "shaders/vertexShader.glsl"
	fragmentShaderFile = "shaders/fragmentShader.glsl"
)

var (
	// The four vertices of a rectangle
	vertices = []float32{
		// Top Right Vertex
		0.5, 0.5, 0.0, // Postion
		1.0, 0.0, 0.0, // Colour
		1.0, 0.0, // Texture Coord

		// Bottom Right Vertex
		0.5, -0.5, 0.0, // Postion
		0.0, 1.0, 0.0, // Colour
		1.0, 1.0, // Texture Coord

		// Bottom Left Vertex
		-0.5, -0.5, 0.0, // Postion
		0.0, 0.0, 1.0, // Colour
		0.0, 1.0,  // Texture Coord
		
		// Top Left Vertex
		-0.5, 0.5, 0.0, // Postion
		1.0, 1.0, 1.0, // Colour
		0.0, 0.0, // Texture Coord
	}

	// The indices into vertices to make a rectangle from two triangles
	indices = []uint32{
		0, 1, 2, // TR -> BR -> BL
		0, 2, 3, // TR -> BL -> TL
	}
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	vao := makeVao(vertices, indices)

	texture1, err := textures.LoadImage("images/IMG_0033.JPG")

	if err != nil {
		panic(0)
	}

	texture2, err := textures.LoadImage("images/pipeline.png")

	if err != nil {
		panic(0)
	}

	texture := texture1
	var count uint = 0

	for !window.ShouldClose() {
		if int(count/10)%2 == 0 {
			texture = texture1
		} else {
			texture = texture2
		}
		draw(window, program, vao, texture, count)
		time.Sleep(50 * time.Millisecond)
		count++
	}
}

func draw(window *glfw.Window, program uint32, vao uint32, texture uint32, count uint) {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.UseProgram(program)

	gl.BindVertexArray(vao)

	dx := gl.GetUniformLocation(program, gl.Str("dx"+"\x00"))
	gl.Uniform1f(dx, float32(count)/100.0)

	gl.DrawElements(gl.TRIANGLES, pointsPerTriangle*numTriangles, gl.UNSIGNED_INT, gl.Ptr(nil))

	gl.BindVertexArray(0)

	window.SwapBuffers()
	glfw.PollEvents()
}

func makeVao(vertices []float32, indices []uint32) uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*sizeOfFloat32, gl.Ptr(vertices), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*sizeOfUint32, gl.Ptr(indices), gl.STATIC_DRAW)

	var offset uintptr = 0
	gl.VertexAttribPointerWithOffset(0, positionSize, gl.FLOAT, false, vertexSize*sizeOfFloat32, offset)
	gl.EnableVertexAttribArray(0)

	offset = positionSize * sizeOfFloat32
	gl.VertexAttribPointerWithOffset(1, colourSize, gl.FLOAT, false, vertexSize*sizeOfFloat32, offset)
	gl.EnableVertexAttribArray(1)

	offset = (positionSize + colourSize) * sizeOfFloat32
	gl.VertexAttribPointerWithOffset(2, texCoordSize, gl.FLOAT, false, vertexSize*sizeOfFloat32, offset)
	gl.EnableVertexAttribArray(2)

	return vao
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Hello", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	glfw.SwapInterval(1)
	window.SetKeyCallback(keyCallBack)

	return window
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic((err))
	}

	gl.ClearColor(0.118, 0.565, 1.0, 1.0)

	vertexShader, err := shaders.LoadShader(vertexShaderFile, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := shaders.LoadShader(fragmentShaderFile, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	gl.UseProgram(prog)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	var nAttributes int32
	gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS, &nAttributes)
	fmt.Println(nAttributes)

	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	return prog
}

func keyCallBack(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, modifier glfw.ModifierKey) {
	if action == glfw.Press {
		fmt.Println(key)
		if key == glfw.KeyEscape {
			window.SetShouldClose(true)
		}
	}
}
