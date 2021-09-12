package main

const INSIDE = 0 // 000000
const LEFT = 1   // 000001
const RIGHT = 2  // 000010
const BOTTOM = 4 // 000100
const TOP = 8    // 001000
const FRONT = 16 // 010000
const BACK = 32  // 100000

func OutCode(triangle *ProjectedTriangle, vertexIdx int) int {
	code := INSIDE

	if triangle.clipVertex[vertexIdx].x < -triangle.clipVertex[vertexIdx].w {
		code |= LEFT
	}

	if triangle.clipVertex[vertexIdx].x > triangle.clipVertex[vertexIdx].w {
		code |= RIGHT
	}

	if triangle.clipVertex[vertexIdx].y < -triangle.clipVertex[vertexIdx].w {
		code |= BOTTOM
	}

	if triangle.clipVertex[vertexIdx].y > triangle.clipVertex[vertexIdx].w {
		code |= TOP
	}

	if triangle.clipVertex[vertexIdx].z < -triangle.clipVertex[vertexIdx].w {
		code |= FRONT
	}

	if triangle.clipVertex[vertexIdx].z > triangle.clipVertex[vertexIdx].w {
		code |= BACK
	}

	return code
}

// Formulas taken from https://stackoverflow.com/questions/60910464/at-what-stage-is-clipping-performed-in-the-graphics-pipeline
func intersectionTop(triangle *ProjectedTriangle, idx0, idx1 int) float64 {
	return (triangle.clipVertex[idx0].y - triangle.clipVertex[idx0].w) / ((triangle.clipVertex[idx0].y - triangle.clipVertex[idx0].w) - (triangle.clipVertex[idx1].y - triangle.clipVertex[idx1].w))
}

func intersectionBottom(triangle *ProjectedTriangle, idx0, idx1 int) float64 {
	return (triangle.clipVertex[idx0].y + triangle.clipVertex[idx0].w) / ((triangle.clipVertex[idx0].y + triangle.clipVertex[idx0].w) - (triangle.clipVertex[idx1].y + triangle.clipVertex[idx1].w))
}

func intersectionRight(triangle *ProjectedTriangle, idx0, idx1 int) float64 {
	return (triangle.clipVertex[idx0].x - triangle.clipVertex[idx0].w) / ((triangle.clipVertex[idx0].x - triangle.clipVertex[idx0].w) - (triangle.clipVertex[idx1].x - triangle.clipVertex[idx1].w))
}

func intersectionLeft(triangle *ProjectedTriangle, idx0, idx1 int) float64 {
	return (triangle.clipVertex[idx0].x + triangle.clipVertex[idx0].w) / ((triangle.clipVertex[idx0].x + triangle.clipVertex[idx0].w) - (triangle.clipVertex[idx1].x + triangle.clipVertex[idx1].w))
}

func intersectionFront(triangle *ProjectedTriangle, idx0, idx1 int) float64 {
	return (triangle.clipVertex[idx0].z - triangle.clipVertex[idx0].w) / ((triangle.clipVertex[idx0].z - triangle.clipVertex[idx0].w) - (triangle.clipVertex[idx1].z - triangle.clipVertex[idx1].w))
}

func intersectionBack(triangle *ProjectedTriangle, idx0, idx1 int) float64 {
	return (triangle.clipVertex[idx0].z + triangle.clipVertex[idx0].w) / ((triangle.clipVertex[idx0].z + triangle.clipVertex[idx0].w) - (triangle.clipVertex[idx1].z + triangle.clipVertex[idx1].w))
}

func findT(triangle *ProjectedTriangle, idx0, idx1 int, plane int) float64 {
	if plane == LEFT {
		return intersectionLeft(triangle, idx0, idx1)
	}

	if plane == RIGHT {
		return intersectionRight(triangle, idx0, idx1)
	}

	if plane == TOP {
		return intersectionTop(triangle, idx0, idx1)
	}

	if plane == BOTTOM {
		return intersectionBottom(triangle, idx0, idx1)
	}

	if plane == FRONT {
		return intersectionFront(triangle, idx0, idx1)
	}

	if plane == BACK {
		return intersectionBack(triangle, idx0, idx1)
	}

	return 0
}

