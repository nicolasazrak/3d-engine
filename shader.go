package main

import (
	"bufio"
	"image"
	"math"
	"os"
)

type Shader interface {
	shade(scene *Scene, triangle *Triangle, l0 float64, l1 float64, l2 float64) (uint8, uint8, uint8)
}

type FastImage struct {
	height float64
	width  float64
	data   []float64
}

func decodeTexture(filename string) image.Image {
	f, err := os.Open("testdata/" + filename)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	i, _, err := image.Decode(bufio.NewReader(f))
	if err != nil {
		panic(err)
	}

	return i
}

func convertTexture(t image.Image) *FastImage {
	height := t.Bounds().Max.Y
	width := t.Bounds().Max.X
	data := make([]float64, height*width*4)

	idx := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := t.At(x, y).RGBA()
			data[idx] = float64(r) / 256
			idx++
			data[idx] = float64(g) / 256
			idx++
			data[idx] = float64(b) / 256
			idx++
			idx++
		}
	}

	return &FastImage{
		height: float64(height),
		width:  float64(width),
		data:   data,
	}
}

func newTextureShader(texturePath string) Shader {
	tF, err := os.Open(texturePath)
	if err != nil {
		panic(err)
	}
	texture, _, err := image.Decode(tF)
	if err != nil {
		panic(err)
	}

	return &TextureShader{
		texture: convertTexture(texture),
	}
}

type TextureShader struct {
	texture *FastImage
}

func (textureShader *TextureShader) shade(scene *Scene, triangle *Triangle, l0 float64, l1 float64, l2 float64) (uint8, uint8, uint8) {
	p := ponderate(triangle.viewVerts, []float64{l0, l1, l2})
	normal := ponderate(triangle.viewNormals, []float64{l0, l1, l2})
	lightNormal := normalize(minus(scene.projectedLight, p))
	intensity := dotProduct(lightNormal, normal)

	if intensity < 0 {
		// Shoudln't be needed if there was occulsion culling or shadows ?
		return 0, 0, 0
	} else {
		vert := ponderate(triangle.uvMapping, []float64{l0, l1, l2})
		x := int(vert.x * textureShader.texture.width)
		y := int(vert.y * textureShader.texture.height)
		idx := (x + y*int(textureShader.texture.width)) * 4
		r := textureShader.texture.data[idx]
		g := textureShader.texture.data[idx+1]
		b := textureShader.texture.data[idx+2]
		return uint8(float64(r) * intensity), uint8(float64(g) * intensity), uint8(float64(b) * intensity)
	}
}

type IntensityShader struct {
}

func (intensity *IntensityShader) shade(scene *Scene, triangle *Triangle, l0 float64, l1 float64, l2 float64) (uint8, uint8, uint8) {
	return uint8(l0 * 255), uint8(l1 * 255), uint8(l2 * 255)
}

type SmoothColorShader struct {
	r float64
	g float64
	b float64
}

func (smoothColor *SmoothColorShader) shade(scene *Scene, triangle *Triangle, l0 float64, l1 float64, l2 float64) (uint8, uint8, uint8) {
	p := ponderate(triangle.viewVerts, []float64{l0, l1, l2})
	normal := ponderate(triangle.viewNormals, []float64{l0, l1, l2})
	lightNormal := normalize(minus(scene.projectedLight, p))
	distance := 1 / math.Sqrt(norm(minus(p, scene.projectedLight)))
	intensity := dotProduct(lightNormal, normal) * distance

	if intensity < 0 {
		// Shoudln't be needed if there was occulsion culling or shadows ?
		return 0, 0, 0
	} else {
		return uint8(intensity * smoothColor.r), uint8(intensity * smoothColor.g), uint8(intensity * smoothColor.b)
	}
}

type FlatGrayScaleShader struct {
}

func (flatGrayScaleShader *FlatGrayScaleShader) shade(scene *Scene, triangle *Triangle, l0 float64, l1 float64, l2 float64) (uint8, uint8, uint8) {
	p := triangle.viewVerts[0]
	normal := triangle.viewNormals[0]
	lightNormal := normalize(minus(scene.projectedLight, p))
	intensity := dotProduct(lightNormal, normal)

	if intensity < 0 {
		// Shoudln't be needed if there was occulsion culling or shadows ?
		return 0, 0, 0
	} else {
		return uint8(intensity * 255), uint8(intensity * 255), uint8(intensity * 255)
	}
}
