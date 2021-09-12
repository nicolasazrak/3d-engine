package main

import (
	"math"
)

type Camera interface {
	project(scene *Scene)
	move(move Vector3)
	rotate(yaw, pitch float64)
	transformInput(inputMove Vector3) Vector3
	getPosition() Vector3
}

type LookAtCamera struct {
	position         Vector3
	target           Vector3
	angle            float64
	viewMatrix       [4][4]float64
	normalMatrix     [4][4]float64
	projectionMatrix [4][4]float64
}

type FPSCamera struct {
	position         Vector3
	pitch            float64
	yaw              float64
	viewMatrix       [4][4]float64
	normalMatrix     [4][4]float64
	projectionMatrix [4][4]float64
}

func buildProjectionMatrix() [4][4]float64 {
	nearPlane := .1
	farPlane := 50.

	leftPlane := -.1
	rightPlane := .1

	topPlane := .1
	bottomPlane := -.1

	fovX := math.Pi / 2
	fovY := math.Pi / 2

	useOpenGlMatrix := true

	// http://www.songho.ca/opengl/gl_projectionmatrix.html
	if useOpenGlMatrix {
		return [4][4]float64{
			{(2 * nearPlane) / (rightPlane - leftPlane), 0, 0, 0},
			{0, (2 * nearPlane) / (topPlane - bottomPlane), 0, 0},
			{(rightPlane + leftPlane) / (rightPlane - leftPlane), (topPlane + bottomPlane) / (topPlane - bottomPlane), -((farPlane + nearPlane) / (farPlane - nearPlane)), -1},
			{0, 0, -((2 * farPlane * nearPlane) / (farPlane - nearPlane)), 0},
		}
	} else {
		return [4][4]float64{
			{1 / math.Tan(fovX/2), 0, 0, 0},
			{0, 1 / math.Tan(fovY/2), 0, 0},
			{0, 0, -(farPlane / (farPlane - nearPlane)), -1},
			{0, 0, -((farPlane * nearPlane) / (farPlane - nearPlane)), 0},
		}
	}
}

