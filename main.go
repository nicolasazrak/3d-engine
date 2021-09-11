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

func (scene *Scene) drawTriangle(model *Model, triangle *ProjectedTriangle) {
	pts := []Vector2{
		{((triangle.clipVertex[0].x / triangle.clipVertex[0].w) + 1) * scene.fWidth * 0.5, ((triangle.clipVertex[0].y / triangle.clipVertex[0].w) + 1) * scene.fHeight * 0.5},
		{((triangle.clipVertex[1].x / triangle.clipVertex[1].w) + 1) * scene.fWidth * 0.5, ((triangle.clipVertex[1].y / triangle.clipVertex[1].w) + 1) * scene.fHeight * 0.5},
		{((triangle.clipVertex[2].x / triangle.clipVertex[2].w) + 1) * scene.fWidth * 0.5, ((triangle.clipVertex[2].y / triangle.clipVertex[2].w) + 1) * scene.fHeight * 0.5},
	}

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

	// TODO re implement https://fgiesen.wordpress.com/2013/02/10/optimizing-the-basic-rasterizer/ it was buggy before
	for y := minbbox.y; y <= maxbbox.y; y++ {
		for x := minbbox.x; x < maxbbox.x; x++ {
			w0 := orient2d(pts[1], pts[2], x, y)
			w1 := orient2d(pts[2], pts[0], x, y)
			w2 := orient2d(pts[0], pts[1], x, y)
			if (w0 | w1 | w2) >= 0 {
				l0 := float64(w0) * area
				l1 := float64(w1) * area
				l2 := float64(w2) * area

				// Should the z-buffer use the ndc value ??
				zPos := 1 / (l0*(1/triangle.viewVerts[0].z) + l1*(1/triangle.viewVerts[1].z) + l2*(1/triangle.viewVerts[2].z))
				idx := int(x) + (int(y))*scene.width

				if zPos > scene.zBuffer[idx] {
					scene.zBuffer[idx] = zPos
					r, g, b := model.shader.shade(scene, triangle, [3]float64{l0, l1, l2}, zPos)
					scene.setAt(int(x), int(y), r, g, b)
				}
			}
		}
	}
}

func (scene *Scene) drawModels() {
	for _, model := range scene.models {
		for _, triangle := range model.projection {
			scene.drawTriangle(model, triangle)
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
	rotationSpeed := scene.lastElapsedMillis * 0.0025

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
		if key == "KeyQ" {
			scene.camera.move(0, -moveSpeed, 0)
		}
		if key == "KeyE" {
			scene.camera.move(0, moveSpeed, 0)
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
	headTexture := newTextureShader("assets/head.texture.tga")
	concreteTexture := newTextureShader("assets/concrete.texture.jpeg")
	brickTexture := newTextureShader("assets/brick.texture.jpg")

	grass := newXZSquare(4, grassTexture).scale(1, 1, 1).scaleUV(2, 1).moveY(-2)
	ceiling := newXZSquare(4, concreteTexture).rotateX(math.Pi).scale(1, 1, 1).scaleUV(2, 1).moveY(2)
	leftWall := newXYSquare(4, brickTexture).rotateY(math.Pi/2).scaleUV(4, 1).moveX(-2)
	rightWall := newXYSquare(4, brickTexture).rotateY(-math.Pi/2).scaleUV(4, 1).moveX(2)
	head := parseModel("assets/head.obj", headTexture)
	backWall := newXYSquare(4, brickTexture).moveZ(-2)

	if false {
		scene.models = append(scene.models, grass)
		scene.models = append(scene.models, leftWall)
		scene.models = append(scene.models, rightWall)
		scene.models = append(scene.models, head)
		scene.models = append(scene.models, ceiling)
		scene.models = append(scene.models, backWall)
	}

	red := &LineShader{170, 30, 30, 255, 255, 255, 0.01}
	blue := &LineShader{30, 30, 130, 255, 255, 255, 0.01}
	green := &LineShader{30, 143, 23, 255, 255, 255, 0.01}
	purple := &LineShader{114, 48, 191, 255, 255, 255, 0.01}
	grey := &FlatShader{100, 100, 100}

	scenario := []string{
		"XXXXXXXXXXXXXXXXXXXXXXXXX",
		"X          B     R      X",
		"X          B     R      X",
		"XRRRRRR    B     RRR    X ",
		"X          B            X",
		"X          B            X",
		"X     BBBBBB            X",
		"X                 G     X",
		"X                 G     X",
		"X  G       G      G     X",
		"XXXXXXXXXXXXXXXXXXXXXXXXX",
	}

	for y, line := range scenario {
		for x, color := range line {
			if color == 'X' {
				scene.models = append(scene.models, newCube(1, purple).moveX(float64(14-x)).moveZ(float64(6-y)))
			} else if color == 'G' {
				scene.models = append(scene.models, newCube(1, green).moveX(float64(14-x)).moveZ(float64(6-y)))
			} else if color == 'B' {
				scene.models = append(scene.models, newCube(1, blue).moveX(float64(14-x)).moveZ(float64(6-y)))
			} else if color == 'R' {
				scene.models = append(scene.models, newCube(1, red).moveX(float64(14-x)).moveZ(float64(6-y)))
			}

			scene.models = append(scene.models, newCube(1, grey).moveX(float64(14-x)).moveZ(float64(6-y)).moveY(-1.))
		}
	}
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
	wnd, cv, err := sdlcanvas.CreateWindow(1024, 1024, "Hello")
	if err != nil {
		panic(err)
	}
	defer wnd.Destroy()

	scene := newScene(cv.Width(), cv.Height(), 4)
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

		if true {
			fmt.Println(wnd.FPS(), "Triangles = ", scene.trianglesDrawn)
		}
	})
}
