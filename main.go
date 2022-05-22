package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 800
	height = 600
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	initOpenGL()

	for !window.ShouldClose() {
		draw(window)
	}
}

func draw(window *glfw.Window) {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	window.SwapBuffers()
	glfw.PollEvents()
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

func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic((err))
	}
	gl.ClearColor(0.118, 0.565, 1.0, 1.0)
}

func keyCallBack(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, modifier glfw.ModifierKey) {
	if action == glfw.Press {
		fmt.Println(key)
		if key == glfw.KeyEscape {
			window.SetShouldClose(true)
		}
	}
}
