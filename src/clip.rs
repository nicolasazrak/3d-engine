use crate::algebra::*;
use crate::model::*;

static INSIDE: u8 = 0; // 000000
static LEFT: u8 = 1; // 000001
static RIGHT: u8 = 2; // 000010
static BOTTOM: u8 = 4; // 000100
static TOP: u8 = 8; // 001000
static FRONT: u8 = 16; // 010000
static BACK: u8 = 32; // 100000

pub fn out_code(triangle: &ProjectedTriangle, vertex_idx: usize) -> u8 {
    let mut code = INSIDE;

    if triangle.clip_vertex[vertex_idx].x < -triangle.clip_vertex[vertex_idx].w {
        code |= LEFT;
    }

    if triangle.clip_vertex[vertex_idx].x > triangle.clip_vertex[vertex_idx].w {
        code |= RIGHT;
    }

    if triangle.clip_vertex[vertex_idx].y < -triangle.clip_vertex[vertex_idx].w {
        code |= BOTTOM;
    }

    if triangle.clip_vertex[vertex_idx].y > triangle.clip_vertex[vertex_idx].w {
        code |= TOP;
    }

    if triangle.clip_vertex[vertex_idx].z < -triangle.clip_vertex[vertex_idx].w {
        code |= FRONT;
    }

    if triangle.clip_vertex[vertex_idx].z > triangle.clip_vertex[vertex_idx].w {
        code |= BACK;
    }

    return code;
}

// Formulas taken from https://stackoverflow.com/questions/60910464/at-what-stage-is-clipping-performed-in-the-graphics-pipeline
pub fn intersection_top(triangle: &ProjectedTriangle, idx0: usize, idx1: usize) -> f32 {
    return (triangle.clip_vertex[idx0].y - triangle.clip_vertex[idx0].w)
        / ((triangle.clip_vertex[idx0].y - triangle.clip_vertex[idx0].w) - (triangle.clip_vertex[idx1].y - triangle.clip_vertex[idx1].w));
}

pub fn intersection_bottom(triangle: &ProjectedTriangle, idx0: usize, idx1: usize) -> f32 {
    return (triangle.clip_vertex[idx0].y + triangle.clip_vertex[idx0].w)
        / ((triangle.clip_vertex[idx0].y + triangle.clip_vertex[idx0].w) - (triangle.clip_vertex[idx1].y + triangle.clip_vertex[idx1].w));
}

pub fn intersection_right(triangle: &ProjectedTriangle, idx0: usize, idx1: usize) -> f32 {
    return (triangle.clip_vertex[idx0].x - triangle.clip_vertex[idx0].w)
        / ((triangle.clip_vertex[idx0].x - triangle.clip_vertex[idx0].w) - (triangle.clip_vertex[idx1].x - triangle.clip_vertex[idx1].w));
}

pub fn intersection_left(triangle: &ProjectedTriangle, idx0: usize, idx1: usize) -> f32 {
    return (triangle.clip_vertex[idx0].x + triangle.clip_vertex[idx0].w)
        / ((triangle.clip_vertex[idx0].x + triangle.clip_vertex[idx0].w) - (triangle.clip_vertex[idx1].x + triangle.clip_vertex[idx1].w));
}

pub fn intersection_front(triangle: &ProjectedTriangle, idx0: usize, idx1: usize) -> f32 {
    return (triangle.clip_vertex[idx0].z - triangle.clip_vertex[idx0].w)
        / ((triangle.clip_vertex[idx0].z - triangle.clip_vertex[idx0].w) - (triangle.clip_vertex[idx1].z - triangle.clip_vertex[idx1].w));
}

pub fn intersection_back(triangle: &ProjectedTriangle, idx0: usize, idx1: usize) -> f32 {
    return (triangle.clip_vertex[idx0].z + triangle.clip_vertex[idx0].w)
        / ((triangle.clip_vertex[idx0].z + triangle.clip_vertex[idx0].w) - (triangle.clip_vertex[idx1].z + triangle.clip_vertex[idx1].w));
}

