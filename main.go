package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 800
	height = 600

	vertexSize    = 3
	sizeOfFloat32 = 4
	sizeOfUint32  = 4
	numAttributes = 2

	vertexShaderSource = `
	#version 410 core
    layout (location = 0) in vec3 aPos;
    layout (location = 1) in vec3 aColour;

	out vec3 ourColour;

	void main()
    {
       gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
	   ourColour = aColour;
    }
` + "\x00"

	fragmentShaderSource = `
	#version 410 core
	in vec3 ourColour;
	out vec4 FragColor;
	void main()
	{
		FragColor = vec4(ourColour, 1.0);
	} 
` + "\x00"
)

var (
	vertices = []float32{
		// positions      // colors
		0.5, -0.5, 0.0, 1.0, 0.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // bottom left
		0.0, 0.5, 0.0, 0.0, 0.0, 1.0, // top
	}
)

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()

	vao := makeVao(vertices)

	for !window.ShouldClose() {
		draw(window, program, vao)
	}
}

func draw(window *glfw.Window, program uint32, vao uint32) {
	// timeValue := glfw.GetTime()
	// greenValue := (math.Sin(timeValue) / 2.0) + 0.5
	// vertexColourLocation := gl.GetUniformLocation(program, gl.Str("ourColour\x00"))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.UseProgram(program)
	// gl.Uniform4f(vertexColourLocation, 0.118, float32(greenValue), 1.0, 1.0)
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, vertexSize)
	gl.BindVertexArray(0)

	window.SwapBuffers()
	glfw.PollEvents()
}

func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*sizeOfFloat32, gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var offset uintptr = 0
	gl.VertexAttribPointerWithOffset(0, vertexSize, gl.FLOAT, false, vertexSize*sizeOfFloat32*numAttributes, offset)
	gl.EnableVertexAttribArray(0)

	offset = vertexSize * sizeOfFloat32
	gl.VertexAttribPointerWithOffset(1, vertexSize, gl.FLOAT, false, vertexSize*sizeOfFloat32*numAttributes, offset)
	gl.EnableVertexAttribArray(1)

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

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
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

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
