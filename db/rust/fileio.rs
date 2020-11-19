use std::fs::File;
use std::io::prelude::*
use std::path::Path;

pub fn create(filename: str) {
  let path = Path::new(filename);
  // Open a file in write-only mode, returns 'io::Result<File>'
  let mut file = match File::create(&path) {
    Err(err) => panic!("Err creating file: {}", err),
    Ok(file) => file,
  };
  // Write
}