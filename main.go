package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/tfriedel6/canvas/sdlcanvas"
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

func parseModel(path string) []Triangle {
	triangles := []Triangle{}
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	str := string(f)
	vertex := []Vector3{}
	for _, line := range strings.Split(str, "\n") {
		if strings.HasPrefix(line, "v ") {
			splitted := strings.Split(line, " ")
			x, err := strconv.ParseFloat(splitted[1], 64)
			if err != nil {
				print(line)
				panic(err)
			}
			y, err := strconv.ParseFloat(splitted[2], 64)
			if err != nil {
				panic(err)
			}
			z, err := strconv.ParseFloat(splitted[3], 64)
			if err != nil {
				panic(err)
			}
			vertex = append(vertex, Vector3{x, y, z})
		}

		if strings.HasPrefix(line, "f ") {
			splitted := strings.Split(line, " ")
			idx1, err := strconv.ParseInt(strings.Split(splitted[1], "/")[0], 10, 32)
			if err != nil {
				panic(err)
			}
			idx2, err := strconv.ParseInt(strings.Split(splitted[2], "/")[0], 10, 32)
			if err != nil {
				panic(err)
			}
			idx3, err := strconv.ParseInt(strings.Split(splitted[3], "/")[0], 10, 32)
			if err != nil {
				panic(err)
			}

			v1 := minus(vertex[idx3-1], vertex[idx1-1])
			v2 := minus(vertex[idx2-1], vertex[idx1-1])
			normal := normalize(crossProduct(v1, v2))

			triangles = append(triangles, Triangle{
				verts:  []Vector3{vertex[idx1-1], vertex[idx2-1], vertex[idx3-1]},
				color:  color.White,
				normal: normal,
			})
		}
	}

	return triangles
}

func drawLine(x0, y0, x1, y1 float64, image *image.RGBA, color color.Color) {
	steep := false
	if math.Abs(x0-x1) < math.Abs(y0-y1) {
		y0, x0 = x0, y0
		y1, x1 = x1, y1
		steep = true
	}
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	dx := x1 - x0
	dy := y1 - y0
	derror2 := math.Abs(dy) * 2
	error2 := float64(0)
	y := int(y0)
	for x := int(x0); x <= int(x1); x++ {
		if steep {
			image.Set(y, x, color)
		} else {
			image.Set(x, y, color)
		}
		error2 += derror2
		if error2 > dx {
			if y1 > y0 {
				y += 1
			} else {
				y -= 1
			}
			error2 -= dx * 2
		}
	}
}

func boundingBox(pts []Vector3) (Vector2, Vector2) {
	minx := float64(999)
	miny := float64(999)
	maxx := float64(-999)
	maxy := float64(-999)

	for _, p := range pts {
		minx = math.Min(minx, p.x)
		maxx = math.Max(maxx, p.x)

		miny = math.Min(miny, p.y)
		maxy = math.Max(maxy, p.y)
	}

	return Vector2{minx, miny}, Vector2{maxx, maxy}
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

func times(a Vector3, b Vector3) float64 {
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

func orient2d(a Vector3, b Vector3, x float64, y float64) float64 {
	return (b.x-a.x)*(y-a.y) - (b.y-a.y)*(x-a.x)
}

func drawTriangle(pts []Vector3, image *image.RGBA, color color.Color, width int, zbuffer []float64) {
	minbbox, maxbbox := boundingBox(pts)
	for x := minbbox.x; x <= maxbbox.x; x++ {
		intX := int(x)
		for y := minbbox.y; y <= maxbbox.y; y++ {
			coordinates := baycentricCoordinates(x, y, pts)
			if coordinates.x < 0 || coordinates.y < 0 || coordinates.z < 0 {
				continue
			}

			zPos := float64(0)
			for _, pt := range pts {
				zPos += pt.z * coordinates.z
			}
			idx := intX + int(y)*width
			if zPos < zbuffer[idx] {
				zbuffer[idx] = zPos
				image.Set(intX, int(y), color)
			}
		}
	}
}

func main() {
	// fmt.Println(crossProduct(Vector3{2, 3, 4}, Vector3{5, 6, 7}))
	main2()
}

func main2() {
	wnd, cv, err := sdlcanvas.CreateWindow(1000, 1000, "Hello")
	if err != nil {
		panic(err)
	}
	defer wnd.Destroy()

	model := parseModel("head.obj")

	wnd.MainLoop(func() {

		start := time.Now()
		width := cv.Width()
		height := cv.Height()
		fWidth := float64(width)
		fHeight := float64(height)
		image := cv.GetImageData(0, 0, width, height)

		// Clear buffer
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				image.Set(x, y, color.Black)
			}
		}

		zbuffer := make([]float64, width*height)
		for i := range zbuffer {
			zbuffer[i] = 999999
		}
		light := Vector3{0, 0, -1}

		for _, triangle := range model {
			verts := []Vector3{}
			for _, vert := range triangle.verts {
				verts = append(verts, Vector3{
					x: (vert.x + 1.) * fWidth / 2.,
					y: (vert.y + 1.) * fHeight / 2.,
					z: 0,
				})
			}

			intensity := times(triangle.normal, light) * .8
			if intensity > 0 {
				c := color.RGBA{uint8(intensity * 255), uint8(intensity * 255), uint8(intensity * 255), 255}
				drawTriangle(verts, image, c, width, zbuffer)
			}
		}

		flipped := imaging.FlipV(image)
		for x := 0; x < height; x++ {
			for y := 0; y < width; y++ {
				image.Set(x, y, flipped.At(x, y))
			}
		}

		cv.PutImageData(image, 0, 0)

		elapsed := time.Since(start)
		fmt.Println(elapsed.String())
	})
}
