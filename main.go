package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const(
	width = 800
	height = 600
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

    for !window.ShouldClose() {
		window.SwapBuffers()
		glfw.PollEvents()
    }
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

func keyCallBack(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, modifier glfw.ModifierKey) {
	if action == glfw.Press {
		fmt.Println(key)
		if key == glfw.KeyEscape {
			window.SetShouldClose(true)
		}
	}
}
