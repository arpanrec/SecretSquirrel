use std::io::{Read, Write};
use std::net::TcpListener;
use std::process;
fn main() {
    let listener = TcpListener::bind("127.0.0.1:7878").unwrap_or_else(|e| {
        println!("Failed to bind to address: {}", e);
        process::exit(1);
    });
    println!("Server listening on 127.0.0.1:7878");
    for stream in listener.incoming() {
        match stream {
            Ok(mut stream) => {
                println!("Client connected: {}", stream.peer_addr().unwrap());
                let mut buffer = [0; 128];
                match stream.read(&mut buffer) {
                    Ok(bytes_read) => {
                        println!(
                            "Received from client: {}",
                            String::from_utf8_lossy(&buffer[..bytes_read])
                        );
                    }
                    Err(e) => {
                        eprintln!("Failed to read from client: {}", e);
                        continue;
                    }
                }

                // Send a response to the client
                if let Err(e) = stream.write_all(b"Hello, Clientasfasfas123214124!") {
                    println!("Failed to write to client: {}", e);
                } else {
                    println!("Sent response to client");
                }
            }
            Err(e) => eprintln!("Connection failed: {}", e),
        }
    }
}
