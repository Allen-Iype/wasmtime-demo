use std::collections::HashMap;
use std::env;
use std::fs;
use std::io::{Read, Write};
use wasm_bindgen::prelude::*;

struct Contract {
    state: HashMap<String, String>,
    state_file: String,
}

impl Contract {
    fn new(state_file: &str) -> Result<Self, std::io::Error> {
        let mut state = HashMap::new();
        if let Ok(contents) = fs::read_to_string(state_file) {
            for line in contents.lines() {
                let parts: Vec<&str> = line.splitn(2, '=').collect();
                if parts.len() == 2 {
                    state.insert(parts[0].to_owned(), parts[1].to_owned());
                }
            }
        }
        Ok(Self {
            state,
            state_file: state_file.to_owned(),
        })
    }

    fn get(&self, key: &str) -> Option<&String> {
        self.state.get(key)
    }

    fn set(&mut self, key: String, value: String) {
        self.state.insert(key, value);
    }

    fn save_state(&self) -> Result<(), std::io::Error> {
        
        let mut contents = String::new();
        for (key, value) in &self.state {
            contents.push_str(&format!("{}={}\n", key, value));
        }
        fs::write(&self.state_file, contents)
    }
}

impl Drop for Contract {
    fn drop(&mut self) {
        if let Err(err) = self.save_state() {
            eprintln!("Error saving contract state: {}", err);
        }
    }
}
#[wasm_bindgen]
pub fn write_state() {
    let state_file = "contract_state.env";
    let mut contract = match Contract::new(state_file) {
        Ok(contract) => contract,
        Err(err) => {
            eprintln!("Error loading contract state: {}", err);
            return;
        }
    };

    // Example usage
    contract.set("name".to_owned(), "Shaji".to_owned());
    contract.set("age".to_owned(), "30".to_owned());

    if let Some(name) = contract.get("name") {
        println!("Name: {}", name);
    }
    if let Some(age) = contract.get("age") {
        println!("Age: {}", age);
    }

}

fn main() {
    let state_file = "contract_state.env";
    let mut contract = match Contract::new(state_file) {
        Ok(contract) => contract,
        Err(err) => {
            eprintln!("Error loading contract state: {}", err);
            return;
        }
    };

    // Example usage
    contract.set("name".to_owned(), "Shaji".to_owned());
    contract.set("age".to_owned(), "30".to_owned());

    if let Some(name) = contract.get("name") {
        println!("Name: {}", name);
    }
    if let Some(age) = contract.get("age") {
        println!("Age: {}", age);
    }
}