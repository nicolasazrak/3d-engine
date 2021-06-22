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

func newBilinearTextureShader(texturePath string) Shader {
	return &BilinearTextureShader{
		texture: getFastTexture(texturePath),
	}
}

type TextureShader struct {
	texture *FastImage
}

func (textureShader *TextureShader) shade(scene *Scene, triangle *Triangle, l0 float64, l1 float64, l2 float64) (uint8, uint8, uint8) {
	p := ponderate(triangle.cameraVerts, []float64{l0, l1, l2})
	normal := ponderate(triangle.cameraNormals, []float64{l0, l1, l2})
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

type BilinearTextureShader struct {
	texture *FastImage
}

func ponderate2(a Vector3, b Vector3, aProportion float64) Vector3 {
	return Vector3{
		x: a.x*aProportion + b.x*(1-aProportion),
		y: a.y*aProportion + b.y*(1-aProportion),
		z: a.z*aProportion + b.z*(1-aProportion),
	}
}

func (bilinearTexture *BilinearTextureShader) shade(scene *Scene, triangle *Triangle, l0 float64, l1 float64, l2 float64) (uint8, uint8, uint8) {
	p := ponderate(triangle.cameraVerts, []float64{l0, l1, l2})
	intensity := 1. // norm(minus(p, scene.projectedLight))

	if intensity < 0 {
		// Shoudln't be needed if there was occulsion culling or shadows ?
		return 0, 0, 0
	} else {
		p0 := triangle.cameraVerts[0]
		p1 := triangle.cameraVerts[1]
		p3 := triangle.cameraVerts[2]
		p2 := triangle.quad.cameraVerts[0]

		n0Vec := minus(p0, p3)
		n1Vec := minus(p1, p0)
		n2Vec := minus(p2, p1)
		n3Vec := minus(p3, p2)

		n0 := normalize(Vector3{x: -n0Vec.y, y: n0Vec.x, z: 0})
		n1 := normalize(Vector3{x: -n1Vec.y, y: n1Vec.x, z: 0})
		n2 := normalize(Vector3{x: -n2Vec.y, y: n2Vec.x, z: 0})
		n3 := normalize(Vector3{x: -n3Vec.y, y: n3Vec.x, z: 0})

		u := dotProduct(minus(p, p0), n0) / (dotProduct(minus(p, p0), n0) + dotProduct(minus(p, p2), n2))
		v := dotProduct(minus(p, p0), n1) / (dotProduct(minus(p, p0), n1) + dotProduct(minus(p, p3), n3))

		x := int(math.Min(u*bilinearTexture.texture.width, bilinearTexture.texture.width-1))
		y := int(math.Min(v*bilinearTexture.texture.height, bilinearTexture.texture.height-1))
		idx := (x + y*int(bilinearTexture.texture.width)) * 4
		r := bilinearTexture.texture.data[idx]
		g := bilinearTexture.texture.data[idx+1]
		b := bilinearTexture.texture.data[idx+2]
		return uint8(math.Min(float64(r)*intensity, 255)), uint8(math.Min(float64(g)*intensity, 255)), uint8(math.Min(float64(b)*intensity, 255))
	}
}

type LineShader struct {
}

func (shader *LineShader) shade(scene *Scene, triangle *Triangle, l0 float64, l1 float64, l2 float64) (uint8, uint8, uint8) {
	if l0 < 0.001 || l1 < 0.001 || l2 < 0.001 {
		return uint8(l0 * 255), uint8(l1 * 255), uint8(l2 * 255)
	} else {
		return 0, 0, 0
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
	p := ponderate(triangle.cameraVerts, []float64{l0, l1, l2})
	normal := ponderate(triangle.cameraNormals, []float64{l0, l1, l2})
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

func (flatGrayScaleShader *FlatGrayScaleShader) shade(scene *Scene, triangle *Triangle, l0 float64, l1 float64, l2 float64) (uint8, uint8, uint8) {
	p := triangle.cameraVerts[0]
	normal := triangle.cameraNormals[0]
	lightNormal := normalize(minus(scene.projectedLight, p))
	intensity := dotProduct(lightNormal, normal)

	if intensity < 0 {
		// Shoudln't be needed if there was occulsion culling or shadows ?
		return 0, 0, 0
	} else {
		return uint8(intensity * 255), uint8(intensity * 255), uint8(intensity * 255)
	}
}
