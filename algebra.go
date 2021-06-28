package main

import (
	"math"
	"unsafe"
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
	ptsminx := math.Min(pts[0].x, math.Min(pts[1].x, pts[2].x))
	ptsmaxx := math.Max(pts[0].x, math.Max(pts[1].x, pts[2].x))

	ptsminy := math.Min(pts[0].y, math.Min(pts[1].y, pts[2].y))
	ptsmaxy := math.Max(pts[0].y, math.Max(pts[1].y, pts[2].y))

	min := Vector2{
		x: math.Max(minx, math.Min(ptsminx, maxx)),
		y: math.Max(miny, math.Min(ptsminy, maxy)),
	}
	max := Vector2{
		x: math.Min(maxx, math.Max(ptsmaxx, minx)),
		y: math.Min(maxy, math.Max(ptsmaxy, miny)),
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

func fastInvSqrt(v float64) float64 {
	x := float32(v)
	xhalf := float32(0.5) * x
	i := *(*int32)(unsafe.Pointer(&x))
	i = int32(0x5f3759df) - int32(i>>1)
	x = *(*float32)(unsafe.Pointer(&i))
	x = x * (1.5 - (xhalf * x * x))
	return float64(x)
}

func fastNormalize(a Vector3) Vector3 {
	n := fastInvSqrt(a.x*a.x + a.y*a.y + a.z*a.z)
	return Vector3{
		x: a.x * n,
		y: a.y * n,
		z: a.z * n,
	}
}

func normalize(a Vector3) Vector3 {
	n := 1 / norm(a)
	return Vector3{
		x: a.x * n,
		y: a.y * n,
		z: a.z * n,
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

	if u.z > -1 && u.z < 1 {
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

func matmult(m [4][4]float64, vec Vector3, h float64) Vector3 {
	x := m[0][0]*vec.x + m[1][0]*vec.y + m[2][0]*vec.z + m[3][0]*h
	y := m[0][1]*vec.x + m[1][1]*vec.y + m[2][1]*vec.z + m[3][1]*h
	z := m[0][2]*vec.x + m[1][2]*vec.y + m[2][2]*vec.z + m[3][2]*h
	w := m[0][3]*vec.x + m[1][3]*vec.y + m[2][3]*vec.z + m[3][3]*h
	div := 1 / w
	return Vector3{x * div, y * div, z * div}
}

func inverseTranspose(dst *[4][4]float64, src [4][4]float64) {
	// https://semath.info/src/inverse-cofactor-ex4.html
	// https://stackoverflow.com/questions/33088577/symbolically-calculate-the-inverse-of-a-4-x-4-matrix-in-matlab
	A11 := src[0][0]
	A12 := src[1][0]
	A13 := src[2][0]
	A14 := src[3][0]

	A21 := src[0][1]
	A22 := src[1][1]
	A23 := src[2][1]
	A24 := src[3][1]

	A31 := src[0][2]
	A32 := src[1][2]
	A33 := src[2][2]
	A34 := src[3][2]

	A41 := src[0][3]
	A42 := src[1][3]
	A43 := src[2][3]
	A44 := src[3][3]

	dst[0][0] = (A22*A33*A44 - A22*A34*A43 - A23*A32*A44 + A23*A34*A42 + A24*A32*A43 - A24*A33*A42) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)
	dst[0][1] = -(A12*A33*A44 - A12*A34*A43 - A13*A32*A44 + A13*A34*A42 + A14*A32*A43 - A14*A33*A42) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)
	dst[0][2] = (A12*A23*A44 - A12*A24*A43 - A13*A22*A44 + A13*A24*A42 + A14*A22*A43 - A14*A23*A42) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)
	dst[0][3] = -(A12*A23*A34 - A12*A24*A33 - A13*A22*A34 + A13*A24*A32 + A14*A22*A33 - A14*A23*A32) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)

	dst[1][0] = -(A21*A33*A44 - A21*A34*A43 - A23*A31*A44 + A23*A34*A41 + A24*A31*A43 - A24*A33*A41) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)
	dst[1][1] = (A11*A33*A44 - A11*A34*A43 - A13*A31*A44 + A13*A34*A41 + A14*A31*A43 - A14*A33*A41) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)
	dst[1][2] = -(A11*A23*A44 - A11*A24*A43 - A13*A21*A44 + A13*A24*A41 + A14*A21*A43 - A14*A23*A41) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)
	dst[1][3] = (A11*A23*A34 - A11*A24*A33 - A13*A21*A34 + A13*A24*A31 + A14*A21*A33 - A14*A23*A31) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)

	dst[2][0] = (A21*A32*A44 - A21*A34*A42 - A22*A31*A44 + A22*A34*A41 + A24*A31*A42 - A24*A32*A41) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)
	dst[2][1] = -(A11*A32*A44 - A11*A34*A42 - A12*A31*A44 + A12*A34*A41 + A14*A31*A42 - A14*A32*A41) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)
	dst[2][2] = (A11*A22*A44 - A11*A24*A42 - A12*A21*A44 + A12*A24*A41 + A14*A21*A42 - A14*A22*A41) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)
	dst[2][3] = -(A11*A22*A34 - A11*A24*A32 - A12*A21*A34 + A12*A24*A31 + A14*A21*A32 - A14*A22*A31) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)

	dst[3][0] = -(A21*A32*A43 - A21*A33*A42 - A22*A31*A43 + A22*A33*A41 + A23*A31*A42 - A23*A32*A41) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)
	dst[3][1] = (A11*A32*A43 - A11*A33*A42 - A12*A31*A43 + A12*A33*A41 + A13*A31*A42 - A13*A32*A41) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)
	dst[3][2] = -(A11*A22*A43 - A11*A23*A42 - A12*A21*A43 + A12*A23*A41 + A13*A21*A42 - A13*A22*A41) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)
	dst[3][3] = (A11*A22*A33 - A11*A23*A32 - A12*A21*A33 + A12*A23*A31 + A13*A21*A32 - A13*A22*A31) / (A11*A22*A33*A44 - A11*A22*A34*A43 - A11*A23*A32*A44 + A11*A23*A34*A42 + A11*A24*A32*A43 - A11*A24*A33*A42 - A12*A21*A33*A44 + A12*A21*A34*A43 + A12*A23*A31*A44 - A12*A23*A34*A41 - A12*A24*A31*A43 + A12*A24*A33*A41 + A13*A21*A32*A44 - A13*A21*A34*A42 - A13*A22*A31*A44 + A13*A22*A34*A41 + A13*A24*A31*A42 - A13*A24*A32*A41 - A14*A21*A32*A43 + A14*A21*A33*A42 + A14*A22*A31*A43 - A14*A22*A33*A41 - A14*A23*A31*A42 + A14*A23*A32*A41)
}
