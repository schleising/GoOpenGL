package shapes

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/schleising/GoOpenGL/screen"
)

const (
	numTriangles = 2
	NumVertices  = 4
	sizeOfUint32 = 4
)

type Rectangle struct {
	XPos     float32
	YPos     float32
	Width    int
	Height   int
	Vertices []Vertex
	Indices  []uint32
	screen   screen.Screen
	Handle   uint32
	texture  uint32
	program  uint32
}

func (r *Rectangle) Create(x, y float32, width, height int, screen screen.Screen) {
	r.XPos = x
	r.YPos = y
	r.Width = width
	r.Height = height
	r.screen = screen
	r.createVertices()
	r.Handle = r.makeVao()
}

func (r *Rectangle) Pos() (float32, float32) {
	return r.XPos, r.YPos
}

func (r *Rectangle) Size() (int, int) {
	return r.Width, r.Height
}

func (r *Rectangle) SetProgram(program uint32) {
	r.program = program
}

func (r *Rectangle) createVertices() {
	// Bottom Left
	glBlx, glBly := mgl32.ScreenToGLCoords(-r.Width/2+r.screen.Width/2, -r.Height/2+r.screen.Height/2, r.screen.Width, r.screen.Height)
	glTlx, glTly := mgl32.ScreenToGLCoords(-r.Width/2+r.screen.Width/2, r.Height/2+r.screen.Height/2, r.screen.Width, r.screen.Height)
	glTrx, glTry := mgl32.ScreenToGLCoords(r.Width/2+r.screen.Width/2, r.Height/2+r.screen.Height/2, r.screen.Width, r.screen.Height)
	glBrx, glBry := mgl32.ScreenToGLCoords(r.Width/2+r.screen.Width/2, -r.Height/2+r.screen.Height/2, r.screen.Width, r.screen.Height)

	blp := Point{
		X: glBlx,
		Y: glBly,
		Z: 0,
	}

	blc := Colour{
		R: 1.0,
		G: 0.0,
		B: 0.0,
	}

	blt := TexCoord{
		S: 0.0,
		T: 0.0,
	}

	blv := Vertex{
		Point:    blp,
		Colour:   blc,
		TexCoord: blt,
	}

	tlp := Point{
		X: glTlx,
		Y: glTly,
		Z: 0,
	}

	tlc := Colour{
		R: 0.0,
		G: 1.0,
		B: 0.0,
	}

	tlt := TexCoord{
		S: 0.0,
		T: 1.0,
	}

	tlv := Vertex{
		Point:    tlp,
		Colour:   tlc,
		TexCoord: tlt,
	}

	trp := Point{
		X: glTrx,
		Y: glTry,
		Z: 0,
	}

	trc := Colour{
		R: 0.0,
		G: 0.0,
		B: 1.0,
	}

	trt := TexCoord{
		S: 1.0,
		T: 1.0,
	}

	trv := Vertex{
		Point:    trp,
		Colour:   trc,
		TexCoord: trt,
	}

	brp := Point{
		X: glBrx,
		Y: glBry,
		Z: 0,
	}

	brc := Colour{
		R: 1.0,
		G: 1.0,
		B: 1.0,
	}

	brt := TexCoord{
		S: 1.0,
		T: 0.0,
	}

	brv := Vertex{
		Point:    brp,
		Colour:   brc,
		TexCoord: brt,
	}

	r.Vertices = []Vertex{blv, tlv, trv, brv}
	r.Indices = []uint32{0, 1, 2, 0, 3, 2}
}

func (r *Rectangle) makeVao() uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, NumVertices*VertexSize, gl.Ptr(r.Vertices), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(r.Indices)*sizeOfUint32, gl.Ptr(r.Indices), gl.STATIC_DRAW)

	var offset uintptr = 0
	gl.VertexAttribPointerWithOffset(0, PointLen, gl.FLOAT, false, VertexSize, offset)
	gl.EnableVertexAttribArray(0)

	offset = PointSize
	gl.VertexAttribPointerWithOffset(1, ColourLen, gl.FLOAT, false, VertexSize, offset)
	gl.EnableVertexAttribArray(1)

	offset = PointSize + ColourSize
	gl.VertexAttribPointerWithOffset(2, TexCoordLen, gl.FLOAT, false, VertexSize, offset)
	gl.EnableVertexAttribArray(2)

	return vao
}

func (r *Rectangle) SetTexture(texture uint32) {
	r.texture = texture
}

func (r *Rectangle) Draw() {
	gl.BindTexture(gl.TEXTURE_2D, r.texture)

	gl.UseProgram(r.program)

	gl.BindVertexArray(r.Handle)

	glX, glY := mgl32.ScreenToGLCoords(int(r.XPos), int(r.YPos), r.screen.Width, r.screen.Height)

	transMat := mgl32.Translate3D(glX, glY, 0.0)
	// rotMat := mgl32.HomogRotate3DZ(mgl32.DegToRad(float32(count)))

	// transMat = transMat.Mul4(rotMat)

	translation := gl.GetUniformLocation(r.program, gl.Str("translation"+"\x00"))
	gl.UniformMatrix4fv(translation, 1, false, &transMat[0])

	gl.DrawElements(gl.TRIANGLES, PointLen*numTriangles, gl.UNSIGNED_INT, gl.Ptr(nil))

	gl.BindVertexArray(0)
}
