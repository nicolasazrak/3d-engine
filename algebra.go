package main

import (
	"image/color"
	"math"
)

type Vector3 struct {
	x float64
	y float64
	z float64
}

type Vector2 struct {
	x float64
	y float64
}

type Triangle struct {
	verts  []Vector3
	color  color.Color
	normal Vector3
}

func boundingBox(pts []Vector3) (Vector2, Vector2) {
	min := Vector2{x: math.Min(pts[0].x, math.Min(pts[1].x, pts[2].x)), y: math.Min(pts[0].y, math.Min(pts[1].y, pts[2].y))}
	max := Vector2{x: math.Max(pts[0].x, math.Max(pts[1].x, pts[2].x)), y: math.Max(pts[0].y, math.Max(pts[1].y, pts[2].y))}
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

func orient2d(a Vector3, b Vector3, x float64, y float64) int {
	return int((b.x-a.x)*(y-a.y) - (b.y-a.y)*(x-a.x))
}
