use crate::algebra::*;
use crate::shader::*;
use std::fs;
use std::rc::Rc;

pub struct Model {
    pub triangles: Vec<Triangle>,
    pub projection: Vec<Box<ProjectedTriangle>>,
    pub shader: Rc<dyn Shader>,
}

#[derive(Debug)]
pub struct Triangle {
    pub world_verts: [Vec3f; 3], // world space
    pub normals: [Vec3f; 3],
    pub uv_mapping: [[f32; 3]; 3],
}

#[derive(Debug)]
pub struct ProjectedTriangle {
    pub view_verts: [Vec3f; 3], // view space relative to camera
    pub clip_vertex: [Vec4f; 3],
    pub view_normals: [Vec3f; 3], // view/camera space
    pub uv_mapping: [[f32; 3]; 3],
    pub light_intensity: [f32; 3],
}

impl Triangle {
    pub fn new() -> Triangle {
        Triangle {
            world_verts: [vec3f(0., 0., 0.); 3],
            normals: [vec3f(0., 0., 0.); 3],
            uv_mapping: [[0.; 3]; 3],
        }
    }
}

impl ProjectedTriangle {
    pub fn new() -> ProjectedTriangle {
        ProjectedTriangle {
            view_verts: [vec3f(0., 0., 0.); 3],
            view_normals: [vec3f(0., 0., 0.); 3],
            clip_vertex: [vec4f(0., 0., 0., 0.); 3],
            uv_mapping: [[0.; 3]; 3],
            light_intensity: [0.; 3],
        }
    }
}

pub fn parse_model(path: &str, shader: Rc<dyn Shader>) -> Model {
    let content = fs::read_to_string(path).unwrap();

    let mut triangles = vec![];
    let mut vertex = vec![];
    let mut normals = vec![];
    let mut textures = vec![];

    for line in content.split("\n") {
        if line.starts_with("v ") {
            let mut splitted = line.split_whitespace();
            splitted.next().unwrap();
            let x = splitted.next().unwrap().parse::<f32>().unwrap();
            let y = splitted.next().unwrap().parse::<f32>().unwrap();
            let z = splitted.next().unwrap().parse::<f32>().unwrap();
            vertex.push(vec3f(x, y, z));
        }

        if line.starts_with("vn ") {
            let mut splitted = line.split_whitespace();
            splitted.next().unwrap();
            let x = splitted.next().unwrap().parse::<f32>().unwrap();
            let y = splitted.next().unwrap().parse::<f32>().unwrap();
            let z = splitted.next().unwrap().parse::<f32>().unwrap();
            normals.push(vec3f(x, y, z));
        }

        if line.starts_with("vt ") {
            let mut splitted = line.split_whitespace();
            splitted.next().unwrap();
            let u = splitted.next().unwrap().parse::<f32>().unwrap();
            let v = splitted.next().unwrap().parse::<f32>().unwrap();
            let w = splitted.next().unwrap().parse::<f32>().unwrap();
            textures.push([u, v, w]);
        }

        if line.starts_with("f ") {
            let mut splitted = line.split_whitespace();
            splitted.next().unwrap();
            let part1 = splitted.next().unwrap();
            let part2 = splitted.next().unwrap();
            let part3 = splitted.next().unwrap();

            let mut part_1_splitted = part1.split("/");
            let mut part_2_splitted = part2.split("/");
            let mut part_3_splitted = part3.split("/");

            let vertex_idx_1 = (part_1_splitted.next().unwrap().parse::<i32>().unwrap() - 1) as usize;
            let vertex_idx_2 = (part_2_splitted.next().unwrap().parse::<i32>().unwrap() - 1) as usize;
            let vertex_idx_3 = (part_3_splitted.next().unwrap().parse::<i32>().unwrap() - 1) as usize;

            let texture_idx_1 = (part_1_splitted.next().unwrap().parse::<i32>().unwrap() - 1) as usize;
            let texture_idx_2 = (part_2_splitted.next().unwrap().parse::<i32>().unwrap() - 1) as usize;
            let texture_idx_3 = (part_3_splitted.next().unwrap().parse::<i32>().unwrap() - 1) as usize;

            let normal_idx_1 = (part_1_splitted.next().unwrap().parse::<i32>().unwrap() - 1) as usize;
            let normal_idx_2 = (part_2_splitted.next().unwrap().parse::<i32>().unwrap() - 1) as usize;
            let normal_idx_3 = (part_3_splitted.next().unwrap().parse::<i32>().unwrap() - 1) as usize;

            let mut triangle = Triangle::new();
            triangle.world_verts = [vertex[vertex_idx_1], vertex[vertex_idx_2], vertex[vertex_idx_3]];
            triangle.normals = [normals[normal_idx_1], normals[normal_idx_2], normals[normal_idx_3]];
            triangle.uv_mapping = [textures[texture_idx_1], textures[texture_idx_2], textures[texture_idx_3]];

            triangles.push(triangle);
        }
    }

    return Model {
        triangles: triangles,
        projection: vec![],
        shader: shader,
    };
}

