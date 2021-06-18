package main

import (
	"math"
)

type Vector4 struct {
	x float64
	y float64
	z float64
	w float64
}

type Vector3 struct {
	x float64
	y float64
	z float64
}

type Vector2 struct {
	x float64
	y float64
}

func boundingBox(pts []Vector2, minx float64, maxx float64, miny float64, maxy float64) (Vector2, Vector2) {
	min := Vector2{
		x: math.Max(minx, math.Min(pts[0].x, math.Min(pts[1].x, pts[2].x))),
		y: math.Max(miny, math.Min(pts[0].y, math.Min(pts[1].y, pts[2].y))),
	}
	max := Vector2{
		x: math.Min(maxx, math.Max(pts[0].x, math.Max(pts[1].x, pts[2].x))),
		y: math.Min(maxy, math.Max(pts[0].y, math.Max(pts[1].y, pts[2].y))),
	}
	return min, max
}

func crossProduct(a Vector3, b Vector3) Vector3 {
	cx := a.y*b.z - a.z*b.y
	cy := a.z*b.x - a.x*b.z
	cz := a.x*b.y - a.y*b.x
	return Vector3{cx, cy, cz}
}

func norm(a Vector3) float64 {
	return math.Sqrt(a.x*a.x + a.y*a.y + a.z*a.z)
}

func normalize(a Vector3) Vector3 {
	n := norm(a)
	return Vector3{
		x: a.x / n,
		y: a.y / n,
		z: a.z / n,
	}
}

func dotProduct(a Vector3, b Vector3) float64 {
	return a.x*b.x + a.y*b.y + a.z*b.z
}

func minus(a Vector3, b Vector3) Vector3 {
	return Vector3{x: a.x - b.x, y: a.y - b.y, z: a.z - b.z}
}

func baycentricCoordinates(x float64, y float64, triangle []Vector3) Vector3 {
	v0 := Vector3{triangle[2].x - triangle[0].x, triangle[1].x - triangle[0].x, triangle[0].x - x}
	v1 := Vector3{triangle[2].y - triangle[0].y, triangle[1].y - triangle[0].y, triangle[0].y - y}
	u := crossProduct(v0, v1)

	if math.Abs(u.z) < 1 {
		return Vector3{-1, -1, -1}
	}

	return Vector3{1. - (u.x+u.y)/u.z, u.y / u.z, u.x / u.z}
}

func orient2d(a Vector2, b Vector2, x float64, y float64) int {
	return int((b.x-a.x)*(y-a.y) - (b.y-a.y)*(x-a.x))
}

func ponderate(pts []Vector3, weights []float64) Vector3 {
	return Vector3{
		x: pts[0].x*weights[0] + pts[1].x*weights[1] + pts[2].x*weights[2],
		y: pts[0].y*weights[0] + pts[1].y*weights[1] + pts[2].y*weights[2],
		z: pts[0].z*weights[0] + pts[1].z*weights[1] + pts[2].z*weights[2],
	}
}

func matmult(m [4][4]float64, vec Vector3) Vector3 {
	x := m[0][0]*vec.x + m[1][0]*vec.y + m[2][0]*vec.z + m[3][0]
	y := m[0][1]*vec.x + m[1][1]*vec.y + m[2][1]*vec.z + m[3][1]
	z := m[0][2]*vec.x + m[1][2]*vec.y + m[2][2]*vec.z + m[3][2]
	w := m[0][3]*vec.x + m[1][3]*vec.y + m[2][3]*vec.z + m[3][3]

	return Vector3{x / w, y / w, z / w}
}
