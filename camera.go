package main

import (
	"math"
)

type Camera struct {
	position     Vector3
	target       Vector3
	angle        float64
	viewMatrix   [4][4]float64
	normalMatrix [4][4]float64
}

func newCamera() *Camera {
	return &Camera{
		position:     Vector3{0, 0, 3},
		target:       Vector3{0, 0, 0},
		angle:        0,
		viewMatrix:   [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		normalMatrix: [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
	}
}

func (cam *Camera) updateViewMatrix() {
	// https://www.3dgep.com/understanding-the-view-matrix/
	// Look at camera

	// target := Vector3{0, 0, 0}
	// target := Vector3{cam.position.x, cam.position.y, cam.position.z - 1}
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

	cam.viewMatrix[3][0] = -math.Max(dotProduct(xaxis, cam.position), 0.0001) // avoids division by 0
	cam.viewMatrix[3][1] = -math.Max(dotProduct(yaxis, cam.position), 0.0001)
	cam.viewMatrix[3][2] = -math.Max(dotProduct(zaxis, cam.position), 0.0001)
	cam.viewMatrix[3][3] = 1

	inverseTranspose(&cam.normalMatrix, cam.viewMatrix)
}

func (cam *Camera) project(scene *Scene) {
	cam.updateViewMatrix()
	scene.projectedLight = matmult(cam.viewMatrix, scene.lightPosition, 1)
	for _, model := range scene.models {
		for _, triangle := range model.triangles {
			for i, vert := range triangle.verts {
				res := matmult(cam.viewMatrix, vert, 1)
				triangle.cameraVerts[i].x = res.x
				triangle.cameraVerts[i].y = res.y
				triangle.cameraVerts[i].z = res.z

				triangle.viewportVerts[i].x = (triangle.cameraVerts[i].x/-res.z + 1.) * scene.fWidth / 2.
				triangle.viewportVerts[i].y = (triangle.cameraVerts[i].y/-res.z + 1.) * scene.fHeight / 2.

			}

			for i, normal := range triangle.normals {
				res := matmult(cam.normalMatrix, normal, 0)
				triangle.cameraNormals[i] = normalize(res)
			}
		}
	}
}
