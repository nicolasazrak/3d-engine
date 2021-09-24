use crate::algebra::*;
use crate::clip::*;
use crate::model::*;
use crate::player::Player;
use std::f32::consts::PI;
use std::fmt::Debug;

pub trait Camera: Debug {
    fn project_model(&self, model: &mut Model, light: &Vec3f);
    fn update(&mut self, player: &Player);
    fn rotate(&mut self, yaw: f32, pitch: f32);
    fn transform_input(&self, mov: &Vec3f) -> Vec3f;
    fn get_position(&self) -> Vec3f;
    fn update_view_matrix(&mut self);
}

#[derive(Debug)]
pub struct FPSCamera {
    position: Vec3f,
    pitch: f32,
    yaw: f32,
    view_matrix: [[f32; 4]; 4],
    normal_matrix: [[f32; 4]; 4],
    projection_matrix: [[f32; 4]; 4],
}

pub fn build_projection_matrix() -> [[f32; 4]; 4] {
    let near_plane = 0.1;
    let far_plane = 50.;

    let left_plane = -0.1;
    let right_plane = 0.1;

    let top_plane = 0.1;
    let bottom_plane = -0.1;

    let fov_x = PI / 2.;
    let fov_y = PI / 2.;

    let use_open_gl_matrix = true;

    // http://www.songho.ca/opengl/gl_projection_matrix.html
    if use_open_gl_matrix {
        [
            [(2. * near_plane) / (right_plane - left_plane), 0., 0., 0.],
            [0., (2. * near_plane) / (top_plane - bottom_plane), 0., 0.],
            [
                (right_plane + left_plane) / (right_plane - left_plane),
                (top_plane + bottom_plane) / (top_plane - bottom_plane),
                -((far_plane + near_plane) / (far_plane - near_plane)),
                -1.,
            ],
            [0., 0., -((2. * far_plane * near_plane) / (far_plane - near_plane)), 0.],
        ]
    } else {
        [
            [1. / (fov_x / 2.).tan(), 0., 0., 0.],
            [0., 1. / (fov_y / 2.).tan(), 0., 0.],
            [0., 0., -(far_plane / (far_plane - near_plane)), -1.],
            [0., 0., -((far_plane * near_plane) / (far_plane - near_plane)), 0.],
        ]
    }
}

/** Look at camera */

/*

pub struct LookAtCamera {
    position: Vec3f,
    target: Vec3f,
    angle: f32,
    view_matrix: [[f32; 4]; 4],
    normal_matrix: [[f32; 4]; 4],
    projection_matrix: [[f32; 4]; 4],
}

impl LookAtCamera {
    pub fn new() -> LookAtCamera {
        LookAtCamera{
            position:         Vec3f{x: 0., y: 0., z: 4.},
            target:           Vec3f{x: 0., y: 0., z: -1.},
            angle:            0.,
            view_matrix:       [[0.;4]; 4],
            normal_matrix:     [[0.;4]; 4],
            projection_matrix: build_projection_matrix(),
        }
    }
}

impl Camera for LookAtCamera {
    pub fn (cam *LookAtCamera) move(move Vec3f) {
        cam.position.x += move.x
        cam.position.y += move.y
        cam.position.z += move.z
    }
    pub fn (cam *LookAtCamera) update_view_matrix() {
        // https://www.3dgep.com/understanding-the-view-matrix/
        // Look at camera
        // Two possible targets
        cam.target = Vec3f{0, 0, 0}
        // cam.target = Vec3f{cam.position.x, cam.position.y, cam.position.z - 1}
        let zaxis = normalize(minus(cam.position, cam.target));
        let xaxis = normalize(crossProduct(Vec3f{0, 1, 0}, zaxis));
        let yaxis = crossProduct(zaxis, xaxis);
        cam.view_matrix[0][0] = xaxis.x
        cam.view_matrix[0][1] = yaxis.x
        cam.view_matrix[0][2] = zaxis.x
        cam.view_matrix[0][3] = 0
        cam.view_matrix[1][0] = xaxis.y
        cam.view_matrix[1][1] = yaxis.y
        cam.view_matrix[1][2] = zaxis.y
        cam.view_matrix[1][3] = 0
        cam.view_matrix[2][0] = xaxis.z
        cam.view_matrix[2][1] = yaxis.z
        cam.view_matrix[2][2] = zaxis.z
        cam.view_matrix[2][3] = 0
        cam.view_matrix[3][0] = -dotProduct(xaxis, cam.position)
        cam.view_matrix[3][1] = -dotProduct(yaxis, cam.position)
        cam.view_matrix[3][2] = -dotProduct(zaxis, cam.position)
        cam.view_matrix[3][3] = 1
        inverseTranspose(&cam.normal_matrix, cam.view_matrix)
    }
    pub fn (cam *LookAtCamera) rotate(yaw, pitch float64) {
        // Not supported...
    }
    pub fn (cam *LookAtCamera) getPosition() Vec3f {
        return cam.position
    }
    pub fn (cam *LookAtCamera) transformInput(inputMove Vec3f) Vec3f {
        return inputMove
    }
    pub fn (cam *LookAtCamera) project(scene *Scene) {
        cam.update_view_matrix()
        scene.projectedLight = matmult(cam.view_matrix, scene.lightPosition, 1)
        let for _, model = range scene.models {;
            let projection = []*ProjectedTriangle{};
            let for _, triangle = range model.triangles {;
                projection = append(projection, projectTriangle(triangle, cam.view_matrix, cam.normal_matrix, cam.projection_matrix, scene.projectedLight)...)
            }
            model.projection = projection
        }
    }
}

*/

