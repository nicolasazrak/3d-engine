package main

import "testing"

func projectedTriangle(vec []Vector4) *ProjectedTriangle {
	t := newProjectedTriangle()
	t.clipVertex = vec
	return t
}

func TestClip0(t *testing.T) {
	vert0 := Vector4{x: 0, y: 1, z: 0, w: 1}
	vert1 := Vector4{x: 1, y: 0, z: 0, w: 1}
	vert2 := Vector4{x: 0, y: -1, z: 0, w: 1}

	triangle := projectedTriangle([]Vector4{vert0, vert1, vert2})

	projection := clipTriangle(triangle)

	if len(projection) != 1 {
		t.Fail()
	}
}

func TestClip1(t *testing.T) {
	vert0 := Vector4{x: 0, y: 4, z: 0, w: 1}
	vert1 := Vector4{x: 0, y: 0, z: 0, w: 1}
	vert2 := Vector4{x: 1, y: 0, z: 0, w: 1}

	triangle := projectedTriangle([]Vector4{vert0, vert1, vert2})

	projection := clipTriangle(triangle)

	if len(projection) != 2 {
		t.Fail()
	}
}

func TestClip2(t *testing.T) {
	vert0 := Vector4{x: 0, y: 0, z: 0, w: 1}
	vert1 := Vector4{x: 0, y: 2, z: 0, w: 1}
	vert2 := Vector4{x: -2, y: 0, z: 0, w: 1}

	triangle := projectedTriangle([]Vector4{vert0, vert1, vert2})

	projection := clipTriangle(triangle)

	if len(projection) != 2 {
		t.Fail()
	}
}
