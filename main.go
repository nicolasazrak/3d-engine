package main

import (
	"fmt"
	"image"
	"os"
	"runtime/pprof"
	"time"

	"github.com/tfriedel6/canvas/sdlcanvas"
)

type Scene struct {
	models         []*Model
	scaleFactor    int
	zBuffer        []float64
	pixBuffer      []uint8
	cleanPixBuffer []uint8
	cleanZBuffer   []float64
	width          int
	height         int
	fWidth         float64
	fHeight        float64
	lightPosition  Vector3
	projectedLight Vector3
	camera         *Camera
	trianglesDrawn int
}

func (scene *Scene) drawTriangle(model *Model, triangle *Triangle) {
	// if triangle.cameraVerts[0].z > 0 && triangle.cameraVerts[1].z > 0 && triangle.cameraVerts[2].z > 0 {
	// 	return
	// }

	farPlane := 20.
	nearPlane := .1
	scene.trianglesDrawn++
	pts := triangle.viewportVerts
	minbbox, maxbbox := boundingBox(pts, 0, scene.fWidth-1, 0, scene.fHeight-1)

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
				sum := float64(w0 + w1 + w2)
				if sum <= -0 {
					// TODO fix this
					continue
				}
				l0 := float64(w0) / sum
				l1 := float64(w1) / sum
				l2 := float64(w2) / sum

				zPos := l0*triangle.cameraVerts[0].z + l1*triangle.cameraVerts[1].z + l2*triangle.cameraVerts[2].z
				zPos = -zPos // WTF ??
				idx := int(x) + (int(y))*scene.width

				if zPos < farPlane && zPos > nearPlane && zPos < scene.zBuffer[idx] {
					scene.zBuffer[idx] = zPos
					r, g, b := model.shader.shade(scene, triangle, l0, l1, l2)
					scene.setAt(int(x), int(y), r, g, b)
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

func (scene *Scene) drawModel(model *Model) {
	for _, triangle := range model.triangles {
		scene.drawTriangle(model, triangle)
	}
}

func (scene *Scene) drawModels() {
	for _, model := range scene.models {
		scene.drawModel(model)
	}
}

func (scene *Scene) setAt(x int, yInverted int, r uint8, g uint8, b uint8) {
	for xScale := 0; xScale < scene.scaleFactor; xScale++ {
		for yScale := 0; yScale < scene.scaleFactor; yScale++ {
			y := scene.height - yInverted - 1

			finalX := x*scene.scaleFactor + xScale
			finalY := y*scene.scaleFactor + yScale

			pixIdx := (finalX + finalY*scene.width*scene.scaleFactor) * 4
			scene.pixBuffer[pixIdx] = r
			scene.pixBuffer[pixIdx+1] = g
			scene.pixBuffer[pixIdx+2] = b
			scene.pixBuffer[pixIdx+3] = uint8(255)
		}
	}
}

func (scene *Scene) toImage() *image.RGBA {
	image := &image.RGBA{
		Pix:    scene.pixBuffer,
		Stride: scene.width * 4 * scene.scaleFactor,
		Rect:   image.Rect(0, 0, scene.width*scene.scaleFactor, scene.height*scene.scaleFactor),
	}

	return image
}

func (scene *Scene) cleanBuffer() {
	scene.trianglesDrawn = 0
	copy(scene.pixBuffer[:], scene.cleanPixBuffer[:])
	copy(scene.zBuffer[:], scene.cleanZBuffer[:])
}

func (scene *Scene) render() *image.RGBA {
	scene.cleanBuffer()
	scene.drawModels()
	return scene.toImage()
}

func newScene(width int, height int, scaleFactor int) *Scene {
	scene := Scene{
		models:        []*Model{},
		zBuffer:       []float64{},
		pixBuffer:     []uint8{},
		scaleFactor:   scaleFactor,
		width:         width / scaleFactor,
		height:        height / scaleFactor,
		fWidth:        float64(width / scaleFactor),
		fHeight:       float64(height / scaleFactor),
		lightPosition: Vector3{2, 2, 1.5},
		camera:        newCamera(),
	}

	scene.pixBuffer = make([]uint8, scene.width*scene.height*4*scene.scaleFactor*scaleFactor)
	scene.cleanPixBuffer = make([]uint8, scene.width*scene.height*4*scene.scaleFactor*scaleFactor)
	scene.zBuffer = make([]float64, scene.width*scene.height)
	scene.cleanZBuffer = make([]float64, scene.width*scene.height)
	for idx := 0; idx < len(scene.cleanPixBuffer); idx += 4 {
		scene.cleanPixBuffer[idx] = uint8(0)
		scene.cleanPixBuffer[idx+1] = uint8(0)
		scene.cleanPixBuffer[idx+2] = uint8(0)
		scene.cleanPixBuffer[idx+3] = uint8(255)
	}
	for i := range scene.zBuffer {
		scene.cleanZBuffer[i] = 999999
	}

	return &scene
}

func main() {
	main2()
	// scene := newScene(100, 100)
	// scene.camera.position = Vector3{2, 0, 0}
	// scene.camera.target = Vector3{0, 0, 0}
	// scene.camera.updateViewMatrix()

	// // fmt.Println(matmult(scene.camera.viewMatrix, Vector3{0, 0, 0}, 1))
	// fmt.Println(normalize(matmult(scene.camera.normalMatrix, normalize(Vector3{0, 0, 1}), 0))) // should be {0,0,1}
}

func main2() {
	wnd, cv, err := sdlcanvas.CreateWindow(1000, 1000, "Hello")
	if err != nil {
		panic(err)
	}
	defer wnd.Destroy()

	scene := newScene(cv.Width(), cv.Height(), 2)
	// scene.models = append(scene.models,
	// 	newXZSquare(2, &IntensityShader{}).moveY(-1),
	// )
	scene.models = append(scene.models,
		// newXZSquare(4, &LineShader{}).moveY(-2),
		newXZSquare(4, newTextureShader("assets/grass.texture.jpg")).moveY(-2),
	)
	scene.models = append(scene.models,
		// newXYSquare(4, &IntensityShader{}).moveZ(-2),
		// newXYSquare(4, &SmoothColorShader{255, 255, 255}).moveZ(-2),
		newXYSquare(4, newTextureShader("assets/brick.texture.jpg")).moveZ(-2),
	)
	// scene.models = append(scene.models,
	// 	parseModel("assets/head.obj", newTextureShader("assets/head.texture.tga")).moveY(-0.1).moveZ(-0.5),
	// )

	if false {
		f, err := os.Create("cpu")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	t := float64(0)

	scene.camera.project(scene)

	pressedKeys := map[string]bool{}
	wnd.KeyDown = func(scancode int, rn rune, name string) {
		pressedKeys[name] = true
	}

	lastFrame := time.Now()
	wnd.MainLoop(func() {
		for key := range pressedKeys {
			delete(pressedKeys, key)
			if key == "KeyD" {
				scene.camera.position.x += 0.1
			}
			if key == "KeyA" {
				scene.camera.position.x -= 0.1
			}
			if key == "KeyW" {
				scene.camera.position.z -= 0.1
			}
			if key == "KeyS" {
				scene.camera.position.z += 0.1
			}
		}

		// scene.lightPosition.z = math.Cos(t/10) * 3
		// scene.lightPosition.x = math.Sin(t/10) * 3
		scene.camera.project(scene)
		cv.PutImageData(scene.render(), 0, 0)

		elapsed := time.Since(lastFrame)
		lastFrame = time.Now()
		t += float64(elapsed.Milliseconds()) / 10
		if true {
			fmt.Println(elapsed.String(), "Triangles = ", scene.trianglesDrawn)
		}
	})
}