func clipTriangleWithOneVertexInside(triangle *ProjectedTriangle, planeToClip int, insideVertex int) *ProjectedTriangle {
	nextIdx := (insideVertex + 1) % 3
	otherIdx := (insideVertex + 2) % 3
	t1 := findT(triangle, insideVertex, nextIdx, planeToClip)
	t2 := findT(triangle, insideVertex, otherIdx, planeToClip)

	return &ProjectedTriangle{
		viewVerts: []Vector3{
			ponderateVec3(triangle.viewVerts[nextIdx], triangle.viewVerts[insideVertex], t1),
			ponderateVec3(triangle.viewVerts[otherIdx], triangle.viewVerts[insideVertex], t2),
			triangle.viewVerts[insideVertex],
		},
		clipVertex: []Vector4{
			ponderateVec4(triangle.clipVertex[nextIdx], triangle.clipVertex[insideVertex], t1),
			ponderateVec4(triangle.clipVertex[otherIdx], triangle.clipVertex[insideVertex], t2),
			triangle.clipVertex[insideVertex],
		},
		viewNormals: []Vector3{
			ponderateVec3(triangle.viewNormals[nextIdx], triangle.viewNormals[insideVertex], t1),
			ponderateVec3(triangle.viewNormals[otherIdx], triangle.viewNormals[insideVertex], t2),
			triangle.viewNormals[insideVertex],
		},
		uvMapping: [][]float64{
			ponderateSlice3(triangle.uvMapping[nextIdx], triangle.uvMapping[insideVertex], t1),
			ponderateSlice3(triangle.uvMapping[otherIdx], triangle.uvMapping[insideVertex], t2),
			triangle.uvMapping[insideVertex],
		},
		lightIntensity: []float64{
			t1*triangle.lightIntensity[nextIdx] + (1-t1)*triangle.lightIntensity[insideVertex],
			t1*triangle.lightIntensity[otherIdx] + (1-t1)*triangle.lightIntensity[insideVertex],
			triangle.lightIntensity[insideVertex],
		},
	}
}

func clipTriangleWithTwoVertexInside(triangle *ProjectedTriangle, planeToClip int, outsideVertex int) (*ProjectedTriangle, *ProjectedTriangle) {
	nextIdx := (outsideVertex + 1) % 3
	otherIdx := (outsideVertex + 2) % 3
	t1 := findT(triangle, outsideVertex, nextIdx, planeToClip)
	t2 := findT(triangle, outsideVertex, otherIdx, planeToClip)

	newVewVert := ponderateVec3(triangle.viewVerts[nextIdx], triangle.viewVerts[outsideVertex], t1)
	newClipVert := ponderateVec4(triangle.clipVertex[nextIdx], triangle.clipVertex[outsideVertex], t1)
	newViewNormal := ponderateVec3(triangle.viewNormals[nextIdx], triangle.viewNormals[outsideVertex], t1)
	newUvMapping := ponderateSlice3(triangle.uvMapping[nextIdx], triangle.uvMapping[outsideVertex], t1)
	newLightIntensity := triangle.lightIntensity[nextIdx]*t1 + (1-t1)*triangle.lightIntensity[outsideVertex]

	triangle1 := ProjectedTriangle{
		viewVerts: []Vector3{
			newVewVert,
			triangle.viewVerts[nextIdx],
			triangle.viewVerts[otherIdx],
		},
		clipVertex: []Vector4{
			newClipVert,
			triangle.clipVertex[nextIdx],
			triangle.clipVertex[otherIdx],
		},
		viewNormals: []Vector3{
			newViewNormal,
			triangle.viewNormals[nextIdx],
			triangle.viewNormals[otherIdx],
		},
		uvMapping: [][]float64{
			newUvMapping,
			triangle.uvMapping[nextIdx],
			triangle.uvMapping[otherIdx],
		},
		lightIntensity: []float64{
			newLightIntensity,
			triangle.lightIntensity[nextIdx],
			triangle.lightIntensity[otherIdx],
		},
	}
	triangle2 := ProjectedTriangle{
		viewVerts: []Vector3{
			ponderateVec3(triangle.viewVerts[otherIdx], triangle.viewVerts[outsideVertex], t2),
			newVewVert,
			triangle.viewVerts[otherIdx],
		},
		clipVertex: []Vector4{
			ponderateVec4(triangle.clipVertex[otherIdx], triangle.clipVertex[outsideVertex], t2),
			newClipVert,
			triangle.clipVertex[otherIdx],
		},
		viewNormals: []Vector3{
			ponderateVec3(triangle.viewNormals[otherIdx], triangle.viewNormals[outsideVertex], t2),
			newViewNormal,
			triangle.viewNormals[otherIdx],
		},
		uvMapping: [][]float64{
			ponderateSlice3(triangle.uvMapping[otherIdx], triangle.uvMapping[outsideVertex], t2),
			newUvMapping,
			triangle.uvMapping[otherIdx],
		},
		lightIntensity: []float64{
			t2*triangle.lightIntensity[otherIdx] + (1-t2)*triangle.lightIntensity[outsideVertex],
			newLightIntensity,
			triangle.lightIntensity[otherIdx],
		},
	}

	return &triangle1, &triangle2
}

