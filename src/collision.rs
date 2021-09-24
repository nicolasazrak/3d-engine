use crate::algebra::*;
use crate::model::Model;

pub struct BoundingBox {
	pmin: Vec3f,
	pmax: Vec3f,
}

impl BoundingBox {
	pub fn test(&self, from: &Vec3f, to: &Vec3f, direction: &Vec3f) -> (bool, Vec3f, f32) {
		if to.x > self.pmin.x && to.x < self.pmax.x && to.y > self.pmin.y && to.y < self.pmax.y && to.z > self.pmin.z && to.z < self.pmax.z {
			let mut min_d = 9999999.;
			let mut normal = vec3f(0., 0., 0.);
			let mut collided = false;

			if from.z < self.pmin.z {
				let div_z_plane = dot_product(direction, &vec3f(0., 0., 1.));
				if div_z_plane > 0. {
					let numerator = dot_product(&minus(&self.pmin, from), &vec3f(0., 0., 1.));
					let d = numerator / div_z_plane;
					let intersection = Vec3f {
						x: from.x + d * direction.x,
						y: from.y + d * direction.y,
						z: from.z + d * direction.z,
					};

					if intersection.y > self.pmin.y && intersection.y < self.pmax.y && intersection.x > self.pmin.x && intersection.x < self.pmax.x {
						if d < min_d {
							min_d = d;
							normal = vec3f(0., 0., 1.);
							collided = true;
						}
					}
				}
			}
			if from.z > self.pmax.z {
				let div_z_plane = dot_product(direction, &vec3f(0., 0., 1.));
				if div_z_plane < 0. {
					let numerator = dot_product(&minus(&self.pmax, from), &vec3f(0., 0., 1.));
					let d = numerator / div_z_plane;
					let intersection = Vec3f {
						x: from.x + d * direction.x,
						y: from.y + d * direction.y,
						z: from.z + d * direction.z,
					};

					if intersection.y > self.pmin.y && intersection.y < self.pmax.y && intersection.x > self.pmin.x && intersection.x < self.pmax.x {
						if d < min_d {
							min_d = d;
							normal = vec3f(0., 0., -1.);
							collided = true;
						}
					}
				}
			}
			if from.x < self.pmin.x {
				let div_x_plane = dot_product(direction, &vec3f(1., 0., 0.));
				if div_x_plane > 0. {
					let numerator = dot_product(&minus(&self.pmin, from), &vec3f(1., 0., 0.));
					let d = numerator / div_x_plane;
					let intersection = Vec3f {
						x: from.x + d * direction.x,
						y: from.y + d * direction.y,
						z: from.z + d * direction.z,
					};

					if intersection.y > self.pmin.y && intersection.y < self.pmax.y && intersection.z > self.pmin.z && intersection.z < self.pmax.z {
						if d < min_d {
							min_d = d;
							normal = vec3f(-1., 0., 0.);
							collided = true;
						}
					}
				}
			}
			if from.x > self.pmin.x {
				let div_x_plane = dot_product(direction, &vec3f(1., 0., 0.));
				if div_x_plane < 0. {
					let numerator = dot_product(&minus(&self.pmax, from), &vec3f(1., 0., 0.));
					let d = numerator / div_x_plane;
					let intersection = Vec3f {
						x: from.x + d * direction.x,
						y: from.y + d * direction.y,
						z: from.z + d * direction.z,
					};

					if intersection.y > self.pmin.y
						&& intersection.y < self.pmax.y
						&& intersection.z > self.pmin.z
						&& intersection.z < self.pmax.z
						&& div_x_plane < 0.
					{
						if d < min_d {
							min_d = d;
							normal = vec3f(1., 0., 0.);
							collided = true;
						}
					}
				}
			}
			if from.y < self.pmin.y {
				let div_y_plane = dot_product(direction, &vec3f(0., 1., 0.));
				if div_y_plane > 0. {
					let numerator = dot_product(&minus(&self.pmax, from), &vec3f(0., 1., 0.));
					let d = numerator / div_y_plane;
					let intersection = Vec3f {
						x: from.x + d * direction.x,
						y: from.y + d * direction.y,
						z: from.z + d * direction.z,
					};

					if intersection.x > self.pmin.x && intersection.x < self.pmax.x && intersection.z > self.pmin.z && intersection.z < self.pmax.z {
						if d < min_d {
							min_d = d;
							normal = vec3f(0., -1., 0.);
							collided = true;
						}
					}
				}
			}
			if from.y > self.pmax.y {
				let div_y_plane = dot_product(direction, &vec3f(0., 1., 0.));
				if div_y_plane < 0. {
					let numerator = dot_product(&minus(&self.pmin, from), &vec3f(0., 1., 0.));
					let d = numerator / div_y_plane;
					let intersection = Vec3f {
						x: from.x + d * direction.x,
						y: from.y + d * direction.y,
						z: from.z + d * direction.z,
					};

					if intersection.x > self.pmin.x && intersection.x < self.pmax.x && intersection.z > self.pmin.z && intersection.z < self.pmax.z {
						if d < min_d {
							min_d = d;
							normal = vec3f(0., 1., 0.);
							collided = true;
						}
					}
				}
			}
			return (collided, normal, min_d);
		} else {
			return (false, vec3f(0., 0., 0.), 0.);
		}
	}
	pub fn from_model(model: &Model) -> BoundingBox {
		let mut pmin = Vec3f {
			x: 999999999.,
			y: 9999999999.,
			z: 99999999.,
		};
		let mut pmax = Vec3f {
			x: -999999999.,
			y: -9999999999.,
			z: -99999999.,
		};
		for triangle in &model.triangles {
			for vert in triangle.world_verts {
				pmax.x = pmax.x.max(vert.x) + 0.01;
				pmax.y = pmax.y.max(vert.y) + 0.01;
				pmax.z = pmax.z.max(vert.z) + 0.01;
				pmin.x = pmin.x.min(vert.x) - 0.01;
				pmin.y = pmin.y.min(vert.y) - 0.01;
				pmin.z = pmin.z.min(vert.z) - 0.01;
			}
		}

		BoundingBox { pmin: pmin, pmax: pmax }
	}
}