pub fn find_t(triangle: &ProjectedTriangle, idx0: usize, idx1: usize, plane: u8) -> f32 {
    if plane == LEFT {
        return intersection_left(triangle, idx0, idx1);
    }

    if plane == RIGHT {
        return intersection_right(triangle, idx0, idx1);
    }

    if plane == TOP {
        return intersection_top(triangle, idx0, idx1);
    }

    if plane == BOTTOM {
        return intersection_bottom(triangle, idx0, idx1);
    }

    if plane == FRONT {
        return intersection_front(triangle, idx0, idx1);
    }

    if plane == BACK {
        return intersection_back(triangle, idx0, idx1);
    }

    return 0.;
}

pub fn clip_triangle_with_one_vertex_inside(triangle: &ProjectedTriangle, plane_to_clip: u8, inside_vertex: usize) -> ProjectedTriangle {
    let next_idx = (inside_vertex + 1) % 3;
    let other_idx = (inside_vertex + 2) % 3;
    let t1 = find_t(triangle, inside_vertex, next_idx, plane_to_clip);
    let t2 = find_t(triangle, inside_vertex, other_idx, plane_to_clip);

    ProjectedTriangle {
        view_verts: [
            ponderate_vec3(&triangle.view_verts[next_idx], &triangle.view_verts[inside_vertex], t1),
            ponderate_vec3(&triangle.view_verts[other_idx], &triangle.view_verts[inside_vertex], t2),
            triangle.view_verts[inside_vertex],
        ],
        clip_vertex: [
            ponderate_vec4(&triangle.clip_vertex[next_idx], &triangle.clip_vertex[inside_vertex], t1),
            ponderate_vec4(&triangle.clip_vertex[other_idx], &triangle.clip_vertex[inside_vertex], t2),
            triangle.clip_vertex[inside_vertex],
        ],
        view_normals: [
            ponderate_vec3(&triangle.view_normals[next_idx], &triangle.view_normals[inside_vertex], t1),
            ponderate_vec3(&triangle.view_normals[other_idx], &triangle.view_normals[inside_vertex], t2),
            triangle.view_normals[inside_vertex],
        ],
        uv_mapping: [
            ponderate_slice3(&triangle.uv_mapping[next_idx], &triangle.uv_mapping[inside_vertex], t1),
            ponderate_slice3(&triangle.uv_mapping[other_idx], &triangle.uv_mapping[inside_vertex], t2),
            triangle.uv_mapping[inside_vertex],
        ],
        light_intensity: [
            t1 * triangle.light_intensity[next_idx] + (1. - t1) * triangle.light_intensity[inside_vertex],
            t1 * triangle.light_intensity[other_idx] + (1. - t1) * triangle.light_intensity[inside_vertex],
            triangle.light_intensity[inside_vertex],
        ],
    }
}

pub fn clip_triangle_with_two_vertex_inside(triangle: &ProjectedTriangle, plane_to_clip: u8, outside_vertex: usize) -> (ProjectedTriangle, ProjectedTriangle) {
    let next_idx = ((outside_vertex + 1) % 3) as usize;
    let other_idx = ((outside_vertex + 2) % 3) as usize;
    let t1 = find_t(triangle, outside_vertex, next_idx, plane_to_clip);
    let t2 = find_t(triangle, outside_vertex, other_idx, plane_to_clip);

    let new_view_vert = ponderate_vec3(&triangle.view_verts[next_idx], &triangle.view_verts[outside_vertex], t1);
    let new_clip_vert = ponderate_vec4(&triangle.clip_vertex[next_idx], &triangle.clip_vertex[outside_vertex], t1);
    let new_view_normal = ponderate_vec3(&triangle.view_normals[next_idx], &triangle.view_normals[outside_vertex], t1);
    let new_uv_mapping = ponderate_slice3(&triangle.uv_mapping[next_idx], &triangle.uv_mapping[outside_vertex], t1);
    let new_light_intensity = triangle.light_intensity[next_idx] * t1 + (1. - t1) * triangle.light_intensity[outside_vertex];

    let triangle1 = ProjectedTriangle {
        view_verts: [new_view_vert, triangle.view_verts[next_idx], triangle.view_verts[other_idx]],
        clip_vertex: [new_clip_vert, triangle.clip_vertex[next_idx], triangle.clip_vertex[other_idx]],
        view_normals: [new_view_normal, triangle.view_normals[next_idx], triangle.view_normals[other_idx]],
        uv_mapping: [new_uv_mapping, triangle.uv_mapping[next_idx], triangle.uv_mapping[other_idx]],
        light_intensity: [new_light_intensity, triangle.light_intensity[next_idx], triangle.light_intensity[other_idx]],
    };
    let triangle2 = ProjectedTriangle {
        view_verts: [
            ponderate_vec3(&triangle.view_verts[other_idx], &triangle.view_verts[outside_vertex], t2),
            new_view_vert,
            triangle.view_verts[other_idx],
        ],
        clip_vertex: [
            ponderate_vec4(&triangle.clip_vertex[other_idx], &triangle.clip_vertex[outside_vertex], t2),
            new_clip_vert,
            triangle.clip_vertex[other_idx],
        ],
        view_normals: [
            ponderate_vec3(&triangle.view_normals[other_idx], &triangle.view_normals[outside_vertex], t2),
            new_view_normal,
            triangle.view_normals[other_idx],
        ],
        uv_mapping: [
            ponderate_slice3(&triangle.uv_mapping[other_idx], &triangle.uv_mapping[outside_vertex], t2),
            new_uv_mapping,
            triangle.uv_mapping[other_idx],
        ],
        light_intensity: [
            t2 * triangle.light_intensity[other_idx] + (1. - t2) * triangle.light_intensity[outside_vertex],
            new_light_intensity,
            triangle.light_intensity[other_idx],
        ],
    };

    return (triangle1, triangle2);
}

