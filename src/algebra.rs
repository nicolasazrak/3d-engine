use std::cmp;

#[derive(Copy, Clone, Debug)]
pub struct Vec4f {
    pub x: f32,
    pub y: f32,
    pub z: f32,
    pub w: f32,
}

#[derive(Copy, Clone, Debug)]
pub struct Vec3f {
    pub x: f32,
    pub y: f32,
    pub z: f32,
}

#[derive(Copy, Clone, Debug)]
pub struct Vec2f {
    pub x: f32,
    pub y: f32,
}

#[derive(Copy, Clone, Debug)]
pub struct Vec2i {
    pub x: i32,
    pub y: i32,
}

impl Vec3f {
    pub fn is_zero(&self) -> bool {
        return self.x == 0. && self.y == 0. && self.z == 0.;
    }
}

pub fn vec4f(x: f32, y: f32, z: f32, w: f32) -> Vec4f {
    Vec4f { x: x, y: y, z: z, w: w }
}

pub fn vec3f(x: f32, y: f32, z: f32) -> Vec3f {
    Vec3f { x: x, y: y, z: z }
}

pub fn vec2i(x: i32, y: i32) -> Vec2i {
    Vec2i { x: x, y: y }
}

pub fn bounding_box(pts: &[Vec2i], minx: i32, maxx: i32, miny: i32, maxy: i32) -> (Vec2i, Vec2i) {
    let ptsminx = cmp::min(pts[0].x, cmp::min(pts[1].x, pts[2].x));
    let ptsmaxx = cmp::max(pts[0].x, cmp::max(pts[1].x, pts[2].x));

    let ptsminy = cmp::min(pts[0].y, cmp::min(pts[1].y, pts[2].y));
    let ptsmaxy = cmp::max(pts[0].y, cmp::max(pts[1].y, pts[2].y));

    let min_p = Vec2i {
        x: cmp::max(minx, cmp::min(ptsminx, maxx)),
        y: cmp::max(miny, cmp::min(ptsminy, maxy)),
    };
    let max_p = Vec2i {
        x: cmp::min(maxx, cmp::max(ptsmaxx, minx)),
        y: cmp::min(maxy, cmp::max(ptsmaxy, miny)),
    };

    return (min_p, max_p);
}

pub fn cross_product(a: &Vec3f, b: &Vec3f) -> Vec3f {
    let cx = a.y * b.z - a.z * b.y;
    let cy = a.z * b.x - a.x * b.z;
    let cz = a.x * b.y - a.y * b.x;
    return Vec3f { x: cx, y: cy, z: cz };
}

pub fn norm(a: &Vec3f) -> f32 {
    return (a.x * a.x + a.y * a.y + a.z * a.z).sqrt();
}

pub fn normalize(a: &Vec3f) -> Vec3f {
    let n = 1. / norm(a);
    return Vec3f {
        x: a.x * n,
        y: a.y * n,
        z: a.z * n,
    };
}

pub fn dot_product(a: &Vec3f, b: &Vec3f) -> f32 {
    return a.x * b.x + a.y * b.y + a.z * b.z;
}

pub fn minus(a: &Vec3f, b: &Vec3f) -> Vec3f {
    return Vec3f {
        x: a.x - b.x,
        y: a.y - b.y,
        z: a.z - b.z,
    };
}

pub fn plus(a: &Vec3f, b: &Vec3f) -> Vec3f {
    return Vec3f {
        x: a.x + b.x,
        y: a.y + b.y,
        z: a.z + b.z,
    };
}

pub fn _baycentric_coordinates(x: f32, y: f32, triangle: &[Vec3f]) -> Vec3f {
    let v0 = Vec3f {
        x: triangle[2].x - triangle[0].x,
        y: triangle[1].x - triangle[0].x,
        z: triangle[0].x - x,
    };
    let v1 = Vec3f {
        x: triangle[2].y - triangle[0].y,
        y: triangle[1].y - triangle[0].y,
        z: triangle[0].y - y,
    };

    let u = cross_product(&v0, &v1);

    if u.z > -1. && u.z < 1. {
        return Vec3f { x: -1., y: -1., z: -1. };
    }

    return Vec3f {
        x: 1. - (u.x + u.y) / u.z,
        y: u.y / u.z,
        z: u.x / u.z,
    };
}

pub fn orient_2d(a: &Vec2i, b: &Vec2i, x: i32, y: i32) -> i32 {
    return (b.x - a.x) * (y - a.y) - (b.y - a.y) * (x - a.x);
}

