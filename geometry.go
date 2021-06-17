package main

func newXZSquare(y float64, size float64, shader Shader) *Model {
	v0 := Vector3{x: 0, y: y, z: 0}
	v1 := Vector3{x: size, y: y, z: 0}
	v2 := Vector3{x: size, y: y, z: size}
	v3 := Vector3{x: 0, y: y, z: size}

	t0 := []Vector3{v0, v2, v1}
	t1 := []Vector3{v2, v0, v3}

	normal := Vector3{x: 0, y: 1, z: 0}

	text0 := Vector3{x: 0, y: 0, z: 0}
	text1 := Vector3{x: 0.999, y: 0, z: 0}
	text2 := Vector3{x: 0.999, y: 0.999, z: 0}
	text3 := Vector3{x: 0, y: 0.999, z: 0}

	triangle1 := Triangle{
		verts:            t0,
		normals:          []Vector3{normal, normal, normal},
		screenProjection: []Vector2{{}, {}, {}},
		uvMapping:        []Vector3{text0, text2, text1},
		viewProjection:   []Vector3{{}, {}, {}},
	}

	triangle2 := Triangle{
		verts:            t1,
		normals:          []Vector3{normal, normal, normal},
		screenProjection: []Vector2{{}, {}, {}},
		uvMapping:        []Vector3{text2, text0, text3},
		viewProjection:   []Vector3{{}, {}, {}},
	}

	return &Model{
		triangles: []*Triangle{&triangle1, &triangle2},
		shader:    shader,
	}
}

func newXYSquare(z float64, size float64, shader Shader) *Model {
	v0 := Vector3{x: 0, y: 0, z: z}
	v1 := Vector3{x: size, y: 0, z: z}
	v2 := Vector3{x: size, y: size, z: z}
	v3 := Vector3{x: 0, y: size, z: z}

	t0 := []Vector3{v0, v1, v2}
	t1 := []Vector3{v2, v3, v0}

	normal := Vector3{x: 0, y: 0, z: -1}

	text0 := Vector3{x: 0, y: 0, z: 0}
	text1 := Vector3{x: 0.999, y: 0, z: 0}
	text2 := Vector3{x: 0.999, y: 0.999, z: 0}
	text3 := Vector3{x: 0, y: 0.999, z: 0}

	triangle1 := Triangle{
		verts:            t0,
		normals:          []Vector3{normal, normal, normal},
		screenProjection: []Vector2{{}, {}, {}},
		uvMapping:        []Vector3{text0, text1, text2},
		viewProjection:   []Vector3{{}, {}, {}},
	}

	triangle2 := Triangle{
		verts:            t1,
		normals:          []Vector3{normal, normal, normal},
		screenProjection: []Vector2{{}, {}, {}},
		uvMapping:        []Vector3{text2, text3, text0},
		viewProjection:   []Vector3{{}, {}, {}},
	}

	return &Model{
		triangles: []*Triangle{&triangle1, &triangle2},
		shader:    shader,
	}
}
