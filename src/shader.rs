extern crate image;

use crate::model::ProjectedTriangle;
use image::GenericImageView;

pub trait Shader {
    fn shade(&self, triangle: &ProjectedTriangle, coordinates: [f32; 3], z: f32) -> (u8, u8, u8);
}

pub struct FlatColor {
    pub r: u8,
    pub g: u8,
    pub b: u8,
}

pub struct TextureShader {
    height: f32,
    width: f32,
    data: Vec<u8>,
}

fn clamp(intensity: f32, color: u8) -> u8 {
    let res = intensity * (color as f32);
    if res > 255. {
        return 255;
    }
    return res as u8;
}

impl Shader for FlatColor {
    fn shade(&self, triangle: &ProjectedTriangle, coordinates: [f32; 3], _z: f32) -> (u8, u8, u8) {
        let l0 = coordinates[0];
        let l1 = coordinates[1];
        let l2 = coordinates[2];

        let intensity = l0 * triangle.light_intensity[0] + l1 * triangle.light_intensity[1] + l2 * triangle.light_intensity[2];
        if intensity < 0. {
            // Shoudln't be needed if there was occulsion culling or shadows ?
            return (0, 0, 0);
        } else {
            return (clamp(intensity, self.r), clamp(intensity, self.g), clamp(intensity, self.b));
        }
    }
}

impl TextureShader {
    pub fn from_file(file: &str) -> TextureShader {
        let i = image::open(file).unwrap();
        let i2 = i.flipv();
        let raw_pixels = i2.as_rgb8().unwrap().as_raw();
        TextureShader {
            height: i.height() as f32,
            width: i.width() as f32,
            data: raw_pixels.clone(),
        }
    }
}

impl Shader for TextureShader {
    fn shade(&self, triangle: &ProjectedTriangle, coordinates: [f32; 3], z: f32) -> (u8, u8, u8) {
        let l0 = coordinates[0];
        let l1 = coordinates[1];
        let l2 = coordinates[2];
        let t = triangle;

        let inv_z_0 = 1. / triangle.view_verts[0].z;
        let inv_z_1 = 1. / triangle.view_verts[1].z;
        let inv_z_2 = 1. / triangle.view_verts[2].z;
        let mut u = l0 * t.uv_mapping[0][0] * inv_z_0 + l1 * t.uv_mapping[1][0] * inv_z_1 + l2 * t.uv_mapping[2][0] * inv_z_2;
        let mut v = l0 * t.uv_mapping[0][1] * inv_z_0 + l1 * t.uv_mapping[1][1] * inv_z_1 + l2 * t.uv_mapping[2][1] * inv_z_2;
        u *= z;
        v *= z;

        if u < 0. || v < 0. {
            // TODO investigate why this happens
            return (0, 0, 0);
        }

        let x = ((u * self.width) as u32) % ((self.width) as u32);
        let y = ((v * self.height) as u32) % ((self.height) as u32);
        let pos = ((x + y * self.height as u32) * 3) as usize;
        let r = (self.data[pos] & 255) as u8;
        let g = ((self.data[pos + 1]) & 255) as u8;
        let b = ((self.data[pos + 2]) & 255) as u8;

        let intensity = l0 * triangle.light_intensity[0] + l1 * triangle.light_intensity[1] + l2 * triangle.light_intensity[2];

        return (clamp(intensity, r), clamp(intensity, g), clamp(intensity, b));
    }
}
