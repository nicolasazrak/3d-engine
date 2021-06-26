package main

import "math"

type Camera interface {
	project(scene *Scene)
	nearPlane() float64
	farPlane() float64
	move(x, y, z float64)
	rotate(yaw, pitch float64)
}

type LookAtCamera struct {
	position     Vector3
	target       Vector3
	angle        float64
	viewMatrix   [4][4]float64
	normalMatrix [4][4]float64
}

type FPSCamera struct {
	position     Vector3
	pitch        float64
	yaw          float64
	viewMatrix   [4][4]float64
	normalMatrix [4][4]float64
}

func newLookAtCamera() Camera {
	return &LookAtCamera{
		position:     Vector3{0, 0, 4},
		target:       Vector3{0, 0, -1},
		angle:        0,
		viewMatrix:   [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		normalMatrix: [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
	}
}

func newFPSCamera() Camera {
	return &FPSCamera{
		position:     Vector3{0, 0, 4},
		pitch:        0.,
		yaw:          0.,
		viewMatrix:   [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		normalMatrix: [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
	}
}

/** Look at camera */

func (cam *LookAtCamera) move(x, y, z float64) {
	cam.position.x += x
	cam.position.y += y
	cam.position.z += z
}

func (cam *LookAtCamera) farPlane() float64 {
	return -30.
}

func (cam *LookAtCamera) nearPlane() float64 {
	return -.1
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

func (cam *LookAtCamera) project(scene *Scene) {
	cam.updateViewMatrix()
	scene.projectedLight = matmult(cam.viewMatrix, scene.lightPosition, 1)
	for _, model := range scene.models {
		for _, triangle := range model.triangles {
			projectTriangle(triangle, scene.fWidth, scene.fHeight, cam.viewMatrix, cam.normalMatrix)
		}
	}
}

/** FPS camera */

func (cam *FPSCamera) project(scene *Scene) {
	cam.updateViewMatrix()
	scene.projectedLight = matmult(cam.viewMatrix, scene.lightPosition, 1)
	for _, model := range scene.models {
		for _, triangle := range model.triangles {
			projectTriangle(triangle, scene.fWidth, scene.fHeight, cam.viewMatrix, cam.normalMatrix)
		}
	}
}

func (cam *FPSCamera) nearPlane() float64 {
	return -.5
}

func (cam *FPSCamera) farPlane() float64 {
	return -30
}

func (cam *FPSCamera) move(x, y, z float64) {
	cam.position.x += z*math.Sin(cam.yaw) + x*math.Sin(cam.yaw+math.Pi/2)
	cam.position.y += y
	cam.position.z += z*math.Cos(cam.yaw) + x*math.Cos(cam.yaw+math.Pi/2)
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

/** General method */

func projectTriangle(triangle *Triangle, width float64, height float64, viewMatrix [4][4]float64, normalMatrix [4][4]float64) {
	for i := 0; i < 3; i++ {
		res := matmult(viewMatrix, triangle.worldVerts[i], 1)
		triangle.viewVerts[i].x = res.x
		triangle.viewVerts[i].y = res.y
		triangle.viewVerts[i].z = res.z

		triangle.viewportVerts[i].x = (triangle.viewVerts[i].x/-res.z + 1.) * width / 2.
		triangle.viewportVerts[i].y = (triangle.viewVerts[i].y/-res.z + 1.) * height / 2.

		res2 := matmult(normalMatrix, triangle.normals[i], 1 /* This should be 0. Why do I need to make it 1? */)
		triangle.viewNormals[i] = normalize(res2)

		triangle.uvMappingCorrected[i][0] = triangle.uvMapping[i][0] / triangle.viewVerts[i].z
		triangle.uvMappingCorrected[i][1] = triangle.uvMapping[i][1] / triangle.viewVerts[i].z
		triangle.uvMappingCorrected[i][2] = triangle.uvMapping[i][2] / triangle.viewVerts[i].z

		triangle.invViewZ[i] = 1 / triangle.viewVerts[i].z
	}
}
