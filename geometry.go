package main

import "fmt"

func newXZSquare(size float64, shader Shader) *Model {
	v0 := Vector3{x: -size / 2, y: 0, z: -size / 2}
	v1 := Vector3{x: size / 2, y: 0, z: -size / 2}
	v2 := Vector3{x: size / 2, y: 0, z: size / 2}
	v3 := Vector3{x: -size / 2, y: 0, z: size / 2}

	t0 := []Vector3{v0, v2, v1}
	t1 := []Vector3{v2, v0, v3}

	normal := Vector3{x: 0, y: 1, z: 0}

	text0 := Vector3{x: 0, y: 0, z: 0}
	text1 := Vector3{x: 0.999, y: 0, z: 0}
	text2 := Vector3{x: 0.999, y: 0.999, z: 0}
	text3 := Vector3{x: 0, y: 0.999, z: 0}

	triangle1 := Triangle{
		verts:         t0,
		normals:       []Vector3{normal, normal, normal},
		viewportVerts: []Vector2{{}, {}, {}},
		uvMapping:     []Vector3{text0, text2, text1},
		cameraVerts:   []Vector3{{}, {}, {}},
		cameraNormals: []Vector3{{}, {}, {}},
	}

	triangle2 := Triangle{
		verts:         t1,
		normals:       []Vector3{normal, normal, normal},
		viewportVerts: []Vector2{{}, {}, {}},
		uvMapping:     []Vector3{text2, text0, text3},
		cameraVerts:   []Vector3{{}, {}, {}},
		cameraNormals: []Vector3{{}, {}, {}},
	}

	if false {
		fmt.Println(triangle1)
	}
	return &Model{
		triangles: []*Triangle{&triangle1, &triangle2},
		// triangles: []*Triangle{&triangle2},
		shader: shader,
	}
}

func newXYSquare(size float64, shader Shader) *Model {
	v0 := Vector3{x: -size / 2, y: -size / 2, z: 0}
	v1 := Vector3{x: size / 2, y: -size / 2, z: 0}
	v2 := Vector3{x: size / 2, y: size / 2, z: 0}
	v3 := Vector3{x: -size / 2, y: size / 2, z: 0}

	t0 := []Vector3{v0, v1, v2}
	t1 := []Vector3{v3, v0, v2}

	normal := Vector3{x: 0, y: 0, z: 1}

	text0 := Vector3{x: 0, y: 0, z: 0}
	text1 := Vector3{x: 0.999, y: 0, z: 0}
	text2 := Vector3{x: 0.999, y: 0.999, z: 0}
	text3 := Vector3{x: 0, y: 0.999, z: 0}

	triangle1 := Triangle{
		verts:         t0,
		normals:       []Vector3{normal, normal, normal},
		viewportVerts: []Vector2{{}, {}, {}},
		uvMapping:     []Vector3{text0, text1, text2},
		cameraVerts:   []Vector3{{}, {}, {}},
		cameraNormals: []Vector3{{}, {}, {}},
	}

	triangle2 := Triangle{
		verts:         t1,
		normals:       []Vector3{normal, normal, normal},
		viewportVerts: []Vector2{{}, {}, {}},
		uvMapping:     []Vector3{text3, text0, text2},
		cameraVerts:   []Vector3{{}, {}, {}},
		cameraNormals: []Vector3{{}, {}, {}},
	}

	return &Model{
		triangles: []*Triangle{&triangle1, &triangle2},
		shader:    shader,
	}
}
