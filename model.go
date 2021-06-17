package main

import (
	"io/ioutil"
	"strconv"
	"strings"

	_ "github.com/ftrvxmtrx/tga"
)

type Model struct {
	triangles []*Triangle
	shader    Shader
}

type Triangle struct {
	verts      []Vector3
	projection []Vector2
	normals    []Vector3
	uvMapping  []Vector3
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

func parseModel(objPath string, shader Shader) *Model {
	f, err := ioutil.ReadFile(objPath)
	if err != nil {
		panic(err)
	}

	str := string(f)

	triangles := []*Triangle{}
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

			triangles = append(triangles, &Triangle{
				verts:      []Vector3{vertex[vertexIdx1], vertex[vertexIdx2], vertex[vertexIdx3]},
				normals:    []Vector3{normals[normalIdx1], normals[normalIdx2], normals[normalIdx3]},
				uvMapping:  []Vector3{textures[textureIdx1], textures[textureIdx2], textures[textureIdx3]},
				projection: []Vector2{{}, {}, {}},
			})
		}
	}

	return &Model{
		triangles: triangles,
		shader:    shader,
	}
}

func (model *Model) moveX(x float64) {
	for _, triangle := range model.triangles {
		for i := range triangle.verts {
			triangle.verts[i].x += x
		}
	}
}

func (model *Model) moveY(y float64) {
	for _, triangle := range model.triangles {
		for i := range triangle.verts {
			triangle.verts[i].y += y
		}
	}
}

func (model *Model) moveZ(z float64) {
	for _, triangle := range model.triangles {
		for i := range triangle.verts {
			triangle.verts[i].z += z
		}
	}
}
