use crate::algebra::*;
use crate::collision::BoundingBox;

pub struct Player {
    position: Vec3f,
}

impl Player {
    pub fn new() -> Player {
        return Player { position: vec3f(0., 0., 4.) };
    }
    pub fn get_position(&self) -> Vec3f {
        self.position
    }
    pub fn handle_mov(&mut self, initial_mov: Vec3f, obstacles: &Vec<BoundingBox>) {
        let mut mov = initial_mov;
        let mut collided = true;
        while collided {
            collided = false;
            let mut min_d = 999999.;
            let mut new_norm = vec3f(0., 0., 0.);
            let dst = plus(&mov, &self.position);

            for obstacle in obstacles {
                let (c, norm, d) = obstacle.test(&self.position, &dst, &mov);
                if c && d < min_d {
                    min_d = d;
                    new_norm = norm;
                    collided = true;
                }
            }

            if collided {
                mov = Vec3f {
                    x: mov.x - new_norm.x.abs() * mov.x,
                    y: mov.y - new_norm.y.abs() * mov.y,
                    z: mov.z - new_norm.z.abs() * mov.z,
                }
            }
        }
        self.position = plus(&self.position, &mov);
    }
}
