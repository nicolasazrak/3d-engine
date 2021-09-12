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
	obstacles      []Collisionable
	setAtCalled    int

	// frame stats
	t                 float64
	lastElapsedMillis float64
	trianglesDrawn    int
	lastFrame         time.Time
}

func (scene *Scene) drawTriangle(model *Model, triangle *ProjectedTriangle) {
	v0 := Vector2{
		x: int(math.Round(((triangle.clipVertex[0].x / triangle.clipVertex[0].w) + 1) * (scene.fWidth - 1) * 0.5)),
		y: int(math.Round(((triangle.clipVertex[0].y / triangle.clipVertex[0].w) + 1) * (scene.fHeight - 1) * 0.5)),
	}
	v1 := Vector2{
		x: int(math.Round(((triangle.clipVertex[1].x / triangle.clipVertex[1].w) + 1) * (scene.fWidth - 1) * 0.5)),
		y: int(math.Round(((triangle.clipVertex[1].y / triangle.clipVertex[1].w) + 1) * (scene.fHeight - 1) * 0.5)),
	}
	v2 := Vector2{
		x: int(math.Round(((triangle.clipVertex[2].x / triangle.clipVertex[2].w) + 1) * (scene.fWidth - 1) * 0.5)),
		y: int(math.Round(((triangle.clipVertex[2].y / triangle.clipVertex[2].w) + 1) * (scene.fHeight - 1) * 0.5)),
	}

	pts := []Vector2{v0, v1, v2}

	minbbox, maxbbox := boundingBox(pts, 0, scene.width-1, 0, scene.height-1)
	if minbbox.x >= maxbbox.x || minbbox.y >= maxbbox.y {
		// pseudo frustrum culling
		return
	}

	area := 1. / float64(orient2d(pts[0], pts[1], pts[2].x, pts[2].y))
	if area <= 0 {
		// pseudo backface culling
		return
	}

	A01 := v0.y - v1.y
	B01 := v1.x - v0.x
	A12 := v1.y - v2.y
	B12 := v2.x - v1.x
	A20 := v2.y - v0.y
	B20 := v0.x - v2.x

	w0_row := orient2d(v1, v2, minbbox.x, minbbox.y)
	w1_row := orient2d(v2, v0, minbbox.x, minbbox.y)
	w2_row := orient2d(v0, v1, minbbox.x, minbbox.y)

	scene.trianglesDrawn++

	invZ0 := 1 / triangle.viewVerts[0].z
	invZ1 := 1 / triangle.viewVerts[1].z
	invZ2 := 1 / triangle.viewVerts[2].z

	// https://fgiesen.wordpress.com/2013/02/10/optimizing-the-basic-rasterizer/
	for y := minbbox.y; y <= maxbbox.y; y++ {
		// Barycentric coordinates at start of row
		w0 := w0_row
		w1 := w1_row
		w2 := w2_row

		for x := minbbox.x; x < maxbbox.x; x++ {
			if (w0 | w1 | w2) >= 0 {
				l0 := float64(w0) * area
				l1 := float64(w1) * area
				l2 := float64(w2) * area

				// Should the z-buffer use the ndc value ??
				zPos := 1 / (l0*invZ0 + l1*invZ1 + l2*invZ2)
				idx := x + (y)*scene.width

				if zPos < 0 && zPos > scene.zBuffer[idx] {
					scene.zBuffer[idx] = zPos
					r, g, b := model.shader.shade(scene, triangle, [3]float64{l0, l1, l2}, zPos)
					scene.setAt(x, y, r, g, b)
				}
			}

			// One step to the right
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

func (scene *Scene) drawModels() {
	for _, model := range scene.models {
		for _, triangle := range model.projection {
			scene.drawTriangle(model, triangle)
		}
	}
}

func (scene *Scene) setAt(x int, yInverted int, r uint8, g uint8, b uint8) {
	scene.setAtCalled++
	y := scene.height - yInverted - 1
	pixIdx := (x + y*scene.width) * 4
	scene.pixBuffer[pixIdx] = r
	scene.pixBuffer[pixIdx+1] = g
	scene.pixBuffer[pixIdx+2] = b
	scene.pixBuffer[pixIdx+3] = uint8(255)
}

func (scene *Scene) toImage() *image.RGBA {
	scaledBuffer := make([]uint8, scene.width*scene.height*4*scene.scaleFactor*scene.scaleFactor)

	scaled := &image.RGBA{
		Pix:    scaledBuffer,
		Stride: scene.width * scene.scaleFactor * 4,
		Rect:   image.Rect(0, 0, scene.width*scene.scaleFactor, scene.height*scene.scaleFactor),
	}

	srcPixIdx := 0
	dstPixIdx := 0
	lineLength := scene.width * 4 * scene.scaleFactor

	for y := 0; y < scene.height; y++ {
		startLinePix := dstPixIdx
		for x := 0; x < scene.width; x++ {
			r := scene.pixBuffer[srcPixIdx]
			g := scene.pixBuffer[srcPixIdx+1]
			b := scene.pixBuffer[srcPixIdx+2]
			for i := 0; i < scene.scaleFactor; i++ {
				scaledBuffer[dstPixIdx] = r
				scaledBuffer[dstPixIdx+1] = g
				scaledBuffer[dstPixIdx+2] = b
				scaledBuffer[dstPixIdx+3] = 255
				dstPixIdx += 4
			}
			srcPixIdx += 4
		}
		for i := 0; i < (scene.scaleFactor - 1); i++ {
			copy(scaledBuffer[dstPixIdx:dstPixIdx+lineLength], scaledBuffer[startLinePix:+startLinePix+lineLength])
			dstPixIdx += lineLength
		}
	}

	return scaled
}

func (scene *Scene) cleanBuffer() {
	scene.trianglesDrawn = 0
	scene.setAtCalled = 0
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

	position := scene.camera.getPosition()
	mov := Vector3{}

	for key := range pressedKeys {
		if key == "KeyD" {
			mov.x += moveSpeed
		}
		if key == "KeyA" {
			mov.x -= moveSpeed
		}
		if key == "KeyW" {
			mov.z -= moveSpeed
		}
		if key == "KeyS" {
			mov.z += moveSpeed
		}
		if key == "KeyQ" {
			mov.y -= moveSpeed
		}
		if key == "KeyE" {
			mov.y += moveSpeed
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

	lightMove := scene.camera.transformInput(Vector3{0, 0, -.3})
	scene.lightPosition = plus(position, lightMove)

	if mov.x == 0 && mov.y == 0 && mov.z == 0 {
		return
	}

	mov = scene.camera.transformInput(mov)

	collided := true
	for collided {
		collided = false
		minD := 999999.
		newNorm := Vector3{}
		dst := plus(mov, position)

		for _, obstacle := range scene.obstacles {
			c, norm, d := obstacle.test(position, dst, mov)
			if c && d < minD {
				minD = d
				newNorm = norm
				collided = true
			}
		}

		if collided {
			mov = Vector3{
				x: mov.x - math.Abs(newNorm.x)*mov.x,
				y: mov.y - math.Abs(newNorm.y)*mov.y,
				z: mov.z - math.Abs(newNorm.z)*mov.z,
			}
		}
	}

	scene.camera.move(mov)
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
		obstacles:     []Collisionable{},
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

	scene.pixBuffer = make([]uint8, scene.width*scene.height*4)
	scene.cleanPixBuffer = make([]uint8, scene.width*scene.height*4)
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
	} else {
		red := &SmoothColorShader{170, 30, 30}
		blue := &SmoothColorShader{30, 30, 130}
		green := &SmoothColorShader{30, 143, 23}
		purple := &SmoothColorShader{114, 48, 191}
		grey := &FlatShader{100, 100, 100}

		scenario := []string{
			"XXXXXXXXXXXXXXXXXXXXXXXXXX",
			"X          B     R      X",
			"X          B     R      X",
			"XRRRRRR    B     RRRR   X",
			"X          B            X",
			" X         BBBBBBBBBB   X",
			" X         B            X",
			" X    BBBBBB      G    X",
			" X                G    X",
			" X                G    X",
			" X    RRRRRGGGGGGGG    X",
			"X          G           X",
			"X          G           X",
			"XXXXXXXXXXXXXXXXXXXXXXXX",
		}

		for y, line := range scenario {
			for x, color := range line {
				var model *Model = nil

				if color == 'X' {
					model = newCube(1, purple).moveX(float64(14 - x)).moveZ(float64(6 - y))
				} else if color == 'G' {
					model = newCube(1, green).moveX(float64(14 - x)).moveZ(float64(6 - y))
				} else if color == 'B' {
					model = newCube(1, blue).moveX(float64(14 - x)).moveZ(float64(6 - y))
				} else if color == 'R' {
					model = newCube(1, red).moveX(float64(14 - x)).moveZ(float64(6 - y))
				}

				if model != nil {
					scene.models = append(scene.models, model)
					scene.obstacles = append(scene.obstacles, newCollisionableFromModel(model))
				}

				scene.models = append(scene.models, newCube(1, grey).moveX(float64(14-x)).moveZ(float64(6-y)).moveY(-1.))
			}
		}

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

	f, err := os.Create("cpu")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	pressedKeys := map[string]bool{}
	wnd.KeyDown = func(scancode int, rn rune, name string) {
		pressedKeys[name] = true
	}
	wnd.KeyUp = func(scancode int, rn rune, name string) {
		delete(pressedKeys, name)
	}
	wnd.MouseDown = func(button, x, y int) {
		fmt.Println("Click x =", x/(scene.scaleFactor/2), "y =", y/(scene.scaleFactor/2))
	}

	// scene.processFrame()
	// scene.render()
	// r := scene.render()

	wnd.MainLoop(func() {
		// cv.PutImageData(r, 0, 0)
		scene.processFrame(pressedKeys)
		cv.PutImageData(scene.render(), 0, 0)

		if true {
			fmt.Println(wnd.FPS(), "Triangles = ", scene.trianglesDrawn, "setAtCalled =", scene.setAtCalled)
		}
	})
}
