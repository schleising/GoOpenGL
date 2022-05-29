package shapes

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	numTriangles           = 2
	NumVertices            = 4
	sizeOfUint32           = 4
	defaultTextureLocation = "images/Loading Icon.png"
)

var (
	DefaultTexture *Texture
)

type Rectangle struct {
	xPos     float32
	yPos     float32
	width    int
	height   int
	vertices []Vertex
	indices  []uint32
	Handle   uint32
	texture  uint32
	scaleX   float32
	scaleY   float32
	angle    float32
}

func NewRectangle(x, y float32, width, height int) *Rectangle {
	r := new(Rectangle)
	r.xPos = x
	r.yPos = y
	r.width = width
	r.height = height
	r.scaleX = 1
	r.scaleY = 1
	r.angle = 0
	r.createVertices()
	r.SetDefaultTexture()
	r.Handle = r.MakeVao()
	return r
}

func (r *Rectangle) UpdateXPos(dx float32) {
	r.xPos += dx / r.scaleX
}

func (r *Rectangle) UpdateYPos(dy float32) {
	r.yPos += dy / r.scaleY
}

func (r *Rectangle) Pos() (float32, float32) {
	return r.xPos, r.yPos
}

func (r *Rectangle) Size() (int, int) {
	return r.width, r.height
}

func (r *Rectangle) Scale(scale float32) {
	r.scaleX = scale
	r.scaleY = scale
}

func (r *Rectangle) createVertices() {
	glBlx := float32(0.0)
	glBly := float32(r.height)

	glTlx := float32(0.0)
	glTly := float32(0.0)

	glTrx := float32(r.width)
	glTry := float32(0.0)

	glBrx := float32(r.width)
	glBry := float32(r.height)

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

func (r *Rectangle) MakeVao() uint32 {
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

func (r *Rectangle) SetDefaultTexture() {

	if DefaultTexture == nil {
		imageChan := make(chan ImageMessage)
		var err error
		DefaultTexture, err = NewTexture(defaultTextureLocation, imageChan)

		if err != nil {
			panic(err)
		}

		DefaultTexture.SetTexCoords(r.width, r.height)
	}

	r.SetTexCoords(DefaultTexture)
}

func (r *Rectangle) RequestTexture(filename string, imageChan chan ImageMessage) {
	go LoadImage(filename, imageChan)
}

func (r *Rectangle) SetTexCoords(texture *Texture) {
	r.texture = texture.Handle

	r.vertices[0].TexCoord.S = texture.TexCoords[0].S
	r.vertices[0].TexCoord.T = texture.TexCoords[0].T
	r.vertices[1].TexCoord.S = texture.TexCoords[1].S
	r.vertices[1].TexCoord.T = texture.TexCoords[1].T
	r.vertices[2].TexCoord.S = texture.TexCoords[2].S
	r.vertices[2].TexCoord.T = texture.TexCoords[2].T
	r.vertices[3].TexCoord.S = texture.TexCoords[3].S
	r.vertices[3].TexCoord.T = texture.TexCoords[3].T
}

func (r *Rectangle) Draw(program uint32) {
	gl.BindTexture(gl.TEXTURE_2D, r.texture)

	gl.UseProgram(program)

	gl.BindVertexArray(r.Handle)

	scaleMat := mgl32.Scale3D(r.scaleX, r.scaleY, 0.0)
	scale := gl.GetUniformLocation(program, gl.Str("scale"+"\x00"))
	gl.UniformMatrix4fv(scale, 1, false, &scaleMat[0])

	rotMat := mgl32.HomogRotate3DZ(r.angle)
	rotation := gl.GetUniformLocation(program, gl.Str("rotation"+"\x00"))
	gl.UniformMatrix4fv(rotation, 1, false, &rotMat[0])

	transMat := mgl32.Translate3D(r.xPos*r.scaleX, r.yPos*r.scaleY, 0.0)
	translation := gl.GetUniformLocation(program, gl.Str("translation"+"\x00"))
	gl.UniformMatrix4fv(translation, 1, false, &transMat[0])

	gl.DrawElements(gl.TRIANGLES, PointLen*numTriangles, gl.UNSIGNED_INT, gl.Ptr(nil))

	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (r *Rectangle) ClickInRect(x, y float32) bool {
	if x >= r.xPos*r.scaleX && x <= ((r.xPos+float32(r.width))*r.scaleX) && y >= r.yPos*r.scaleY && y <= ((r.yPos+float32(r.height))*r.scaleY) {
		return true
	} else {
		return false
	}
}
