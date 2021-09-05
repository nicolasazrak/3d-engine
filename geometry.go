package main

import "math"

func newXZSquare(size float64, shader Shader) *Model {
	v0 := Vector3{x: -size / 2, y: 0, z: -size / 2}
	v1 := Vector3{x: size / 2, y: 0, z: -size / 2}
	v2 := Vector3{x: size / 2, y: 0, z: size / 2}
	v3 := Vector3{x: -size / 2, y: 0, z: size / 2}

	t0 := []Vector3{v1, v0, v2}
	t1 := []Vector3{v3, v2, v0}

	normal := Vector3{x: 0, y: 1, z: 0}

	textt0v0 := []float64{0, 0, 0}
	textt1v0 := []float64{0, 0, 0}
	textt0v1 := []float64{0.999, 0, 0}
	textt0v2 := []float64{0.999, 0.999, 0}
	textt1v2 := []float64{0.999, 0.999, 0}
	textt1v3 := []float64{0, 0.999, 0}

	triangle1 := newTriangle()
	triangle1.worldVerts = t0
	triangle1.normals = []Vector3{normal, normal, normal}
	triangle1.uvMapping = [][]float64{textt0v1, textt0v0, textt0v2}

	triangle2 := newTriangle()
	triangle2.worldVerts = t1
	triangle2.normals = []Vector3{normal, normal, normal}
	triangle2.uvMapping = [][]float64{textt1v3, textt1v2, textt1v0}

	return &Model{
		triangles: []*Triangle{triangle1, triangle2},
		shader:    shader,
	}
}

func newXYSquare(size float64, shader Shader) *Model {
	pos := size / 2
	neg := -size / 2

	v0 := Vector3{x: neg, y: pos, z: 0}
	v1 := Vector3{x: neg, y: neg, z: 0}
	v2 := Vector3{x: pos, y: neg, z: 0}
	v3 := Vector3{x: pos, y: pos, z: 0}

	t0 := []Vector3{v1, v2, v0}
	t1 := []Vector3{v3, v0, v2}

	normal := Vector3{x: 0, y: 0, z: 1}

	textt0v0 := []float64{0, 0.999, 0}
	textt1v0 := []float64{0, 0.999, 0}
	textt0v1 := []float64{0, 0, 0}
	textt0v2 := []float64{0.999, 0, 0}
	textt1v2 := []float64{0.999, 0, 0}
	textt1v3 := []float64{0.999, 0.999, 0}

	triangle1 := newTriangle()
	triangle1.worldVerts = t0
	triangle1.normals = []Vector3{normal, normal, normal}
	triangle1.uvMapping = [][]float64{textt0v1, textt0v2, textt0v0}

	triangle2 := newTriangle()
	triangle2.worldVerts = t1
	triangle2.normals = []Vector3{normal, normal, normal}
	triangle2.uvMapping = [][]float64{textt1v3, textt1v0, textt1v2}

	return &Model{
		triangles: []*Triangle{triangle1, triangle2},
		shader:    shader,
	}
}

func newYZSquare(size float64, shader Shader) *Model {
	pos := size / 2
	neg := -size / 2

	v0 := Vector3{x: 0, y: pos, z: neg}
	v1 := Vector3{x: 0, y: neg, z: neg}
	v2 := Vector3{x: 0, y: neg, z: pos}
	v3 := Vector3{x: 0, y: pos, z: pos}

	t0 := []Vector3{v1, v2, v0}
	t1 := []Vector3{v3, v0, v2}

	normal := Vector3{x: 0, y: 0, z: 1}

	textt0v0 := []float64{0, 0.999, 0}
	textt1v0 := []float64{0, 0.999, 0}
	textt0v1 := []float64{0, 0, 0}
	textt0v2 := []float64{0.999, 0, 0}
	textt1v2 := []float64{0.999, 0, 0}
	textt1v3 := []float64{0.999, 0.999, 0}

	triangle1 := newTriangle()
	triangle1.worldVerts = t0
	triangle1.normals = []Vector3{normal, normal, normal}
	triangle1.uvMapping = [][]float64{textt0v1, textt0v2, textt0v0}

	triangle2 := newTriangle()
	triangle2.worldVerts = t1
	triangle2.normals = []Vector3{normal, normal, normal}
	triangle2.uvMapping = [][]float64{textt1v3, textt1v0, textt1v2}

	return &Model{
		triangles: []*Triangle{triangle1, triangle2},
		shader:    shader,
	}
}

func newCube(size float64, shader Shader) *Model {
	triangles := []*Triangle{}

	triangles = append(triangles, newXZSquare(size, shader).rotateX(math.Pi).moveY(-size/2).triangles...) // bottom
	triangles = append(triangles, newXZSquare(size, shader).moveY(size/2).triangles...)                   // top

	triangles = append(triangles, newYZSquare(size, shader).rotateY(math.Pi).moveX(size/2).triangles...) // right
	triangles = append(triangles, newYZSquare(size, shader).moveX(-size/2).triangles...)                 // left

	triangles = append(triangles, newXYSquare(size, shader).rotateY(math.Pi).moveZ(-size/2).triangles...) // back
	triangles = append(triangles, newXYSquare(size, shader).moveZ(size/2).triangles...)                   // front

	return &Model{
		triangles: triangles,
		shader:    shader,
	}
}
