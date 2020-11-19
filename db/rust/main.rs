use std::fs::File;
use std::io::prelude::*;
use std::path::Path;
use std::env;

fn main() -> std::io::Result<()> {
    let path = env::current_dir()?;
    println!("The current directory is {}", path.display());
    match env::current_exe() {
      Ok(exe_path) => println!("{}", exe_path.display()),
      Err(e) => println!("Failed: {}", e),
    };
    Ok(())
}