pub fn ponderate_vec3(vec1: &Vec3f, vec2: &Vec3f, t: f32) -> Vec3f {
    return Vec3f {
        x: vec1.x * t + (1. - t) * vec2.x,
        y: vec1.y * t + (1. - t) * vec2.y,
        z: vec1.z * t + (1. - t) * vec2.z,
    };
}

pub fn ponderate_vec4(vec1: &Vec4f, vec2: &Vec4f, t: f32) -> Vec4f {
    return Vec4f {
        x: vec1.x * t + (1. - t) * vec2.x,
        y: vec1.y * t + (1. - t) * vec2.y,
        z: vec1.z * t + (1. - t) * vec2.z,
        w: vec1.w * t + (1. - t) * vec2.w,
    };
}

pub fn ponderate_slice3(slice1: &[f32; 3], slice2: &[f32; 3], t: f32) -> [f32; 3] {
    return [
        slice1[0] * t + (1. - t) * slice2[0],
        slice1[1] * t + (1. - t) * slice2[1],
        slice1[2] * t + (1. - t) * slice2[2],
    ];
}

pub fn matmult(m: &[[f32; 4]; 4], vec: &Vec3f, h: f32) -> Vec3f {
    let res = matmult4(m, vec, h);
    let div = 1. / res.w;
    return Vec3f {
        x: res.x * div,
        y: res.y * div,
        z: res.z * div,
    };
}

pub fn matmult4(m: &[[f32; 4]; 4], vec: &Vec3f, h: f32) -> Vec4f {
    let x = m[0][0] * vec.x + m[1][0] * vec.y + m[2][0] * vec.z + m[3][0] * h;
    let y = m[0][1] * vec.x + m[1][1] * vec.y + m[2][1] * vec.z + m[3][1] * h;
    let z = m[0][2] * vec.x + m[1][2] * vec.y + m[2][2] * vec.z + m[3][2] * h;
    let w = m[0][3] * vec.x + m[1][3] * vec.y + m[2][3] * vec.z + m[3][3] * h;
    return Vec4f { x, y, z, w };
}

pub fn matmult4h(m: &[[f32; 4]; 4], vec: &Vec4f) -> Vec4f {
    let x = m[0][0] * vec.x + m[1][0] * vec.y + m[2][0] * vec.z + m[3][0] * vec.w;
    let y = m[0][1] * vec.x + m[1][1] * vec.y + m[2][1] * vec.z + m[3][1] * vec.w;
    let z = m[0][2] * vec.x + m[1][2] * vec.y + m[2][2] * vec.z + m[3][2] * vec.w;
    let w = m[0][3] * vec.x + m[1][3] * vec.y + m[2][3] * vec.z + m[3][3] * vec.w;
    return Vec4f { x, y, z, w };
}

