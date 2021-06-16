package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/tfriedel6/canvas/sdlcanvas"
)

type Scene struct {
	models         []Triangle
	zBuffer        []float64
	colorBuffer    []color.Color
	width          int
	height         int
	fWidth         float64
	fHeight        float64
	lightPosition  Vector3
	cameraPosition Vector3
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

func (scene *Scene) drawLine(x0, y0, x1, y1 float64, image *image.RGBA, color color.Color) {
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

func (scene *Scene) drawTriangle(pts []Vector3, color color.Color) {
	minbbox, maxbbox := boundingBox(pts)
	A01 := int(pts[0].y - pts[1].y)
	B01 := int(pts[1].x - pts[0].x)
	A12 := int(pts[1].y - pts[2].y)
	B12 := int(pts[2].x - pts[1].x)
	A20 := int(pts[2].y - pts[0].y)
	B20 := int(pts[0].x - pts[2].x)

	w1_row := orient2d(pts[2], pts[0], minbbox.x, minbbox.y)
	w0_row := orient2d(pts[1], pts[2], minbbox.x, minbbox.y)
	w2_row := orient2d(pts[0], pts[1], minbbox.x, minbbox.y)

	for y := minbbox.y; y <= maxbbox.y; y++ {
		w0 := w0_row
		w1 := w1_row
		w2 := w2_row

		for x := minbbox.x; x <= maxbbox.x; x++ {
			if (w0 | w1 | w2) >= 0 {
				zPos := float64(w0)*pts[0].z + float64(w1)*pts[1].z + float64(w2)*pts[2].z
				idx := int(x) + (int(y))*scene.width
				if zPos < scene.zBuffer[idx] {
					scene.zBuffer[idx] = zPos
					scene.colorBuffer[idx] = color
				}
			}

			w0 += A12
			w1 += A20
			w2 += A01
		}

		// One row step
		w0_row += B12
		w1_row += B20
		w2_row += B01
	}
}

func (scene *Scene) drawModel() {
	for _, triangle := range scene.models {
		verts := []Vector3{}
		for _, vert := range triangle.verts {
			verts = append(verts, Vector3{
				x: (vert.x + 1.) * scene.fWidth / 2.,
				y: (vert.y + 1.) * scene.fHeight / 2.,
				z: 0,
			})
		}

		lightNormal := normalize(minus(scene.lightPosition, triangle.verts[0]))
		intensity := dotProduct(lightNormal, triangle.normal) * .8
		c := color.RGBA{uint8(intensity * 255), uint8(intensity * 255), uint8(intensity * 255), 255}
		if intensity < 0 {
			// Shoudln't be needed if there was occulsion culling or shadows ?
			c = color.RGBA{0, 0, 0, 0}
		}

		scene.drawTriangle(verts, c)
	}
}

func (scene *Scene) toImage() *image.RGBA {
	pix := make([]uint8, scene.width*scene.height*4)
	for x := 0; x < scene.width; x++ {
		for yInverted := 0; yInverted < scene.height; yInverted++ {
			y := scene.height - yInverted - 1
			colorIdx := x + y*scene.width
			pixIdx := (x + yInverted*scene.width) * 4
			r, g, b, _ := scene.colorBuffer[colorIdx].RGBA()
			pix[pixIdx] = uint8(r)
			pix[pixIdx+1] = uint8(g)
			pix[pixIdx+2] = uint8(b)
			pix[pixIdx+3] = uint8(255)
		}
	}

	image := &image.RGBA{
		Pix:    pix,
		Stride: scene.width * 4,
		Rect:   image.Rect(0, 0, scene.width, scene.height),
	}

	return image
}

func (scene *Scene) cleanBuffer() {
	for x := 0; x < scene.width; x++ {
		for y := 0; y < scene.height; y++ {
			scene.colorBuffer[x+y*scene.width] = color.Black
		}
	}

	for i := range scene.zBuffer {
		scene.zBuffer[i] = 999999
	}
}

func (scene *Scene) render() *image.RGBA {
	scene.cleanBuffer()
	scene.drawModel()
	return scene.toImage()
}

func main() {
	wnd, cv, err := sdlcanvas.CreateWindow(1000, 1000, "Hello")
	if err != nil {
		panic(err)
	}
	defer wnd.Destroy()

	scene := Scene{
		models:         []Triangle{},
		zBuffer:        []float64{},
		colorBuffer:    []color.Color{},
		width:          cv.Width(),
		height:         cv.Height(),
		fWidth:         float64(cv.Width()),
		fHeight:        float64(cv.Height()),
		lightPosition:  Vector3{0, 0, -10},
		cameraPosition: Vector3{0, 0, 0},
	}

	scene.models = append(scene.models, parseModel("head.obj")...)
	scene.colorBuffer = make([]color.Color, scene.width*scene.height)
	scene.zBuffer = make([]float64, scene.width*scene.height)

	f, err := os.Create("cpu")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	t := float64(0)
	wnd.MainLoop(func() {
		start := time.Now()

		t++
		scene.lightPosition.z = math.Cos(t/10) * 3
		scene.lightPosition.x = math.Sin(t/10) * 3

		cv.PutImageData(scene.render(), 0, 0)
		elapsed := time.Since(start)
		fmt.Println(elapsed.String())
	})
}
