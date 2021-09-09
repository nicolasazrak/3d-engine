package main

import (
	"math"
)

type Collisionable interface {
	test(from Vector3, to Vector3, direction Vector3) (bool, Vector3, float64)
}

type SquaredBoundingBox struct {
	pmin Vector3
	pmax Vector3
}

func (bb *SquaredBoundingBox) test(from Vector3, to Vector3, direction Vector3) (bool, Vector3, float64) {
	if to.x > bb.pmin.x && to.x < bb.pmax.x &&
		to.y > bb.pmin.y && to.y < bb.pmax.y &&
		to.z > bb.pmin.z && to.z < bb.pmax.z {

		divXPlane := dotProduct(direction, Vector3{1, 0, 0})
		if divXPlane > 0 {
			numerator := dotProduct(minus(bb.pmin, from), Vector3{1, 0, 0})
			d := numerator / divXPlane
			intersection := Vector3{
				x: from.x + d*direction.x,
				y: from.y + d*direction.y,
				z: from.z + d*direction.z,
			}

			if intersection.y > bb.pmin.y && intersection.y < bb.pmax.y &&
				intersection.z > bb.pmin.z && intersection.z < bb.pmax.z {
				return true, Vector3{-1, 0, 0}, d
			}
		}

		if divXPlane < 0 {
			numerator := dotProduct(minus(bb.pmax, from), Vector3{1, 0, 0})
			d := numerator / divXPlane
			intersection := Vector3{
				x: from.x + d*direction.x,
				y: from.y + d*direction.y,
				z: from.z + d*direction.z,
			}

			if intersection.y > bb.pmin.y && intersection.y < bb.pmax.y &&
				intersection.z > bb.pmin.z && intersection.z < bb.pmax.z &&
				divXPlane < 0 {
				return true, Vector3{1, 0, 0}, d
			}
		}

		divYPlane := dotProduct(direction, Vector3{0, 1, 0})
		if divYPlane > 0 {
			numerator := dotProduct(minus(bb.pmax, from), Vector3{0, 1, 0})
			d := numerator / divYPlane
			intersection := Vector3{
				x: from.x + d*direction.x,
				y: from.y + d*direction.y,
				z: from.z + d*direction.z,
			}

			if intersection.x > bb.pmin.x && intersection.x < bb.pmax.x &&
				intersection.z > bb.pmin.z && intersection.z < bb.pmax.z {
				return true, Vector3{0, -1, 0}, d
			}
		}

		if divYPlane < 0 {
			numerator := dotProduct(minus(bb.pmin, from), Vector3{0, 1, 0})
			d := numerator / divYPlane
			intersection := Vector3{
				x: from.x + d*direction.x,
				y: from.y + d*direction.y,
				z: from.z + d*direction.z,
			}

			if intersection.x > bb.pmin.x && intersection.x < bb.pmax.x &&
				intersection.z > bb.pmin.z && intersection.z < bb.pmax.z {
				return true, Vector3{0, 1, 0}, d
			}
		}

		divZPlane := dotProduct(direction, Vector3{0, 0, 1})
		if divZPlane > 0 {
			numerator := dotProduct(minus(bb.pmax, from), Vector3{0, 0, 1})
			d := numerator / divZPlane
			intersection := Vector3{
				x: from.x + d*direction.x,
				y: from.y + d*direction.y,
				z: from.z + d*direction.z,
			}

			if intersection.y > bb.pmin.y && intersection.y < bb.pmax.y &&
				intersection.x > bb.pmin.x && intersection.x < bb.pmax.x {
				return true, Vector3{0, 0, -1}, d
			}
		}

		if divZPlane < 0 {
			numerator := dotProduct(minus(bb.pmin, from), Vector3{0, 0, 1})
			d := numerator / divZPlane
			intersection := Vector3{
				x: from.x + d*direction.x,
				y: from.y + d*direction.y,
				z: from.z + d*direction.z,
			}

			if intersection.y > bb.pmin.y && intersection.y < bb.pmax.y &&
				intersection.x > bb.pmin.x && intersection.x < bb.pmax.x {
				return true, Vector3{0, 0, 1}, d
			}
		}

		return false, from, 0
	} else {
		return false, from, 0
	}
}

type TriangleCollisionable struct {
	triangle *Triangle
}

func (t *TriangleCollisionable) test(from Vector3, to Vector3, direction Vector3) (bool, Vector3) {
	EPSILON := 0.0000001
	vertex0 := t.triangle.worldVerts[0]
	vertex1 := t.triangle.worldVerts[1]
	vertex2 := t.triangle.worldVerts[2]
	edge1 := minus(vertex1, vertex0)
	edge2 := minus(vertex2, vertex0)
	h := crossProduct(direction, edge2)
	a := dotProduct(edge1, h)
	if a > -EPSILON && a < EPSILON {
		return false, from
	}

	f := 1.0 / a
	s := minus(direction, vertex0)
	u := f * dotProduct(s, h)
	if u < 0.0 || u > 1.0 {
		return false, from
	}

	q := crossProduct(s, edge1)
	v := f * dotProduct(direction, q)
	if v < 0.0 || u+v > 1.0 {
		return false, from
	}
	// At this stage we can compute t to find out where the intersection point is on the line.
	ttt := f * dotProduct(edge2, q)
	if ttt > EPSILON {
		return true, t.triangle.normals[0]
	}

	// This means that there is a line intersection but not a ray intersection.
	return false, from
}

func newCollisionableFromModel(model *Model) Collisionable {
	pmin := Vector3{x: 999999999, y: 9999999999, z: 99999999}
	pmax := Vector3{x: -999999999, y: -9999999999, z: -99999999}

	for _, triangle := range model.triangles {
		for _, vert := range triangle.worldVerts {
			pmax.x = math.Max(pmax.x, vert.x) + 0.01
			pmax.y = math.Max(pmax.y, vert.y) + 0.01
			pmax.z = math.Max(pmax.z, vert.z) + 0.01

			pmin.x = math.Min(pmin.x, vert.x) - 0.01
			pmin.y = math.Min(pmin.y, vert.y) - 0.01
			pmin.z = math.Min(pmin.z, vert.z) - 0.01
		}
	}

	return &SquaredBoundingBox{
		pmin: pmin,
		pmax: pmax,
	}
}
