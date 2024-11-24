use std::io::{Read, Write};
use std::net::TcpStream;
use std::process;
#[test]
fn it_works() {
    // Connect to the server on localhost:7878
    let mut stream = TcpStream::connect("127.0.0.1:7878").unwrap_or_else(|e| {
        eprintln!("Failed to connect to server: {}", e);
        process::exit(1);
    });
    println!("Connected to server");

    // Send a message to the server
    if let Err(e) = stream.write_all(b"Hello, Server!") {
        eprintln!("Failed to send message to server: {}", e);
        process::exit(1);
    }
    println!("Sent message to server");

    // Read the response from the server
    let mut buffer = [0; 128];
    match stream.read(&mut buffer) {
        Ok(bytes_read) => {
            println!(
                "Received from server: {}",
                String::from_utf8_lossy(&buffer[..bytes_read])
            );
        }
        Err(e) => {
            eprintln!("Failed to read from server: {}", e);
            process::exit(1);
        }
    }
}