impl Model {
    pub fn new(triangles: Vec<Triangle>, shader: Rc<dyn Shader>) -> Model {
        Model {
            triangles: triangles,
            projection: vec![],
            shader: shader,
        }
    }

    pub fn move_x(&mut self, x: f32) -> &mut Model {
        for idx in 0..self.triangles.len() {
            for vert_idx in [0, 1, 2] {
                self.triangles[idx].world_verts[vert_idx].x += x
            }
        }
        self
    }

    pub fn move_y(&mut self, y: f32) -> &mut Model {
        for idx in 0..self.triangles.len() {
            for vert_idx in [0, 1, 2] {
                self.triangles[idx].world_verts[vert_idx].y += y
            }
        }
        self
    }

    pub fn move_z(&mut self, z: f32) -> &mut Model {
        for idx in 0..self.triangles.len() {
            for vert_idx in [0, 1, 2] {
                self.triangles[idx].world_verts[vert_idx].z += z
            }
        }
        self
    }

    pub fn rotate_y(&mut self, v: f32) -> &mut Model {
        // TODO rotate normal
        for idx in 0..self.triangles.len() {
            for vert_idx in [0, 1, 2] {
                let x = self.triangles[idx].world_verts[vert_idx].x;
                let y = self.triangles[idx].world_verts[vert_idx].y;
                let z = self.triangles[idx].world_verts[vert_idx].z;
                self.triangles[idx].world_verts[vert_idx].x = x * v.cos() + z * v.sin();
                self.triangles[idx].world_verts[vert_idx].y = y;
                self.triangles[idx].world_verts[vert_idx].z = x * v.sin() + z * v.cos();
            }
        }
        self
    }

    pub fn rotate_x(&mut self, v: f32) -> &mut Model {
        // TODO rotate normal
        for idx in 0..self.triangles.len() {
            for vert_idx in [0, 1, 2] {
                let x = self.triangles[idx].world_verts[vert_idx].x;
                let y = self.triangles[idx].world_verts[vert_idx].y;
                let z = self.triangles[idx].world_verts[vert_idx].z;
                self.triangles[idx].world_verts[vert_idx].x = x;
                self.triangles[idx].world_verts[vert_idx].y = y * v.cos() + z * v.sin();
                self.triangles[idx].world_verts[vert_idx].z = y * -v.sin() + z * v.cos();
            }
        }
        self
    }
    pub fn rotate_z(&mut self, v: f32) -> &mut Model {
        // TODO rotate normal
        for idx in 0..self.triangles.len() {
            for vert_idx in [0, 1, 2] {
                let x = self.triangles[idx].world_verts[vert_idx].x;
                let y = self.triangles[idx].world_verts[vert_idx].y;
                let z = self.triangles[idx].world_verts[vert_idx].z;
                self.triangles[idx].world_verts[vert_idx].x = y * v.sin() + x.cos();
                self.triangles[idx].world_verts[vert_idx].y = y * v.cos() + x.sin();
                self.triangles[idx].world_verts[vert_idx].z = z;
            }
        }
        self
    }
    /*
    pub fn scale(&mut self, x: f32, y: f32, z : f32) -> Model {
        for _, triangle := range model.triangles {
            for i := range triangle.worldVerts {
                triangle.worldVerts[i].x *= x
                triangle.worldVerts[i].y *= y
                triangle.worldVerts[i].z *= z
            }
        }
        return self
    }
    pub fn scaleUV(&mut self, u, v : f32) -> Model {
        for _, triangle := range model.triangles {
            for t := range triangle.uvMapping {
                triangle.uvMapping[t][0] *= u
                triangle.uvMapping[t][1] *= v
            }
        }
        return self
    }
    */
}
