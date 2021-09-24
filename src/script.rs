use crate::collision::BoundingBox;
use crate::geometry::*;
use crate::model::parse_model;
use crate::scene::Scene;
use crate::shader::*;
use std::cell::RefCell;
use std::rc::Rc;

pub fn load_base_scenario(scene: &mut Scene) {
    let red: Rc<dyn Shader> = Rc::new(FlatColor { r: 170, g: 30, b: 30 });
    let blue: Rc<dyn Shader> = Rc::new(FlatColor { r: 30, g: 30, b: 130 });
    let green: Rc<dyn Shader> = Rc::new(FlatColor { r: 30, g: 143, b: 23 });
    let purple: Rc<dyn Shader> = Rc::new(FlatColor { r: 114, g: 48, b: 191 });
    let grey: Rc<dyn Shader> = Rc::new(FlatColor { r: 100, g: 100, b: 100 });

    let scenario = [
        "XXXXXXXXXXXXXXXXXXXXXXXXX",
        "X          B     R      X",
        "X          B     R      X",
        "XRRRRRR    B     RRRR   X",
        "X          B            X",
        " X         BBBBBBBBBB   X",
        " X         B            X",
        " X    BBBBBB      G    X",
        " X                G    X",
        " X                G    X",
        " X    RRRRRGGGGGGGG    X",
        "X          G           X",
        "X          G           X",
        "XXXXXXXXXXXXXXXXXXXXXXXX",
    ];

    for y in 0..scenario.len() {
        let line = scenario[y as usize];
        for x in 0..line.len() {
            let code = line.chars().nth(x).unwrap();
            let mut caja = None;

            if code == 'X' {
                caja = Some(Rc::new(RefCell::new(new_cube(1., Rc::clone(&purple)))));
            } else if code == 'G' {
                caja = Some(Rc::new(RefCell::new(new_cube(1., Rc::clone(&green)))));
            } else if code == 'B' {
                caja = Some(Rc::new(RefCell::new(new_cube(1., Rc::clone(&blue)))));
            } else if code == 'R' {
                caja = Some(Rc::new(RefCell::new(new_cube(1., Rc::clone(&red)))));
            }

            match caja {
                None => {}
                Some(b) => {
                    b.borrow_mut().move_x((14. - (x as f32)) as f32);
                    b.borrow_mut().move_z((6. - (y as f32)) as f32);
                    scene.obstacles.push(BoundingBox::from_model(&b.borrow()));
                    scene.models.push(b);
                }
            }

            let piso = Rc::new(RefCell::new(new_cube(1., Rc::clone(&grey))));
            piso.borrow_mut().move_x((14. - (x as f32)) as f32);
            piso.borrow_mut().move_z((6. - (y as f32)) as f32);
            piso.borrow_mut().move_y(-1.);
            scene.models.push(piso);
        }
    }

    let head_texture = Rc::new(TextureShader::from_file("assets/head.texture.tga"));
    let mut head = parse_model("assets/head.obj", head_texture);
    head.move_z(-2.5);
    head.move_y(0.2);
    scene.models.push(Rc::new(RefCell::new(head)));
}
