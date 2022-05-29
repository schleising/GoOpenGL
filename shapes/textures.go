package shapes

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Texture struct {
	Handle    uint32
	Width     int
	Height    int
	TexCoords []TexCoord
}

type ImageMessage struct {
	Filename string
	Rgba     *image.RGBA
}

func NewTexture(filename string, imageChan chan ImageMessage) (*Texture, error) {
	go LoadImage(filename, imageChan)

	imageMessage := <-imageChan

	texture := CreateTextureFromImageRgba(imageMessage.Rgba)

	return texture, nil
}

func LoadImage(filename string, imageChan chan ImageMessage) {
	loadedImage, err := os.Open(filename)
	fmt.Printf("Loading %v\n", filename)

	if err != nil {
		fmt.Printf("Error loading image: %v\n", filename)
		imageChan <- ImageMessage{Filename: filename, Rgba: nil}
	}

	defer loadedImage.Close()

	img, _, err := image.Decode(loadedImage)

	if err != nil {
		fmt.Printf("Error decoding image: %v\n", filename)
		imageChan <- ImageMessage{Filename: filename, Rgba: nil}
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)

	fmt.Printf("Loaded  %v\n", filename)
	imageChan <- ImageMessage{Filename: filename, Rgba: rgba}
}

func CreateTextureFromImageRgba(rgba *image.RGBA) *Texture {
	var texture Texture
	gl.GenTextures(1, &texture.Handle)
	gl.BindTexture(gl.TEXTURE_2D, texture.Handle)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	borderColor := []float32{1.0, 1.0, 1.0, 0.0}
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Rect.Dx()), int32(rgba.Rect.Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	texture.Width = rgba.Rect.Dx()
	texture.Height = rgba.Rect.Dy()

	return &texture
}

func (t *Texture) SetTexCoords(containerWidth, containerHeight int) {
	xScale := float32(containerWidth) / float32(t.Width)
	yScale := float32(containerHeight) / float32(t.Height)

	var scale float32
	if xScale < yScale {
		scale = xScale
	} else {
		scale = yScale
	}

	scaledWidth := float32(t.Width) * scale
	scaledHeight := float32(t.Height) * scale

	widthFactor := (float32(containerWidth)/scaledWidth - 1) / 2
	heightFactor := (float32(containerHeight)/scaledHeight - 1) / 2

	normalisedStartX := -widthFactor
	normalisedStartY := -heightFactor
	normalisedEndX := 1 + widthFactor
	normalisedEndY := 1 + heightFactor

	var blt, tlt, trt, brt TexCoord
	blt.S = normalisedStartX
	blt.T = normalisedEndY

	tlt.S = normalisedStartX
	tlt.T = normalisedStartY

	trt.S = normalisedEndX
	trt.T = normalisedStartY

	brt.S = normalisedEndX
	brt.T = normalisedEndY

	t.TexCoords = []TexCoord{blt, tlt, trt, brt}
}