func newLookAtCamera() Camera {
	return &LookAtCamera{
		position:         Vector3{0, 0, 4},
		target:           Vector3{0, 0, -1},
		angle:            0,
		viewMatrix:       [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		normalMatrix:     [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		projectionMatrix: buildProjectionMatrix(),
	}
}

func newFPSCamera() Camera {
	return &FPSCamera{
		position:         Vector3{0, 0, 4},
		pitch:            0.,
		yaw:              0.,
		viewMatrix:       [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		normalMatrix:     [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		projectionMatrix: buildProjectionMatrix(),
	}
}

/** Look at camera */

func (cam *LookAtCamera) move(move Vector3) {
	cam.position.x += move.x
	cam.position.y += move.y
	cam.position.z += move.z
}

func (cam *LookAtCamera) updateViewMatrix() {
	// https://www.3dgep.com/understanding-the-view-matrix/
	// Look at camera

	// Two possible targets
	cam.target = Vector3{0, 0, 0}
	// cam.target = Vector3{cam.position.x, cam.position.y, cam.position.z - 1}

	zaxis := normalize(minus(cam.position, cam.target))
	xaxis := normalize(crossProduct(Vector3{0, 1, 0}, zaxis))
	yaxis := crossProduct(zaxis, xaxis)

	cam.viewMatrix[0][0] = xaxis.x
	cam.viewMatrix[0][1] = yaxis.x
	cam.viewMatrix[0][2] = zaxis.x
	cam.viewMatrix[0][3] = 0

	cam.viewMatrix[1][0] = xaxis.y
	cam.viewMatrix[1][1] = yaxis.y
	cam.viewMatrix[1][2] = zaxis.y
	cam.viewMatrix[1][3] = 0

	cam.viewMatrix[2][0] = xaxis.z
	cam.viewMatrix[2][1] = yaxis.z
	cam.viewMatrix[2][2] = zaxis.z
	cam.viewMatrix[2][3] = 0

	cam.viewMatrix[3][0] = -dotProduct(xaxis, cam.position)
	cam.viewMatrix[3][1] = -dotProduct(yaxis, cam.position)
	cam.viewMatrix[3][2] = -dotProduct(zaxis, cam.position)
	cam.viewMatrix[3][3] = 1

	inverseTranspose(&cam.normalMatrix, cam.viewMatrix)
}

func (cam *LookAtCamera) rotate(yaw, pitch float64) {
	// Not supported...
}

func (cam *LookAtCamera) getPosition() Vector3 {
	return cam.position
}

func (cam *LookAtCamera) transformInput(inputMove Vector3) Vector3 {
	return inputMove
}

func (cam *LookAtCamera) project(scene *Scene) {
	cam.updateViewMatrix()
	scene.projectedLight = matmult(cam.viewMatrix, scene.lightPosition, 1)
	for _, model := range scene.models {
		projection := []*ProjectedTriangle{}
		for _, triangle := range model.triangles {
			projection = append(projection, projectTriangle(triangle, cam.viewMatrix, cam.normalMatrix, cam.projectionMatrix, scene.projectedLight)...)
		}
		model.projection = projection
	}
}

/** FPS camera */

func (cam *FPSCamera) project(scene *Scene) {
	cam.updateViewMatrix()
	scene.projectedLight = matmult(cam.viewMatrix, scene.lightPosition, 1)
	for _, model := range scene.models {
		model.projection = []*ProjectedTriangle{}
		for _, triangle := range model.triangles {
			model.projection = append(model.projection, projectTriangle(triangle, cam.viewMatrix, cam.normalMatrix, cam.projectionMatrix, scene.projectedLight)...)
		}
	}
}

func (cam *FPSCamera) transformInput(inputMove Vector3) Vector3 {
	return Vector3{
		x: inputMove.z*math.Sin(cam.yaw) + inputMove.x*math.Sin(cam.yaw+math.Pi/2),
		y: inputMove.y,
		z: inputMove.z*math.Cos(cam.yaw) + inputMove.x*math.Cos(cam.yaw+math.Pi/2),
	}
}

func (cam *FPSCamera) move(move Vector3) {
	cam.position.x += move.x
	cam.position.y += move.y
	cam.position.z += move.z
}

func (cam *FPSCamera) rotate(yaw, pitch float64) {
	cam.pitch += pitch
	cam.yaw += yaw
}

func (cam *FPSCamera) updateViewMatrix() {
	// I assume the values are already converted to radians.
	cosPitch := math.Cos(cam.pitch)
	sinPitch := math.Sin(cam.pitch)
	cosYaw := math.Cos(cam.yaw)
	sinYaw := math.Sin(cam.yaw)

	xaxis := Vector3{cosYaw, 0, -sinYaw}
	yaxis := Vector3{sinYaw * sinPitch, cosPitch, cosYaw * sinPitch}
	zaxis := Vector3{sinYaw * cosPitch, -sinPitch, cosPitch * cosYaw}

	cam.viewMatrix[0][0] = xaxis.x
	cam.viewMatrix[0][1] = yaxis.x
	cam.viewMatrix[0][2] = zaxis.x
	cam.viewMatrix[0][3] = 0

	cam.viewMatrix[1][0] = xaxis.y
	cam.viewMatrix[1][1] = yaxis.y
	cam.viewMatrix[1][2] = zaxis.y
	cam.viewMatrix[1][3] = 0

	cam.viewMatrix[2][0] = xaxis.z
	cam.viewMatrix[2][1] = yaxis.z
	cam.viewMatrix[2][2] = zaxis.z
	cam.viewMatrix[2][3] = 0

	cam.viewMatrix[3][0] = -dotProduct(xaxis, cam.position)
	cam.viewMatrix[3][1] = -dotProduct(yaxis, cam.position)
	cam.viewMatrix[3][2] = -dotProduct(zaxis, cam.position)
	cam.viewMatrix[3][3] = 1

	inverseTranspose(&cam.normalMatrix, cam.viewMatrix)
}

func (cam *FPSCamera) getPosition() Vector3 {
	return cam.position
}

/** General method */

func projectTriangle(originalTriangle *Triangle, viewMatrix [4][4]float64, normalMatrix [4][4]float64, projectionMatrix [4][4]float64, light Vector3) []*ProjectedTriangle {
	projection := newProjectedTriangle()

	for i := 0; i < 3; i++ {
		view := matmult4(viewMatrix, originalTriangle.worldVerts[i], 1)
		clip := matmult4h(projectionMatrix, view)
		normal := matmult(normalMatrix, originalTriangle.normals[i], 1 /* This should be 0. Why do I need to make it 1? */)

		projection.clipVertex[i] = clip
		projection.viewVerts[i].x = view.x / view.w
		projection.viewVerts[i].y = view.y / view.w
		projection.viewVerts[i].z = view.z / view.w
		projection.viewNormals[i] = normalize(normal)
		projection.uvMapping = originalTriangle.uvMapping
		projection.lightIntensity[i] = 1. / (norm(minus(projection.viewVerts[i], light)))
	}

	// TODO add backface culling
	return clipTriangle(projection)
}
