package main

import (
	"fmt"
	"image"
	"math"
	"os"
	"runtime/pprof"
	"time"

	"github.com/tfriedel6/canvas/sdlcanvas"
)

type Scene struct {
	// buffers
	zBuffer        []float64
	pixBuffer      []uint8
	cleanPixBuffer []uint8
	cleanZBuffer   []float64

	// window
	width       int
	height      int
	fWidth      float64
	fHeight     float64
	scaleFactor int

	// camera, lights and models
	lightPosition  Vector3
	projectedLight Vector3
	camera         Camera
	models         []*Model

	// frame stats
	t                 float64
	lastElapsedMillis float64
	trianglesDrawn    int
	lastFrame         time.Time
}

func (scene *Scene) drawTriangle(model *Model, triangle *Triangle) {
	pts := triangle.viewportVerts

	minbbox, maxbbox := boundingBox(pts, 0, scene.fWidth-1, 0, scene.fHeight-1)
	if minbbox.x >= maxbbox.x || minbbox.y >= maxbbox.y {
		// pseudo frustrum culling
		// fmt.Println("bounding box culled")
		return
	}

	area := 1. / float64(orient2d(pts[0], pts[1], pts[2].x, pts[2].y))
	if area <= 0 {
		// pseudo backface culling
		// fmt.Println("wrong orientation")
		return
	}

	scene.trianglesDrawn++

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
				l0 := float64(w0) * area
				l1 := float64(w1) * area
				l2 := float64(w2) * area

				zPos := 1 / (l0*triangle.invViewZ[0] + l1*triangle.invViewZ[1] + l2*triangle.invViewZ[2])
				idx := int(x) + (int(y))*scene.width

				if zPos > scene.camera.farPlane() && zPos < scene.camera.nearPlane() && zPos > scene.zBuffer[idx] {
					scene.zBuffer[idx] = zPos
					r, g, b := model.shader.shade(scene, triangle, [3]float64{l0, l1, l2}, zPos)
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

func (scene *Scene) isFrontFacing(triangle *Triangle) bool {
	return triangle.viewNormals[0].z >= 0 ||
		triangle.viewNormals[1].z >= 0 ||
		triangle.viewNormals[2].z >= 0
}

func (scene *Scene) isInFrustum(triangle *Triangle) bool {
	return math.Abs(triangle.viewVerts[0].z) > -scene.camera.nearPlane() ||
		math.Abs(triangle.viewVerts[1].z) > -scene.camera.nearPlane() ||
		math.Abs(triangle.viewVerts[2].z) > -scene.camera.nearPlane()
}

func (scene *Scene) clip(triangle *Triangle) []*Triangle {
	// TODO
	return []*Triangle{triangle}
}

func (scene *Scene) drawModels() {
	for _, model := range scene.models {
		for _, triangle := range model.triangles {
			if scene.isFrontFacing(triangle) && scene.isInFrustum(triangle) {
				triangles := scene.clip(triangle)
				for _, clipped := range triangles {
					scene.drawTriangle(model, clipped)
				}
			}
		}
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
	scene.camera.project(scene)
	scene.cleanBuffer()
	scene.drawModels()
	return scene.toImage()
}

func (scene *Scene) handleKeys(pressedKeys map[string]bool) {
	moveSpeed := scene.lastElapsedMillis * 0.003
	rotationSpeed := scene.lastElapsedMillis * 0.001

	for key := range pressedKeys {
		if key == "KeyD" {
			scene.camera.move(moveSpeed, 0, 0)
		}
		if key == "KeyA" {
			scene.camera.move(-moveSpeed, 0, 0)
		}
		if key == "KeyW" {
			scene.camera.move(0, 0, -moveSpeed)
		}
		if key == "KeyS" {
			scene.camera.move(0, 0, moveSpeed)
		}
		if key == "ArrowUp" {
			scene.camera.rotate(0, rotationSpeed)
		}
		if key == "ArrowDown" {
			scene.camera.rotate(0, -rotationSpeed)
		}
		if key == "ArrowLeft" {
			scene.camera.rotate(rotationSpeed, 0)
		}
		if key == "ArrowRight" {
			scene.camera.rotate(-rotationSpeed, 0)
		}
	}
}

func (scene *Scene) processFrame(pressedKeys map[string]bool) {
	scene.lightPosition.z = math.Cos(scene.t/50) * 3
	scene.lightPosition.x = math.Sin(scene.t/50) * 3

	elapsed := time.Since(scene.lastFrame)
	scene.lastFrame = time.Now()
	scene.t += float64(elapsed.Milliseconds()) / 10
	scene.lastElapsedMillis = float64(elapsed.Milliseconds())

	scene.handleKeys(pressedKeys)
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
		camera:        newFPSCamera(),
		lastFrame:     time.Now(),
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
		scene.cleanZBuffer[i] = -999999
	}

	return &scene
}

func addModels(scene *Scene) {
	grassTexture := newTextureShader("assets/grass.texture.jpg")
	brickTexture := newTextureShader("assets/brick.texture.jpg")
	headTexture := newTextureShader("assets/head.texture.tga")
	concreteTexture := newTextureShader("assets/concrete.texture.jpeg")

	grass := newXZSquare(4, grassTexture).scale(1, 1, 1).scaleUV(2, 1).moveY(-2)
	ceiling := newXZSquare(4, concreteTexture).rotateX(math.Pi).scale(1, 1, 1).scaleUV(2, 1).moveY(2)
	leftWall := newXYSquare(4, brickTexture).rotateY(math.Pi/2).scaleUV(4, 1).moveX(-2)
	rightWall := newXYSquare(4, brickTexture).rotateY(-math.Pi/2).scaleUV(4, 1).moveX(2)
	backWall := newXYSquare(4, brickTexture).moveZ(-2)
	head := parseModel("assets/head.obj", headTexture)

	scene.models = append(scene.models, grass)
	scene.models = append(scene.models, leftWall)
	scene.models = append(scene.models, rightWall)
	scene.models = append(scene.models, head)
	scene.models = append(scene.models, ceiling)
	scene.models = append(scene.models, backWall)
}

func takeProfile() func() {
	f, err := os.Create("cpu")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	return func() {
		pprof.StopCPUProfile()
	}
}

func main() {
	wnd, cv, err := sdlcanvas.CreateWindow(1000, 1000, "Hello")
	if err != nil {
		panic(err)
	}
	defer wnd.Destroy()

	scene := newScene(cv.Width(), cv.Height(), 2)
	addModels(scene)

	// endProfile := takeProfile()
	// defer endProfile()

	pressedKeys := map[string]bool{}
	wnd.KeyDown = func(scancode int, rn rune, name string) {
		pressedKeys[name] = true
	}
	wnd.KeyUp = func(scancode int, rn rune, name string) {
		delete(pressedKeys, name)
	}

	// scene.processFrame()
	// scene.render()
	// r := scene.render()

	wnd.MainLoop(func() {
		// cv.PutImageData(r, 0, 0)
		scene.processFrame(pressedKeys)
		cv.PutImageData(scene.render(), 0, 0)

		if false {
			fmt.Println(wnd.FPS(), "Triangles = ", scene.trianglesDrawn)
		}
	})
}
