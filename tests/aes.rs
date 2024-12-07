use secretsquirrel::{aes256_dec, aes256_enc};

#[tokio::test]
async fn aes_test() {
    let key_str_base64 = "bCpu2ln1ivhkBlo1iYWfewMdi4yvQHEDnmClTj0ZNPU=";
    let iv_str_base64 = "+0Vfhn16YpMKYQNOvnP/AA==";
    let plaintext = "hello world! this is my plaintext.";
    let encrypted_text = "d618sNKZ9ouOIn4M5IiIanT5T14cJTJMxJ0d9xmo/hRf+TtuHB6G6tIkzq4viTSo";
    let enc = aes256_enc(key_str_base64, iv_str_base64, plaintext);
    println!("{:?}", enc);
    assert_eq!(enc, encrypted_text);
    let dec = aes256_dec(key_str_base64, iv_str_base64, &enc);
    println!("{:?}", dec);
    assert_eq!(dec, plaintext);
}
