package shapes

const (
	float32Size = 4
	PointLen    = 3
	ColourLen   = 3
	TexCoordLen = 2

	PointSize    = PointLen * float32Size
	ColourSize   = ColourLen * float32Size
	TexCoordSize = TexCoordLen * float32Size

	VertexSize = PointSize + ColourSize + TexCoordSize
)

type Point struct {
	X float32
	Y float32
	Z float32
}

type Colour struct {
	R float32
	G float32
	B float32
}

type TexCoord struct {
	S float32
	T float32
}

type Vertex struct {
	Point    Point
	Colour   Colour
	TexCoord TexCoord
}