pub fn inverse_transpose(dst: &mut [[f32; 4]; 4], src: &[[f32; 4]; 4]) {
    // https://semath.info/src/inverse-cofactor-ex4.html
    // https://stackoverflow.com/questions/33088577/symbolically-calculate-the-inverse-of-a-4-x-4-matrix-in-matlab
    let a_11 = src[0][0];
    let a_12 = src[1][0];
    let a_13 = src[2][0];
    let a_14 = src[3][0];

    let a_21 = src[0][1];
    let a_22 = src[1][1];
    let a_23 = src[2][1];
    let a_24 = src[3][1];

    let a_31 = src[0][2];
    let a_32 = src[1][2];
    let a_33 = src[2][2];
    let a_34 = src[3][2];

    let a_41 = src[0][3];
    let a_42 = src[1][3];
    let a_43 = src[2][3];
    let a_44 = src[3][3];

    dst[0][0] = (a_22 * a_33 * a_44 - a_22 * a_34 * a_43 - a_23 * a_32 * a_44 + a_23 * a_34 * a_42 + a_24 * a_32 * a_43 - a_24 * a_33 * a_42)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);
    dst[0][1] = -(a_12 * a_33 * a_44 - a_12 * a_34 * a_43 - a_13 * a_32 * a_44 + a_13 * a_34 * a_42 + a_14 * a_32 * a_43 - a_14 * a_33 * a_42)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);
    dst[0][2] = (a_12 * a_23 * a_44 - a_12 * a_24 * a_43 - a_13 * a_22 * a_44 + a_13 * a_24 * a_42 + a_14 * a_22 * a_43 - a_14 * a_23 * a_42)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);
    dst[0][3] = -(a_12 * a_23 * a_34 - a_12 * a_24 * a_33 - a_13 * a_22 * a_34 + a_13 * a_24 * a_32 + a_14 * a_22 * a_33 - a_14 * a_23 * a_32)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);

    dst[1][0] = -(a_21 * a_33 * a_44 - a_21 * a_34 * a_43 - a_23 * a_31 * a_44 + a_23 * a_34 * a_41 + a_24 * a_31 * a_43 - a_24 * a_33 * a_41)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);
    dst[1][1] = (a_11 * a_33 * a_44 - a_11 * a_34 * a_43 - a_13 * a_31 * a_44 + a_13 * a_34 * a_41 + a_14 * a_31 * a_43 - a_14 * a_33 * a_41)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);
    dst[1][2] = -(a_11 * a_23 * a_44 - a_11 * a_24 * a_43 - a_13 * a_21 * a_44 + a_13 * a_24 * a_41 + a_14 * a_21 * a_43 - a_14 * a_23 * a_41)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);
    dst[1][3] = (a_11 * a_23 * a_34 - a_11 * a_24 * a_33 - a_13 * a_21 * a_34 + a_13 * a_24 * a_31 + a_14 * a_21 * a_33 - a_14 * a_23 * a_31)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);

    dst[2][0] = (a_21 * a_32 * a_44 - a_21 * a_34 * a_42 - a_22 * a_31 * a_44 + a_22 * a_34 * a_41 + a_24 * a_31 * a_42 - a_24 * a_32 * a_41)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);
    dst[2][1] = -(a_11 * a_32 * a_44 - a_11 * a_34 * a_42 - a_12 * a_31 * a_44 + a_12 * a_34 * a_41 + a_14 * a_31 * a_42 - a_14 * a_32 * a_41)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);
    dst[2][2] = (a_11 * a_22 * a_44 - a_11 * a_24 * a_42 - a_12 * a_21 * a_44 + a_12 * a_24 * a_41 + a_14 * a_21 * a_42 - a_14 * a_22 * a_41)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);
    dst[2][3] = -(a_11 * a_22 * a_34 - a_11 * a_24 * a_32 - a_12 * a_21 * a_34 + a_12 * a_24 * a_31 + a_14 * a_21 * a_32 - a_14 * a_22 * a_31)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);

    dst[3][0] = -(a_21 * a_32 * a_43 - a_21 * a_33 * a_42 - a_22 * a_31 * a_43 + a_22 * a_33 * a_41 + a_23 * a_31 * a_42 - a_23 * a_32 * a_41)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);
    dst[3][1] = (a_11 * a_32 * a_43 - a_11 * a_33 * a_42 - a_12 * a_31 * a_43 + a_12 * a_33 * a_41 + a_13 * a_31 * a_42 - a_13 * a_32 * a_41)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);
    dst[3][2] = -(a_11 * a_22 * a_43 - a_11 * a_23 * a_42 - a_12 * a_21 * a_43 + a_12 * a_23 * a_41 + a_13 * a_21 * a_42 - a_13 * a_22 * a_41)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);
    dst[3][3] = (a_11 * a_22 * a_33 - a_11 * a_23 * a_32 - a_12 * a_21 * a_33 + a_12 * a_23 * a_31 + a_13 * a_21 * a_32 - a_13 * a_22 * a_31)
        / (a_11 * a_22 * a_33 * a_44 - a_11 * a_22 * a_34 * a_43 - a_11 * a_23 * a_32 * a_44 + a_11 * a_23 * a_34 * a_42 + a_11 * a_24 * a_32 * a_43
            - a_11 * a_24 * a_33 * a_42
            - a_12 * a_21 * a_33 * a_44
            + a_12 * a_21 * a_34 * a_43
            + a_12 * a_23 * a_31 * a_44
            - a_12 * a_23 * a_34 * a_41
            - a_12 * a_24 * a_31 * a_43
            + a_12 * a_24 * a_33 * a_41
            + a_13 * a_21 * a_32 * a_44
            - a_13 * a_21 * a_34 * a_42
            - a_13 * a_22 * a_31 * a_44
            + a_13 * a_22 * a_34 * a_41
            + a_13 * a_24 * a_31 * a_42
            - a_13 * a_24 * a_32 * a_41
            - a_14 * a_21 * a_32 * a_43
            + a_14 * a_21 * a_33 * a_42
            + a_14 * a_22 * a_31 * a_43
            - a_14 * a_22 * a_33 * a_41
            - a_14 * a_23 * a_31 * a_42
            + a_14 * a_23 * a_32 * a_41);
}
