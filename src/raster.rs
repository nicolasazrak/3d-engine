use crate::algebra::*;
use crate::model::*;

pub fn draw_triangle(model: &Model, triangle: &ProjectedTriangle, width: u32, height: u32, pixel_buffer: &mut [u32], z_buffer: &mut [f32]) {
    let f_width = width as f32;
    let f_height = height as f32;

    let v0 = vec2i(
        (((triangle.clip_vertex[0].x / triangle.clip_vertex[0].w) + 1.) * (f_width - 1.) * 0.5).round() as i32,
        (((triangle.clip_vertex[0].y / triangle.clip_vertex[0].w) + 1.) * (f_height - 1.) * 0.5).round() as i32,
    );
    let v1 = vec2i(
        (((triangle.clip_vertex[1].x / triangle.clip_vertex[1].w) + 1.) * (f_width - 1.) * 0.5).round() as i32,
        (((triangle.clip_vertex[1].y / triangle.clip_vertex[1].w) + 1.) * (f_height - 1.) * 0.5).round() as i32,
    );
    let v2 = vec2i(
        (((triangle.clip_vertex[2].x / triangle.clip_vertex[2].w) + 1.) * (f_width - 1.) * 0.5).round() as i32,
        (((triangle.clip_vertex[2].y / triangle.clip_vertex[2].w) + 1.) * (f_height - 1.) * 0.5).round() as i32,
    );

    let pts = [v0, v1, v2];

    let (minbbox, maxbbox) = bounding_box(&pts, 0, (width - 1) as i32, 0, (height - 1) as i32);
    if minbbox.x >= maxbbox.x || minbbox.y >= maxbbox.y {
        // pseudo frustrum culling
        return;
    }

    let area = 1. / (orient_2d(&pts[0], &pts[1], pts[2].x, pts[2].y) as f32);
    if area <= 0. {
        // pseudo backface culling
        return;
    }

    let a_01 = v0.y - v1.y;
    let b_01 = v1.x - v0.x;
    let a_12 = v1.y - v2.y;
    let b_12 = v2.x - v1.x;
    let a_20 = v2.y - v0.y;
    let b_20 = v0.x - v2.x;

    let mut w0_row = orient_2d(&v1, &v2, minbbox.x, minbbox.y);
    let mut w1_row = orient_2d(&v2, &v0, minbbox.x, minbbox.y);
    let mut w2_row = orient_2d(&v0, &v1, minbbox.x, minbbox.y);

    // scene.trianglesDrawn++

    let inv_z0 = 1. / triangle.view_verts[0].z;
    let inv_z1 = 1. / triangle.view_verts[1].z;
    let inv_z2 = 1. / triangle.view_verts[2].z;

    // https://fgiesen.wordpress.com/2013/02/10/optimizing-the-basic-rasterizer/
    for y in minbbox.y..maxbbox.y {
        // Barycentric coordinates at start of row
        let mut w0 = w0_row;
        let mut w1 = w1_row;
        let mut w2 = w2_row;

        for x in minbbox.x..maxbbox.x {
            if (w0 | w1 | w2) >= 0 {
                let l0 = (w0 as f32) * area;
                let l1 = (w1 as f32) * area;
                let l2 = (w2 as f32) * area;

                // Should the z-buffer use the ndc value ??
                let z_pos = 1. / (l0 * inv_z0 + l1 * inv_z1 + l2 * inv_z2);
                let idx = (x + ((height as i32) - y - 1) * (width as i32)) as usize;

                if z_pos < 0. && z_pos > z_buffer[idx] {
                    z_buffer[idx] = z_pos;
                    let (r, g, b) = model.shader.shade(triangle, [l0, l1, l2], z_pos);
                    set_buffer(pixel_buffer, idx, r, g, b);
                }
            }

            // One step to the right
            w0 += a_12;
            w1 += a_20;
            w2 += a_01;
        }

        // One row step
        w0_row += b_12;
        w1_row += b_20;
        w2_row += b_01;
    }
}

fn set_buffer(pixel_buffer: &mut [u32], pixel_buffer_idx: usize, r: u8, g: u8, b: u8) {
    pixel_buffer[pixel_buffer_idx] = ((r as u32) << 16) | ((g as u32) << 8) | (b as u32);
}
