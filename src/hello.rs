use std::env;
use std::fs;
use std::io::{Read, Write};
use reqwest;

// fn process(input_fname: &str, output_fname: &str) -> Result<(), String> {
//     let mut input_file =
//         fs::File::open(input_fname).map_err(|err| format!("error opening input {}: {}", input_fname, err))?;
//     let mut contents = Vec::new();
//     input_file
//         .read_to_end(&mut contents)
//         .map_err(|err| format!("read error: {}", err))?;

//     let mut output_file = fs::File::create(output_fname)
//         .map_err(|err| format!("error opening output {}: {}", output_fname, err))?;
//     output_file
//         .write_all(&contents)
//         .map_err(|err| format!("write error: {}", err))
// }

// fn process1(input: &str)-> String{
//     //print input
//     let string = input.to_string();
//     // let mut file = fs::File::create("test.txt");
//     // file.write_all(&string).map_err(|err| format!("write error: {}", err));
//     println!("input: {}", input);

//     let host = "localhost";
//     let port = 20001;
//     let path = "/api/get-account-info";
//     let query_param = "did=bafybmid3slphcoalnccv6hei3k2dg34qdzltovzl6uarhs7an363e4apmy";

//     // let request = format!(
//     //     "GET {} HTTP/1.1\r\nHost: {}\r\n\r\n",
//     //     path, host
//     // );
//     let request = format!("GET {}?{} HTTP/1.1\r\nHost: {}:{}\r\n\r\n", path, query_param, host, port);

//     if let Ok(mut stream) = TcpStream::connect((host, port)) {
//         if let Err(err) = stream.write_all(request.as_bytes()) {
//             eprintln!("Failed to send request: {}", err);
//             return err.to_string();
//         }

//         let mut response = String::new();
//         if let Err(err) = stream.read_to_string(&mut response) {
//             eprintln!("Failed to read response: {}", err);
//             return err.to_string();
//         }

//         println!("Response:\n{}", response);
//     } else {
//         eprintln!("Failed to connect to the server");
//     }

    
    
    
    
//     return string;
// }



// #[tokio::main]
// async fn main() -> Result<(), Box<dyn std::error::Error>> {
//     let url = "http://api.example.com/api/endpoint";
//     let query_params = [("param1", "value1"), ("param2", "value2")];

//     let client = reqwest::Client::new();
//     let response = client.get(url).query(&query_params).send().await?;

//     if response.status().is_success() {
//         let body = response.text().await?;
//         println!("Response: {}", body);
//     } else {
//         println!("Request failed with status code: {}", response.status());
//     }

//     Ok(())
// }


fn call_api() {
    println!("Starting call_api function");
    let url = "http://localhost:20001/api/get-account-info";
    let query_params = [("did","bafybmid3slphcoalnccv6hei3k2dg34qdzltovzl6uarhs7an363e4apmy")];
    let client = reqwest::Client::new();
    let response = client.get(url).query(&query_params).send().await?;
    println!("Response is {}",response);
    if response.status().is_success() {
        let body = response.text().await?;
        println!("Response: {}", body);
    } else {
        println!("Request failed with status code: {}", response.status());
    }
}
fn main() {
    let args: Vec<String> = env::args().collect();
    let program = args[0].clone();
    println!("program : {:?}",program);
    //call an api
    // let url = "http://localhost:20001/api/get-account-info/did?=bafybmid3slphcoalnccv6hei3k2dg34qdzltovzl6uarhs7an363e4apmy";
    // let resp = reqwest::blocking::get(url).unwrap();


    // if args.len() < 3 {
    //     eprintln!("usage: {} <from> <to>", program);
    //     return;
    // }
    println!("args is : {:?}",args);
    call_api();
    // let result = process1(&args[1]);
    // println!("result: {}", result);
    // if let Err(err) = process1(&args[1]) {
    //     eprintln!("{}", err)
    // }
}