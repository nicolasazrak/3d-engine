extern crate image;
extern crate minifb;

// use image::png::PngEncoder;
// use image::ColorType::Rgba8;
// use std::fs::File;
use minifb::{Key, Window, WindowOptions};
use std::time::Instant;

mod algebra;
mod camera;
mod clip;
mod collision;
mod geometry;
mod input;
mod model;
mod player;
mod raster;
mod scene;
mod script;
mod shader;

pub fn main() {
    let mut scene = scene::Scene::new(500, 500);
    script::load_base_scenario(&mut scene);

    let mut window = Window::new("Test", scene.width as usize, scene.height as usize, WindowOptions::default()).unwrap_or_else(|e| {
        panic!("{}", e);
    });

    // Limit to max ~60 fps update rate
    window.limit_update_rate(Some(std::time::Duration::from_micros(16600 / 2)));

    let mut start = Instant::now();
    while window.is_open() && !window.is_key_down(Key::Escape) {
        let t = start.elapsed().as_millis() as f32;
        start = Instant::now();
        let (mut mov, rot) = input::handle_input(t, &window);
        scene.camera.rotate(rot.x, rot.y);

        let light_mov = scene.camera.transform_input(&algebra::vec3f(0., 0., -0.3));
        scene.light = algebra::plus(&scene.camera.get_position(), &light_mov);

        if !mov.is_zero() {
            mov = scene.camera.transform_input(&mov);
            scene.player.borrow_mut().handle_mov(mov, &scene.obstacles);
            scene.camera.update(&scene.player.borrow());
        }

        let buffer = scene.process_frame();

        window.update_with_buffer(&buffer, scene.width as usize, scene.height as usize).unwrap();
        let took = start.elapsed().as_millis();
        println!("Frame took: {}ms, FPS: {}", took, 1000. / t);
    }

    // let pixel_buffer = scene.process_frame();
    // let file = File::create("out.png")?;
    // let out = PngEncoder::new(file);
    // out.encode(&pixel_buffer, scene.width, scene.height, Rgba8)
    //     .expect("Couldn't write image!");
    // Ok(())
}
