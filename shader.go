package main

import (
	"bufio"
	"image"
	"os"
)

type Shader interface {
	shade(scene *Scene, triangle *ProjectedTriangle, coordinates [3]float64, z float64) (uint8, uint8, uint8)
}

type FastImage struct {
	height float64
	width  float64
	data   []uint32
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
	data := make([]uint32, height*width)

	idx := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := t.At(x, y).RGBA()
			v := uint32(0)
			v = uint32(r/256) & 255
			v |= (uint32(g/256) & 255) << 8
			v |= (uint32(b/256) & 255) << 16
			data[idx] = v
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

func (textureShader *TextureShader) shade(scene *Scene, triangle *ProjectedTriangle, coordinates [3]float64, z float64) (uint8, uint8, uint8) {
	l0 := coordinates[0]
	l1 := coordinates[1]
	l2 := coordinates[2]

	t := triangle

	u := l0*t.uvMapping[0][0]*(1/triangle.viewVerts[0].z) + l1*t.uvMapping[1][0]*(1/triangle.viewVerts[1].z) + l2*t.uvMapping[2][0]*(1/triangle.viewVerts[2].z)
	v := l0*t.uvMapping[0][1]*(1/triangle.viewVerts[0].z) + l1*t.uvMapping[1][1]*(1/triangle.viewVerts[1].z) + l2*t.uvMapping[2][1]*(1/triangle.viewVerts[2].z)
	u *= z
	v *= z

	if u < 0 || v < 0 {
		return 0, 0, 0
	}

	x := int(u*textureShader.texture.width) % int(textureShader.texture.width)
	y := int(v*textureShader.texture.height) % int(textureShader.texture.height)
	idx := (x + y*int(textureShader.texture.width))
	data := textureShader.texture.data[idx]

	r := float64(data & 255)
	g := float64((data >> 8) & 255)
	b := float64((data >> 16) & 255)

	intensity := l0*triangle.lightIntensity[0] + l1*triangle.lightIntensity[1] + l2*triangle.lightIntensity[2]

	return clampColor(r * intensity), clampColor(g * intensity), clampColor(b * intensity)
}

type LineShader struct {
	r         uint8
	g         uint8
	b         uint8
	lineR     uint8
	lineG     uint8
	lineB     uint8
	thickness float64
}

func (shader *LineShader) shade(scene *Scene, triangle *ProjectedTriangle, coordinates [3]float64, z float64) (uint8, uint8, uint8) {
	l0 := coordinates[0]
	l1 := coordinates[1]
	l2 := coordinates[2]

	if l0 < shader.thickness || l1 < shader.thickness || l2 < shader.thickness {
		return uint8(float64(shader.lineR)), uint8(float64(shader.lineG)), uint8(float64(shader.lineB))
	} else {
		return shader.r, shader.g, shader.b
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

func (smoothColor *SmoothColorShader) shade(scene *Scene, triangle *ProjectedTriangle, coordinates [3]float64, z float64) (uint8, uint8, uint8) {
	l0 := coordinates[0]
	l1 := coordinates[1]
	l2 := coordinates[2]

	intensity := l0*triangle.lightIntensity[0] + l1*triangle.lightIntensity[1] + l2*triangle.lightIntensity[2]
	intensity += .4

	if intensity < 0 {
		// Shoudln't be needed if there was occulsion culling or shadows ?
		return 0, 0, 0
	} else {
		return clampColor(intensity * smoothColor.r), clampColor(intensity * smoothColor.g), clampColor(intensity * smoothColor.b)
	}
}

type FlatGrayScaleShader struct {
}

func (flatGrayScaleShader *FlatGrayScaleShader) shade(scene *Scene, triangle *ProjectedTriangle, coordinates [3]float64, z float64) (uint8, uint8, uint8) {
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

type FlatShader struct {
	r uint8
	g uint8
	b uint8
}

func (shader *FlatShader) shade(scene *Scene, triangle *ProjectedTriangle, coordinates [3]float64, z float64) (uint8, uint8, uint8) {
	return uint8(shader.r), uint8(shader.g), uint8(shader.b)
}
