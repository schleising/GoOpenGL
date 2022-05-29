package shapes

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	numTriangles           = 2
	NumVertices            = 4
	sizeOfUint32           = 4
	defaultTextureLocation = "images/1x1_#000000ff.png"
)

type Rectangle struct {
	XPos     float32
	YPos     float32
	Width    int
	Height   int
	vertices []Vertex
	indices  []uint32
	handle   uint32
	texture  uint32
}

func NewRectangle(x, y float32, width, height int) *Rectangle {
	r := new(Rectangle)
	r.XPos = x
	r.YPos = y
	r.Width = width
	r.Height = height
	r.createVertices()
	r.SetTexture(defaultTextureLocation)
	r.handle = r.makeVao()
	return r
}

func (r *Rectangle) Pos() (float32, float32) {
	return r.XPos, r.YPos
}

func (r *Rectangle) Size() (int, int) {
	return r.Width, r.Height
}

func (r *Rectangle) createVertices() {
	glBlx := float32(0.0)
	glBly := float32(r.Height)

	glTlx := float32(0.0)
	glTly := float32(0.0)

	glTrx := float32(r.Width)
	glTry := float32(0.0)

	glBrx := float32(r.Width)
	glBry := float32(r.Height)

	blp := Point{X: glBlx, Y: glBly, Z: 0}
	blv := Vertex{Point: blp}

	tlp := Point{X: glTlx, Y: glTly, Z: 0}
	tlv := Vertex{Point: tlp}

	trp := Point{X: glTrx, Y: glTry, Z: 0}
	trv := Vertex{Point: trp}

	brp := Point{X: glBrx, Y: glBry, Z: 0}
	brv := Vertex{Point: brp}

	r.vertices = []Vertex{blv, tlv, trv, brv}
	r.indices = []uint32{0, 1, 2, 0, 3, 2}
}

func (r *Rectangle) makeVao() uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, NumVertices*VertexSize, gl.Ptr(r.vertices), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(r.indices)*sizeOfUint32, gl.Ptr(r.indices), gl.STATIC_DRAW)

	var offset uintptr = 0
	gl.VertexAttribPointerWithOffset(0, PointLen, gl.FLOAT, false, VertexSize, offset)
	gl.EnableVertexAttribArray(0)

	offset = PointSize
	gl.VertexAttribPointerWithOffset(1, ColourLen, gl.FLOAT, false, VertexSize, offset)
	gl.EnableVertexAttribArray(1)

	offset = PointSize + ColourSize
	gl.VertexAttribPointerWithOffset(2, TexCoordLen, gl.FLOAT, false, VertexSize, offset)
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	return vao
}

func (r *Rectangle) SetTexture(filename string) {
	texture, err := NewTexture(filename, r.Width, r.Height)

	if err == nil {
		r.texture = texture.Handle

		r.vertices[0].TexCoord.S = texture.TexCoords[0].S
		r.vertices[0].TexCoord.T = texture.TexCoords[0].T
		r.vertices[1].TexCoord.S = texture.TexCoords[1].S
		r.vertices[1].TexCoord.T = texture.TexCoords[1].T
		r.vertices[2].TexCoord.S = texture.TexCoords[2].S
		r.vertices[2].TexCoord.T = texture.TexCoords[2].T
		r.vertices[3].TexCoord.S = texture.TexCoords[3].S
		r.vertices[3].TexCoord.T = texture.TexCoords[3].T

		r.handle = r.makeVao()
	}
}

func (r *Rectangle) Draw(program uint32) {
	gl.BindTexture(gl.TEXTURE_2D, r.texture)

	gl.UseProgram(program)

	gl.BindVertexArray(r.handle)

	transMat := mgl32.Translate3D(r.XPos, r.YPos, 0.0)
	translation := gl.GetUniformLocation(program, gl.Str("translation"+"\x00"))
	gl.UniformMatrix4fv(translation, 1, false, &transMat[0])

	gl.DrawElements(gl.TRIANGLES, PointLen*numTriangles, gl.UNSIGNED_INT, gl.Ptr(nil))

	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (r *Rectangle) ClickInRect(x, y float32) bool {
	if x >= r.XPos && x <= r.XPos+float32(r.Width) && y >= r.YPos && y <= r.YPos+float32(r.Height) {
		return true
	} else {
		return false
	}
}
