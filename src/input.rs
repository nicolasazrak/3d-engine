use crate::algebra::*;
use minifb::{Key, Window, WindowOptions};

pub fn handle_input(t: f32, window: &Window) -> (Vec3f, Vec2f) {
    let move_speed = 0.003;
    let rotation_speed = 0.005;

    let mut mov = Vec3f { x: 0., y: 0., z: 0. };
    let mut rot = Vec2f { x: 0., y: 0. };

    if window.is_key_down(Key::A) {
        mov.x -= move_speed * t;
    }
    if window.is_key_down(Key::S) {
        mov.z += move_speed * t;
    }
    if window.is_key_down(Key::Q) {
        mov.y -= move_speed * t;
    }
    if window.is_key_down(Key::E) {
        mov.y += move_speed * t;
    }
    if window.is_key_down(Key::D) {
        mov.x += move_speed * t;
    }
    if window.is_key_down(Key::W) {
        mov.z -= move_speed * t;
    }

    if window.is_key_down(Key::Up) {
        rot.y += rotation_speed * t;
    }
    if window.is_key_down(Key::Down) {
        rot.y -= rotation_speed * t;
    }
    if window.is_key_down(Key::Left) {
        rot.x += rotation_speed * t;
    }
    if window.is_key_down(Key::Right) {
        rot.x -= rotation_speed * t;
    }

    return (mov, rot);
}
