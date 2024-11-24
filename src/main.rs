use std::io::{Read, Write};
use std::str;
use std::time::Duration;
use std::{
    io::BufReader,
    net::{TcpListener, TcpStream},
    thread::spawn,
};
fn main() {
    let listener = TcpListener::bind("127.0.0.1:7878").unwrap();

    for stream in listener.incoming() {
        let stream = stream.unwrap();
        spawn(
            || match stream.set_read_timeout(Some(Duration::from_secs(3))) {
                Ok(_) => {
                    handle_connection(stream);
                }
                Err(e) => {
                    println!("Error setting read timeout: {:?}", e);
                }
            },
        );
    }
}

fn handle_connection(mut stream: TcpStream) {
    let buf_reader = BufReader::new(&mut stream);
    let mut header_done: bool = false;
    let mut headers: Vec<String> = Vec::new();
    let mut body: Vec<u8> = Vec::new();
    let mut remaining_content_length: i64 = 0;
    let mut current_line: Vec<u8> = Vec::new();

    for byte in buf_reader.bytes() {
        let byte = byte.unwrap();
        println!("Processing byte: {:?}", byte);
        current_line.push(byte);
        if !header_done && current_line.ends_with(&[13, 10]) {
            // Remove last two bytes
            current_line.pop();
            current_line.pop();
            let line = String::from_utf8(current_line.clone()).unwrap();
            current_line.clear();
            if !header_done && line == "" {
                println!("Found empty line");
                header_done = true;
            }
            if !header_done && line.starts_with("Content-Length") {
                let parts: Vec<&str> = line.split(":").collect();
                remaining_content_length = parts[1].trim().parse().unwrap();
            }
            println!("Found header: {:?}", line);
            headers.push(line);
            if header_done && remaining_content_length == 0 {
                break;
            }
        } else if header_done && remaining_content_length > 0 {
            remaining_content_length -= 1;
            body.push(byte.clone());
            if remaining_content_length == 0 {
                println!("Body done");
                break;
            }
        } else {
            println!("What is this?");
        }
    }

    println!("Headers: {:?}", headers);
    // println!("Body: {:?}", String::from_utf8(body.clone()).unwrap());
    println!("Body: {:?}", str::from_utf8(&body).unwrap());
    stream.write_all(b"HTTP/1.1 200 OK\r\n\r\n").unwrap();
}
