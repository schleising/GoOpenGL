package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/schleising/GoOpenGL/shaders"
	"github.com/schleising/GoOpenGL/shapes"
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
	rectList []*shapes.Rectangle
	program  uint32
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program = initOpenGL()

	rect1 := shapes.NewRectangle(500, 100, 200, 150)

	texture1, err := textures.LoadImage("images/IMG_0033.JPG")

	if err != nil {
		panic(0)
	}

	rect1.SetTexture(texture1)
	rect1.SetProgram(program)

	rectList = append(rectList, rect1)

	rect2 := shapes.NewRectangle(50, 200, 400, 300)

	texture2, err := textures.LoadImage("images/pipeline.png")

	if err != nil {
		panic(0)
	}

	rect2.SetTexture(texture2)
	rect2.SetProgram(program)

	rectList = append(rectList, rect2)

	var count uint = 0

	draw(window)

	for !window.ShouldClose() {
		time.Sleep(16 * time.Millisecond)
		count++
		glfw.PollEvents()

	}
}

func draw(window *glfw.Window) {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	width, height := window.GetSize()

	projMat := mgl32.Ortho2D(0, float32(width), float32(height), 0)
	projection := gl.GetUniformLocation(program, gl.Str("projection"+"\x00"))
	gl.UniformMatrix4fv(projection, 1, false, &projMat[0])

	for _, rect := range rectList {
		rect.Draw()
	}

	window.SwapBuffers()
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
	window.SetScrollCallback(scrollCallback)
	window.SetSizeCallback(sizeCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)

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

	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	return prog
}

func keyCallBack(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, modifier glfw.ModifierKey) {
	if action == glfw.Press {
		if key == glfw.KeyEscape {
			window.SetShouldClose(true)
		}
		if key == glfw.KeyUp {
			for _, rect := range rectList {
				rect.YPos += 1
			}
		}
	}
}

func scrollCallback(window *glfw.Window, xoffset, yoffset float64) {
	for _, rect := range rectList {
		rect.YPos = rect.YPos + float32(yoffset*10)
	}
	draw(window)
}

func sizeCallback(window *glfw.Window, width int, height int) {
	draw(window)
}

func mouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	x, y := window.GetCursorPos()

	if action == glfw.Press && button == glfw.MouseButtonLeft {
		for _, rect := range rectList {
			if rect.ClickInRect(float32(x), float32(y)) {
				fmt.Println("Click In Rect")
			}
		}
	}
}
