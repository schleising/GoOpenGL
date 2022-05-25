package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/schleising/GoOpenGL/screen"
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

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	var screen screen.Screen
	screen.SetScreenSize(width, height)

	var rect shapes.Rectangle
	rect.Create(200, 150, 400, 300, screen)

	vao := rect.Handle

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

	transMat := mgl32.Translate3D(float32(count)/100.0, float32(count)/100.0, 0.0)
	rotMat := mgl32.HomogRotate3DZ(mgl32.DegToRad(float32(count)))

	transMat = transMat.Mul4(rotMat)

	translation := gl.GetUniformLocation(program, gl.Str("translation"+"\x00"))
	gl.UniformMatrix4fv(translation, 1, false, &transMat[0])

	gl.DrawElements(gl.TRIANGLES, pointsPerTriangle*numTriangles, gl.UNSIGNED_INT, gl.Ptr(nil))

	gl.BindVertexArray(0)

	window.SwapBuffers()
	glfw.PollEvents()
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
