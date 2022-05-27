package shaders

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func LoadShader(filename string, shaderType uint32) (uint32, error) {
	src, err := ioutil.ReadFile(filename)

	if err != nil {
		return 0, err
	}

	handle := gl.CreateShader(shaderType)

	glSource, freeFn := gl.Strs(string(src) + "\x00")
	defer freeFn()

	gl.ShaderSource(handle, 1, glSource, nil)
	gl.CompileShader(handle)

	var status int32
	gl.GetShaderiv(handle, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(handle, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(handle, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", filename, log)
	}

	return handle, nil
}
