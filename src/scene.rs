use crate::algebra::*;
use crate::camera::*;
use crate::collision::*;
use crate::model::*;
use crate::player::Player;
use crate::raster::*;
use std::cell::RefCell;
use std::rc::Rc;
use std::vec::Vec;

pub struct Scene {
    pub width: u32,
    pub height: u32,
    pub models: Vec<Rc<RefCell<Model>>>,
    pub camera: Box<dyn Camera>,
    pub light: Vec3f,
    pub obstacles: Vec<BoundingBox>,
    pub player: Rc<RefCell<Player>>,
}

impl Scene {
    pub fn new(width: u32, height: u32) -> Scene {
        let mut scene = Scene {
            width: width,
            height: height,
            models: vec![],
            camera: Box::new(FPSCamera::new()),
            light: vec3f(0., 0., 0.),
            obstacles: vec![],
            player: Rc::new(RefCell::new(Player::new())),
        };
        scene.camera.update(&scene.player.borrow());
        return scene;
    }

    pub fn make_pixel_buffer(&self) -> Vec<u32> {
        vec![0 as u32; (self.width * self.height) as usize]
    }

    pub fn make_z_buffer(&self) -> Vec<f32> {
        vec![-999999 as f32; (self.width * self.height) as usize]
    }

    pub fn process_frame(&mut self) -> Vec<u32> {
        self.camera.update_view_matrix();

        let mut pixel_buffer = self.make_pixel_buffer();
        let mut z_buffer = self.make_z_buffer();

        let mut drawn = 0;
        for i in 0..self.models.len() {
            let model = &mut self.models[i];
            self.camera.project_model(&mut model.borrow_mut(), &self.light);
            for projected in &model.borrow().projection {
                drawn += 1;
                draw_triangle(&model.borrow(), &projected, self.width, self.height, &mut pixel_buffer, &mut z_buffer);
            }
        }

        // println!("Drawn {} triangles", drawn);

        pixel_buffer
    }
}