func getInsidePlaneVertex(vertCodes []int, plane int) int {
	insideVertex := -1
	for i, code := range vertCodes {
		if code&plane == 0 {
			insideVertex = i
			break
		}
	}

	return insideVertex
}

func getOutsidePlaneVertex(vertCodes []int, plane int) int {
	outsideVertex := -1
	for i, code := range vertCodes {
		if code&plane != 0 {
			outsideVertex = i
			break
		}
	}

	return outsideVertex
}

func clipTriangle(triangle *ProjectedTriangle) []*ProjectedTriangle {
	// https://gabrielgambetta.com/computer-graphics-from-scratch/11-clipping.html

	baseVertex0Opt := OutCode(triangle, 0)
	baseVertex1Opt := OutCode(triangle, 1)
	baseVertex2Opt := OutCode(triangle, 2)
	if (baseVertex0Opt | baseVertex1Opt | baseVertex2Opt) == INSIDE {
		// Base check inside, doesn't need clipping
		return []*ProjectedTriangle{triangle}
	}

	projection := []*ProjectedTriangle{triangle}
	for _, planeToClip := range []int{LEFT, RIGHT, TOP, BOTTOM, FRONT, BACK} {

		planeTriangles := []*ProjectedTriangle{}

		// Crop each triangle against the given plane
		for _, projected := range projection {
			vertex0Opt := OutCode(projected, 0)
			vertex1Opt := OutCode(projected, 1)
			vertex2Opt := OutCode(projected, 2)
			vertCodes := []int{vertex0Opt, vertex1Opt, vertex2Opt}

			insidePlane := 0
			for _, code := range vertCodes {
				if code&planeToClip == 0 {
					insidePlane++
				}
			}

			if insidePlane == 0 {
				// Trivial case reject
				continue
			}

			if insidePlane == 1 {
				// crop single tringle
				insideVertex := getInsidePlaneVertex(vertCodes, planeToClip)
				newTriangle := clipTriangleWithOneVertexInside(projected, planeToClip, insideVertex)
				planeTriangles = append(planeTriangles, newTriangle)
			}

			if insidePlane == 2 {
				// split in two triangles
				outsideVertex := getOutsidePlaneVertex(vertCodes, planeToClip)
				t1, t2 := clipTriangleWithTwoVertexInside(projected, planeToClip, outsideVertex)
				planeTriangles = append(planeTriangles, t1)
				planeTriangles = append(planeTriangles, t2)
			}

			if insidePlane == 3 {
				// keep the triangle as is
				planeTriangles = append(planeTriangles, projected)
			}
		}

		projection = planeTriangles
	}

	return projection
}
