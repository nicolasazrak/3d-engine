package main

import (
	"bufio"
	"image"
	"os"
)

type Shader interface {
	shade(scene *Scene, triangle *Triangle, coordinates [3]float64, z float64) (uint8, uint8, uint8)
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

func getFastTexture(texturePath string) *FastImage {
	tF, err := os.Open(texturePath)
	if err != nil {
		panic(err)
	}
	texture, _, err := image.Decode(tF)
	if err != nil {
		panic(err)
	}
	return convertTexture(texture)
}

func newTextureShader(texturePath string) Shader {
	return &TextureShader{
		texture: getFastTexture(texturePath),
	}
}

type TextureShader struct {
	texture *FastImage
}

func clampColor(a float64) uint8 {
	if a > 255 {
		return 255
	}
	return uint8(a)
}

func (textureShader *TextureShader) shade(scene *Scene, triangle *Triangle, coordinates [3]float64, z float64) (uint8, uint8, uint8) {
	l0 := coordinates[0]
	l1 := coordinates[1]
	l2 := coordinates[2]

	t := triangle
	p := ponderate(triangle.viewVerts, coordinates[:])
	p.z = z
	intensity := 1. / norm(minus(p, scene.projectedLight))
	intensity += 0.4 // ambient

	if intensity < 0 {
		// Shoudln't be needed if there was occulsion culling or shadows ?
		return 0, 0, 0
	} else {
		u := l0*t.uvMappingCorrected[0][0] + l1*t.uvMappingCorrected[1][0] + l2*t.uvMappingCorrected[2][0]
		v := l0*t.uvMappingCorrected[0][1] + l1*t.uvMappingCorrected[1][1] + l2*t.uvMappingCorrected[2][1]
		u *= z
		v *= z

		if u < 0 || v < 0 || u >= 1 || v >= 1 {
			// TODO remove this
			return 0, 0, 0
		}

		x := int(u * textureShader.texture.width)
		y := int(v * textureShader.texture.height)
		idx := (x + y*int(textureShader.texture.width)) * 4
		r := textureShader.texture.data[idx]
		g := textureShader.texture.data[idx+1]
		b := textureShader.texture.data[idx+2]

		return clampColor(r * intensity), clampColor(g * intensity), clampColor(b * intensity)
	}
}

type LineShader struct {
}

func (shader *LineShader) shade(scene *Scene, triangle *Triangle, coordinates [3]float64, z float64) (uint8, uint8, uint8) {
	l0 := coordinates[0]
	l1 := coordinates[1]
	l2 := coordinates[2]

	if l0 < 0.001 || l1 < 0.001 || l2 < 0.001 {
		return uint8(l0 * 255), uint8(l1 * 255), uint8(l2 * 255)
	} else {
		return 0, 0, 0
	}
}

type IntensityShader struct {
}

func (intensity *IntensityShader) shade(scene *Scene, triangle *Triangle, coordinates [3]float64, z float64) (uint8, uint8, uint8) {
	l0 := coordinates[0]
	l1 := coordinates[1]
	l2 := coordinates[2]

	return uint8(l0 * 255), uint8(l1 * 255), uint8(l2 * 255)
}

type SmoothColorShader struct {
	r float64
	g float64
	b float64
}

func (smoothColor *SmoothColorShader) shade(scene *Scene, triangle *Triangle, coordinates [3]float64, z float64) (uint8, uint8, uint8) {
	p := ponderate(triangle.viewVerts, coordinates[:])
	normal := ponderate(triangle.viewNormals, coordinates[:])
	lightNormal := normalize(minus(scene.projectedLight, p))
	intensity := dotProduct(lightNormal, normal)

	if intensity < 0 {
		// Shoudln't be needed if there was occulsion culling or shadows ?
		return 0, 0, 0
	} else {
		return uint8(intensity * smoothColor.r), uint8(intensity * smoothColor.g), uint8(intensity * smoothColor.b)
	}
}

type FlatGrayScaleShader struct {
}

func (flatGrayScaleShader *FlatGrayScaleShader) shade(scene *Scene, triangle *Triangle, coordinates [3]float64, z float64) (uint8, uint8, uint8) {
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
