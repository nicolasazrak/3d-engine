package main

import (
	"bufio"
	"image"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	_ "github.com/ftrvxmtrx/tga"
)

type FastImage struct {
	height float64
	width  float64
	data   []float64
}

type Model struct {
	triangles []Triangle
	texture   *FastImage
}

type Triangle struct {
	verts     []Vector3
	normals   []Vector3
	uvMapping []Vector3
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

func parseInt(s string) int {
	vertexIdx1, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		print(s)
		panic(err)
	}
	return int(vertexIdx1)
}

func parseFloat(s string) float64 {
	vertexIdx1, err := strconv.ParseFloat(s, 64)
	if err != nil {
		print(s)
		panic(err)
	}
	return vertexIdx1
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

func parseModel(objPath string, texturePath string) *Model {
	f, err := ioutil.ReadFile(objPath)
	if err != nil {
		panic(err)
	}

	str := string(f)

	triangles := []Triangle{}
	vertex := []Vector3{}
	normals := []Vector3{}
	textures := []Vector3{}

	for _, line := range strings.Split(str, "\n") {
		if strings.HasPrefix(line, "v ") {
			splitted := strings.Split(line, " ")
			x := parseFloat(splitted[1])
			y := parseFloat(splitted[2])
			z := parseFloat(splitted[3])
			vertex = append(vertex, Vector3{x, y, z})
		}

		if strings.HasPrefix(line, "vn ") {
			splitted := strings.Split(line, " ")
			x := parseFloat(splitted[2])
			y := parseFloat(splitted[3])
			z := parseFloat(splitted[4])
			normals = append(normals, normalize(Vector3{x, y, z}))
		}

		if strings.HasPrefix(line, "vt ") {
			splitted := strings.Split(line, " ")
			x := parseFloat(splitted[2])
			y := 1 - parseFloat(splitted[3])
			z := parseFloat(splitted[4])
			textures = append(textures, Vector3{x, y, z})
		}

		if strings.HasPrefix(line, "f ") {
			splitted := strings.Split(line, " ")
			vertexIdx1 := parseInt(strings.Split(splitted[1], "/")[0]) - 1
			vertexIdx2 := parseInt(strings.Split(splitted[2], "/")[0]) - 1
			vertexIdx3 := parseInt(strings.Split(splitted[3], "/")[0]) - 1

			textureIdx1 := parseInt(strings.Split(splitted[1], "/")[1]) - 1
			textureIdx2 := parseInt(strings.Split(splitted[2], "/")[1]) - 1
			textureIdx3 := parseInt(strings.Split(splitted[3], "/")[1]) - 1

			normalIdx1 := parseInt(strings.Split(splitted[1], "/")[2]) - 1
			normalIdx2 := parseInt(strings.Split(splitted[2], "/")[2]) - 1
			normalIdx3 := parseInt(strings.Split(splitted[3], "/")[2]) - 1

			triangles = append(triangles, Triangle{
				verts:     []Vector3{vertex[vertexIdx1], vertex[vertexIdx2], vertex[vertexIdx3]},
				normals:   []Vector3{normals[normalIdx1], normals[normalIdx2], normals[normalIdx3]},
				uvMapping: []Vector3{textures[textureIdx1], textures[textureIdx2], textures[textureIdx3]},
			})
		}
	}

	tF, err := os.Open(texturePath)
	if err != nil {
		panic(err)
	}
	texture, _, err := image.Decode(tF)
	if err != nil {
		panic(err)
	}

	return &Model{
		triangles: triangles,
		texture:   convertTexture(texture),
	}
}
