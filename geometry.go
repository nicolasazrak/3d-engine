package main

func newXZSquare(size float64, shader Shader) *Model {
	v0 := Vector3{x: -size / 2, y: 0, z: -size / 2}
	v1 := Vector3{x: size / 2, y: 0, z: -size / 2}
	v2 := Vector3{x: size / 2, y: 0, z: size / 2}
	v3 := Vector3{x: -size / 2, y: 0, z: size / 2}

	t0 := []Vector3{v1, v0, v2}
	t1 := []Vector3{v3, v2, v0}

	normal := Vector3{x: 0, y: 1, z: 0}

	text0 := []float64{0, 0, 0}
	text1 := []float64{0.999, 0, 0}
	text2 := []float64{0.999, 0.999, 0}
	text3 := []float64{0, 0.999, 0}

	triangle1 := newTriangle()
	triangle1.worldVerts = t0
	triangle1.normals = []Vector3{normal, normal, normal}
	triangle1.uvMapping = [][]float64{text1, text0, text2}

	triangle2 := newTriangle()
	triangle2.worldVerts = t1
	triangle2.normals = []Vector3{normal, normal, normal}
	triangle2.uvMapping = [][]float64{text3, text2, text0}

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

	text0 := []float64{0, 0.999, 0}
	text1 := []float64{0, 0, 0}
	text2 := []float64{0.999, 0, 0}
	text3 := []float64{0.999, 0.999, 0}

	triangle1 := newTriangle()
	triangle1.worldVerts = t0
	triangle1.normals = []Vector3{normal, normal, normal}
	triangle1.uvMapping = [][]float64{text1, text2, text0}

	triangle2 := newTriangle()
	triangle2.worldVerts = t1
	triangle2.normals = []Vector3{normal, normal, normal}
	triangle2.uvMapping = [][]float64{text3, text0, text2}

	return &Model{
		triangles: []*Triangle{triangle1, triangle2},
		shader:    shader,
	}
}
