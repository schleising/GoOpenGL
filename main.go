package main

import (
	"fmt"
	"runtime"

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
	vertices = []float32{
		// positions   // colors // texture coords
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom left
		-0.5, 0.5, 0.0, 1.0, 1.0, 1.0, 0.0, 1.0, // top left
	}

	indices = []uint32{
		0, 1, 2,
		0, 2, 3,
	}
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	vao := makeVao(vertices, indices)

	texture, err := textures.LoadImage("images/pipeline.png")

	if err != nil {
		panic(0)
	}

	for !window.ShouldClose() {
		draw(window, program, vao, texture)
	}
}

func draw(window *glfw.Window, program uint32, vao uint32, texture uint32) {
	// timeValue := glfw.GetTime()
	// greenValue := (math.Sin(timeValue) / 2.0) + 0.5
	// vertexColourLocation := gl.GetUniformLocation(program, gl.Str("ourColour\x00"))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.UseProgram(program)
	// gl.Uniform4f(vertexColourLocation, 0.118, float32(greenValue), 1.0, 1.0)
	gl.BindVertexArray(vao)
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

	glfw.WindowHint(glfw.Resizable, glfw.False)
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
