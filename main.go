package main

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/schleising/GoOpenGL/shaders"
	"github.com/schleising/GoOpenGL/shapes"
)

const (
	startWidth  = 800
	startHeight = 600

	thumbnailsPerRow = 8

	vertexShaderFile   = "shaders/vertexShader.glsl"
	fragmentShaderFile = "shaders/fragmentShader.glsl"

	baseFolder = "Pictures/Test Images"
)

var (
	rectList     []*shapes.Rectangle
	program      uint32
	xLastDrag    float64
	yLastDrag    float64
	activeRect   *shapes.Rectangle
)

func main() {
	runtime.LockOSThread()

	homeFolder, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	currentFolder := filepath.Join(homeFolder, baseFolder)

	window := initGlfw()
	defer glfw.Terminate()

	width, _ := window.GetSize()

	program = initOpenGL()

	generateThumbnails(currentFolder, width)

	var count uint = 0

	draw(window)

	for !window.ShouldClose() {
		time.Sleep(16 * time.Millisecond)
		count++
		glfw.PollEvents()

	}
}

func generateThumbnails(path string, windowWidth int) {
	matches, err := filepath.Glob(filepath.Join(path, "*.jpg"))

	if err != nil {
		panic(err)
	}

	thumbnailSize := float64(windowWidth) / float64(thumbnailsPerRow)

	for count, _ := range matches {
		xPos := thumbnailSize * float64(count%thumbnailsPerRow)
		yPos := thumbnailSize * float64(count/thumbnailsPerRow)

		rect := shapes.NewRectangle(float32(xPos), float32(yPos), int(thumbnailSize), int(thumbnailSize))
		// rect.SetTexture(file)
		rectList = append(rectList, rect)
	}
}

func draw(window *glfw.Window) {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	width, height := window.GetSize()

	projMat := mgl32.Ortho2D(0, float32(width), float32(height), 0)
	projection := gl.GetUniformLocation(program, gl.Str("projection"+"\x00"))
	gl.UniformMatrix4fv(projection, 1, false, &projMat[0])

	for _, rect := range rectList {
		rect.Draw(program)
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

	window, err := glfw.CreateWindow(startWidth, startHeight, "Hello", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	glfw.SwapInterval(1)
	window.SetKeyCallback(keyCallBack)
	window.SetScrollCallback(scrollCallback)
	window.SetSizeCallback(sizeCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetCursorPosCallback(cursorPosCallback)

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
	}
}

func scrollCallback(window *glfw.Window, xoffset, yoffset float64) {
	for _, rect := range rectList {
		rect.UpdateYPos(float32(yoffset * 10))
	}
	draw(window)
}

func sizeCallback(window *glfw.Window, width int, height int) {
	xScale := float64(width) / startWidth
	// yScale := float64(height) / startHeight
	scale := xScale

	for _, rect := range rectList {
		rect.Scale(float32(scale))
	}

	draw(window)
}

func mouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	x, y := window.GetCursorPos()

	if action == glfw.Press && button == glfw.MouseButtonLeft {
		for _, rect := range rectList {
			if rect.ClickInRect(float32(x), float32(y)) {
				activeRect = rect
				xLastDrag = x
				yLastDrag = y
			}
		}
	} else if action == glfw.Release {
		activeRect = nil
		xLastDrag = 0
		yLastDrag = 0
	}
}

func cursorPosCallback(window *glfw.Window, xpos float64, ypos float64) {
	if activeRect != nil {
		dx := xpos - xLastDrag
		dy := ypos - yLastDrag
		xLastDrag = xpos
		yLastDrag = ypos
		activeRect.UpdateXPos(float32(dx))
		activeRect.UpdateYPos(float32(dy))
		draw(window)
	}
}
