use crate::algebra::*;
use crate::model::*;
use crate::shader::*;
use std::f32::consts::PI;
use std::rc::Rc;

pub fn new_xz_square(size: f32, shader: Rc<dyn Shader>) -> Model {
    let v0 = vec3f(-size / 2., 0., -size / 2.);
    let v1 = vec3f(size / 2., 0., -size / 2.);
    let v2 = vec3f(size / 2., 0., size / 2.);
    let v3 = vec3f(-size / 2., 0., size / 2.);

    let t0 = [v1, v0, v2];
    let t1 = [v3, v2, v0];

    let normal = vec3f(0., 1., 0.);

    let textt0v0 = [0., 0., 0.];
    let textt1v0 = [0., 0., 0.];
    let textt0v1 = [0.999, 0., 0.];
    let textt0v2 = [0.999, 0.999, 0.];
    let textt1v2 = [0.999, 0.999, 0.];
    let textt1v3 = [0., 0.999, 0.];

    let mut triangle1 = Triangle::new();
    triangle1.world_verts = t0;
    triangle1.normals = [normal, normal, normal];
    triangle1.uv_mapping = [textt0v1, textt0v0, textt0v2];

    let mut triangle2 = Triangle::new();
    triangle2.world_verts = t1;
    triangle2.normals = [normal, normal, normal];
    triangle2.uv_mapping = [textt1v3, textt1v2, textt1v0];

    return Model::new(vec![triangle1, triangle2], shader);
}

pub fn new_xy_square(size: f32, shader: Rc<dyn Shader>) -> Model {
    let pos = size / 2.;
    let neg = -size / 2.;

    let v0 = Vec3f { x: neg, y: pos, z: 0. };
    let v1 = Vec3f { x: neg, y: neg, z: 0. };
    let v2 = Vec3f { x: pos, y: neg, z: 0. };
    let v3 = Vec3f { x: pos, y: pos, z: 0. };

    let t0 = [v1, v2, v0];
    let t1 = [v3, v0, v2];

    let normal = Vec3f { x: 0., y: 0., z: 1. };

    let textt0v0 = [0., 0.999, 0.];
    let textt1v0 = [0., 0.999, 0.];
    let textt0v1 = [0., 0., 0.];
    let textt0v2 = [0.999, 0., 0.];
    let textt1v2 = [0.999, 0., 0.];
    let textt1v3 = [0.999, 0.999, 0.];

    let mut triangle1 = Triangle::new();
    triangle1.world_verts = t0;
    triangle1.normals = [normal, normal, normal];
    triangle1.uv_mapping = [textt0v1, textt0v2, textt0v0];

    let mut triangle2 = Triangle::new();
    triangle2.world_verts = t1;
    triangle2.normals = [normal, normal, normal];
    triangle2.uv_mapping = [textt1v3, textt1v0, textt1v2];

    return Model::new(vec![triangle1, triangle2], shader);
}

pub fn new_yz_square(size: f32, shader: Rc<dyn Shader>) -> Model {
    let pos = size / 2.;
    let neg = -size / 2.;

    let v0 = Vec3f { x: 0., y: pos, z: neg };
    let v1 = Vec3f { x: 0., y: neg, z: neg };
    let v2 = Vec3f { x: 0., y: neg, z: pos };
    let v3 = Vec3f { x: 0., y: pos, z: pos };

    let t0 = [v1, v2, v0];
    let t1 = [v3, v0, v2];

    let normal = Vec3f { x: 0., y: 0., z: 1. };

    let textt0v0 = [0., 0.999, 0.];
    let textt1v0 = [0., 0.999, 0.];
    let textt0v1 = [0., 0., 0.];
    let textt0v2 = [0.999, 0., 0.];
    let textt1v2 = [0.999, 0., 0.];
    let textt1v3 = [0.999, 0.999, 0.];

    let mut triangle1 = Triangle::new();
    triangle1.world_verts = t0;
    triangle1.normals = [normal, normal, normal];
    triangle1.uv_mapping = [textt0v1, textt0v2, textt0v0];

    let mut triangle2 = Triangle::new();
    triangle2.world_verts = t1;
    triangle2.normals = [normal, normal, normal];
    triangle2.uv_mapping = [textt1v3, textt1v0, textt1v2];

    return Model::new(vec![triangle1, triangle2], shader);
}

pub fn new_cube(size: f32, shader: Rc<dyn Shader>) -> Model {
    let mut bottom = new_xz_square(size, Rc::clone(&shader));
    bottom.rotate_x(PI).move_y(-size / 2.);

    let mut up = new_xz_square(size, Rc::clone(&shader));
    up.move_y(size / 2.);

    let mut right = new_yz_square(size, Rc::clone(&shader));
    right.rotate_y(PI).move_x(size / 2.);

    let mut left = new_yz_square(size, Rc::clone(&shader));
    left.move_x(-size / 2.);

    let mut back = new_xy_square(size, Rc::clone(&shader));
    back.rotate_y(PI).move_z(-size / 2.);

    let mut front = new_xy_square(size, Rc::clone(&shader));
    front.move_z(size / 2.);

    let mut triangles = vec![];
    triangles.append(&mut bottom.triangles);
    triangles.append(&mut up.triangles);
    triangles.append(&mut left.triangles);
    triangles.append(&mut right.triangles);
    triangles.append(&mut front.triangles);
    triangles.append(&mut back.triangles);

    return Model::new(triangles, shader);
}
