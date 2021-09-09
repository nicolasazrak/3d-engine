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

		minD := 9999999.
		normal := Vector3{}
		collided := false

		if from.z < bb.pmin.z {
			divZPlane := dotProduct(direction, Vector3{0, 0, 1})
			if divZPlane > 0 {
				numerator := dotProduct(minus(bb.pmin, from), Vector3{0, 0, 1})
				d := numerator / divZPlane
				intersection := Vector3{
					x: from.x + d*direction.x,
					y: from.y + d*direction.y,
					z: from.z + d*direction.z,
				}

				if intersection.y > bb.pmin.y && intersection.y < bb.pmax.y &&
					intersection.x > bb.pmin.x && intersection.x < bb.pmax.x {

					if d < minD {
						minD = d
						normal = Vector3{0, 0, 1}
						collided = true
					}
				}
			}
		}

		if from.z > bb.pmax.z {
			divZPlane := dotProduct(direction, Vector3{0, 0, 1})
			if divZPlane < 0 {
				numerator := dotProduct(minus(bb.pmax, from), Vector3{0, 0, 1})
				d := numerator / divZPlane
				intersection := Vector3{
					x: from.x + d*direction.x,
					y: from.y + d*direction.y,
					z: from.z + d*direction.z,
				}

				if intersection.y > bb.pmin.y && intersection.y < bb.pmax.y &&
					intersection.x > bb.pmin.x && intersection.x < bb.pmax.x {
					if d < minD {
						minD = d
						normal = Vector3{0, 0, -1}
						collided = true
					}
				}
			}
		}

		if from.x < bb.pmin.x {
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
					if d < minD {
						minD = d
						normal = Vector3{-1, 0, 0}
						collided = true
					}
				}
			}
		}

		if from.x > bb.pmin.x {
			divXPlane := dotProduct(direction, Vector3{1, 0, 0})
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
					if d < minD {
						minD = d
						normal = Vector3{1, 0, 0}
						collided = true
					}
				}
			}
		}

		if from.y < bb.pmin.y {
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
					if d < minD {
						minD = d
						normal = Vector3{0, -1, 0}
						collided = true
					}
				}
			}
		}

		if from.y > bb.pmax.y {
			divYPlane := dotProduct(direction, Vector3{0, 1, 0})
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
					if d < minD {
						minD = d
						normal = Vector3{0, 1, 0}
						collided = true
					}
				}
			}
		}
		return collided, normal, minD
	} else {
		return false, from, 0
	}
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
