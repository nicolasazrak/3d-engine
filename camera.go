package main

type Camera struct {
	position   Vector3
	angle      float64
	viewMatrix [4][4]float64
}

func newCamera() *Camera {
	return &Camera{
		position:   Vector3{0, 0, 1},
		angle:      0,
		viewMatrix: [4][4]float64{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
	}
}

func (cam *Camera) updateViewMatrix() {
	// Look at camera

	target := Vector3{0, 0, 0}
	zaxis := normalize(minus(cam.position, target))
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
}

func (cam *Camera) project(scene *Scene, triangles []*Triangle) {
	cam.updateViewMatrix()
	scene.projectedLight = matmult(cam.viewMatrix, scene.lightPosition)
	scene.projectedLight = scene.lightPosition
	for _, triangle := range triangles {
		for i, vert := range triangle.verts {
			res := matmult(cam.viewMatrix, vert)
			triangle.viewVerts[i].x = res.x
			triangle.viewVerts[i].y = res.y
			triangle.viewVerts[i].z = res.z

			triangle.screenProjection[i].x = (triangle.viewVerts[i].x/-res.z + 1.) * scene.fWidth / 2.
			triangle.screenProjection[i].y = (triangle.viewVerts[i].y/-res.z + 1.) * scene.fHeight / 2.
			// fmt.Println(triangle.viewProjection)
		}

		for i, normal := range triangle.normals {
			triangle.viewNormals[i] = matmult(cam.viewMatrix, normal)
		}
	}
}