pub fn get_inside_plane_vertex(vert_codes: [u8; 3], plane: u8) -> usize {
    for i in 0..vert_codes.len() {
        let code = vert_codes[i];
        if code & plane == 0 {
            return i;
        }
    }
    panic!("get_inside_plane_vertex failed");
}

pub fn get_outside_plane_vertex(vert_codes: [u8; 3], plane: u8) -> usize {
    for i in 0..vert_codes.len() {
        let code = vert_codes[i];
        if code & plane != 0 {
            return i;
        }
    }

    panic!("get_outside_plane_vertex failed");
}

pub fn clip_triangle(triangle: ProjectedTriangle) -> Vec<Box<ProjectedTriangle>> {
    // https://gabrielgambetta.com/computer-graphics-from-scratch/11-clipping.html

    let base_vertex_0_code = out_code(&triangle, 0);
    let base_vertex_1_code = out_code(&triangle, 1);
    let base_vertex_2_code = out_code(&triangle, 2);
    if (base_vertex_0_code | base_vertex_1_code | base_vertex_2_code) == INSIDE {
        // Base check inside, doesn't need clipping
        return vec![Box::new(triangle)];
    }

    if (base_vertex_0_code & base_vertex_1_code & base_vertex_2_code) != INSIDE {
        // completly outside, doesn't need clipping
        return vec![];
    }

    let mut projection = vec![Box::new(triangle)];
    for plane_to_clip in [LEFT, RIGHT, TOP, BOTTOM, FRONT, BACK] {
        let mut plane_triangles = vec![];

        // Crop each triangle against the given plane
        for projected in projection {
            let vertex_0_code = out_code(&projected, 0);
            let vertex_1_code = out_code(&projected, 1);
            let vertex_2_code = out_code(&projected, 2);
            let vert_codes = [vertex_0_code, vertex_1_code, vertex_2_code];

            let mut inside_plane = 0;
            for code in vert_codes {
                if code & plane_to_clip == 0 {
                    inside_plane += 1;
                }
            }

            if inside_plane == 0 {
                // Trivial case reject
                continue;
            }

            if inside_plane == 1 {
                // crop single tringle
                let inside_vertex = get_inside_plane_vertex(vert_codes, plane_to_clip);
                let new_triangle = clip_triangle_with_one_vertex_inside(&projected, plane_to_clip, inside_vertex);
                plane_triangles.push(Box::new(new_triangle));
            }

            if inside_plane == 2 {
                // split in two triangles
                let outside_vertex = get_outside_plane_vertex(vert_codes, plane_to_clip);
                let (t1, t2) = clip_triangle_with_two_vertex_inside(&projected, plane_to_clip, outside_vertex);
                plane_triangles.push(Box::new(t1));
                plane_triangles.push(Box::new(t2));
            }

            if inside_plane == 3 {
                // keep the triangle as is
                plane_triangles.push(projected);
            }
        }

        projection = plane_triangles;
    }

    return projection;
}