/** FPS camera */

impl FPSCamera {
    pub fn new() -> FPSCamera {
        FPSCamera {
            position: Vec3f { x: 0., y: 0., z: 0. },
            pitch: 0.,
            yaw: 0.,
            view_matrix: [[0.; 4]; 4],
            normal_matrix: [[0.; 4]; 4],
            projection_matrix: build_projection_matrix(),
        }
    }
}

impl Camera for FPSCamera {
    fn project_model(&self, model: &mut Model, light: &Vec3f) {
        model.projection = vec![];
        for triangle in &model.triangles {
            let mut t_projection = project_triangle(&triangle, &self.view_matrix, &self.normal_matrix, &self.projection_matrix, light);
            model.projection.append(&mut t_projection);
        }
    }

    fn transform_input(&self, mov: &Vec3f) -> Vec3f {
        return Vec3f {
            x: mov.z * (self.yaw.sin()) + mov.x * (self.yaw + PI / 2.).sin(),
            y: mov.y,
            z: mov.z * (self.yaw.cos()) + mov.x * (self.yaw + PI / 2.).cos(),
        };
    }
    fn update(&mut self, player: &Player) {
        self.position = player.get_position();
    }
    fn rotate(&mut self, yaw: f32, pitch: f32) {
        self.pitch += pitch;
        self.yaw += yaw;
    }
    fn update_view_matrix(&mut self) {
        // I assume the values are already converted to radians.
        let cos_pitch = self.pitch.cos();
        let sin_pitch = self.pitch.sin();
        let cos_yaw = self.yaw.cos();
        let sin_yaw = self.yaw.sin();
        let xaxis = vec3f(cos_yaw, 0., -sin_yaw);
        let yaxis = vec3f(sin_yaw * sin_pitch, cos_pitch, cos_yaw * sin_pitch);
        let zaxis = vec3f(sin_yaw * cos_pitch, -sin_pitch, cos_pitch * cos_yaw);
        self.view_matrix[0][0] = xaxis.x;
        self.view_matrix[0][1] = yaxis.x;
        self.view_matrix[0][2] = zaxis.x;
        self.view_matrix[0][3] = 0.;
        self.view_matrix[1][0] = xaxis.y;
        self.view_matrix[1][1] = yaxis.y;
        self.view_matrix[1][2] = zaxis.y;
        self.view_matrix[1][3] = 0.;
        self.view_matrix[2][0] = xaxis.z;
        self.view_matrix[2][1] = yaxis.z;
        self.view_matrix[2][2] = zaxis.z;
        self.view_matrix[2][3] = 0.;
        self.view_matrix[3][0] = -dot_product(&xaxis, &self.position);
        self.view_matrix[3][1] = -dot_product(&yaxis, &self.position);
        self.view_matrix[3][2] = -dot_product(&zaxis, &self.position);
        self.view_matrix[3][3] = 1.;
        inverse_transpose(&mut self.normal_matrix, &self.view_matrix)
    }
    fn get_position(&self) -> Vec3f {
        return self.position;
    }
}

/** General method */

pub fn project_triangle(
    original_triangle: &Triangle,
    view_matrix: &[[f32; 4]; 4],
    normal_matrix: &[[f32; 4]; 4],
    projection_matrix: &[[f32; 4]; 4],
    light: &Vec3f,
) -> Vec<Box<ProjectedTriangle>> {
    let mut projection = ProjectedTriangle::new();

    for i in 0..3 {
        let view = matmult4(view_matrix, &original_triangle.world_verts[i], 1.);
        let clip = matmult4h(projection_matrix, &view);
        let normal = normalize(&matmult(
            &normal_matrix,
            &original_triangle.normals[i],
            1., /* This should be 0. Why do I need to make it 1? */
        ));

        projection.clip_vertex[i] = clip;
        projection.view_verts[i].x = view.x / view.w;
        projection.view_verts[i].y = view.y / view.w;
        projection.view_verts[i].z = view.z / view.w;
        projection.view_normals[i] = normal;
        projection.uv_mapping = original_triangle.uv_mapping;

        // TODO only calculate if in frustrum
        // TODO use fast sqrt
        projection.light_intensity[i] = 1. / (norm(&minus(&original_triangle.world_verts[i], light)).sqrt());
    }

    // TODO add backface culling
    let clipped_projection = clip_triangle(projection);
    clipped_projection
}
