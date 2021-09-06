package main

import (
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	_ "github.com/ftrvxmtrx/tga"
)

type Model struct {
	triangles []*Triangle
	shader    Shader
}

type Triangle struct {
	worldVerts    []Vector3 // world space
	viewVerts     []Vector3 // view space relative to camera
	viewportVerts []Vector2 // ndc space relative to viewports [-1,1]
	normals       []Vector3
	viewNormals   []Vector3 // view/camera space
	inFrustrum    bool
	uvMapping     [][]float64
	invViewZ      []float64
}

func newTriangle() *Triangle {
	return &Triangle{
		worldVerts:    []Vector3{{}, {}, {}},
		normals:       []Vector3{{}, {}, {}},
		viewportVerts: []Vector2{{}, {}, {}},
		uvMapping:     [][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		invViewZ:      []float64{0, 0, 0},
		viewVerts:     []Vector3{{}, {}, {}},
		viewNormals:   []Vector3{{}, {}, {}},
	}
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
	textures := [][]float64{}

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
			u := parseFloat(splitted[2])
			v := 1 - parseFloat(splitted[3])
			w := parseFloat(splitted[4])
			textures = append(textures, []float64{u, v, w})
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

			triangle := newTriangle()
			triangle.worldVerts = []Vector3{vertex[vertexIdx1], vertex[vertexIdx2], vertex[vertexIdx3]}
			triangle.normals = []Vector3{normals[normalIdx1], normals[normalIdx2], normals[normalIdx3]}
			triangle.uvMapping = [][]float64{textures[textureIdx1], textures[textureIdx2], textures[textureIdx3]}

			triangles = append(triangles, triangle)
		}
	}

	return &Model{
		triangles: triangles,
		shader:    shader,
	}
}

func (model *Model) moveX(x float64) *Model {
	for _, triangle := range model.triangles {
		for i := range triangle.worldVerts {
			triangle.worldVerts[i].x += x
		}
	}
	return model
}

func (model *Model) moveY(y float64) *Model {
	for _, triangle := range model.triangles {
		for i := range triangle.worldVerts {
			triangle.worldVerts[i].y += y
		}
	}
	return model
}

func (model *Model) moveZ(z float64) *Model {
	for _, triangle := range model.triangles {
		for i := range triangle.worldVerts {
			triangle.worldVerts[i].z += z
		}
	}
	return model
}

func (model *Model) scale(x, y, z float64) *Model {
	for _, triangle := range model.triangles {
		for i := range triangle.worldVerts {
			triangle.worldVerts[i].x *= x
			triangle.worldVerts[i].y *= y
			triangle.worldVerts[i].z *= z
		}
	}
	return model
}

func (model *Model) scaleUV(u, v float64) *Model {
	for _, triangle := range model.triangles {
		for t := range triangle.uvMapping {
			triangle.uvMapping[t][0] *= u
			triangle.uvMapping[t][1] *= v
		}
	}
	return model
}

func (model *Model) rotateY(v float64) *Model {
	for _, triangle := range model.triangles {
		for i := range triangle.worldVerts {
			x := triangle.worldVerts[i].x
			y := triangle.worldVerts[i].y
			z := triangle.worldVerts[i].z
			triangle.worldVerts[i].x = x*math.Cos(v) + z*math.Sin(v)
			triangle.worldVerts[i].y = y
			triangle.worldVerts[i].z = x*-math.Sin(v) + z*math.Cos(v)
		}
	}
	return model
}

func (model *Model) rotateX(v float64) *Model {
	for _, triangle := range model.triangles {
		for i := range triangle.worldVerts {
			x := triangle.worldVerts[i].x
			y := triangle.worldVerts[i].y
			z := triangle.worldVerts[i].z
			triangle.worldVerts[i].x = x
			triangle.worldVerts[i].y = y*math.Cos(v) + z*math.Sin(v)
			triangle.worldVerts[i].z = y*-math.Sin(v) + z*math.Cos(v)
		}
	}
	return model
}

func (model *Model) rotateZ(v float64) *Model {
	for _, triangle := range model.triangles {
		for i := range triangle.worldVerts {
			x := triangle.worldVerts[i].x
			y := triangle.worldVerts[i].y
			z := triangle.worldVerts[i].z
			// TODO check this?
			triangle.worldVerts[i].x = y*-math.Sin(v) + x*math.Cos(v)
			triangle.worldVerts[i].y = y*math.Cos(v) + x*math.Sin(v)
			triangle.worldVerts[i].z = z
		}
	}
	return model
}
