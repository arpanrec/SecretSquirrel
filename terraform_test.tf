terraform {
  backend "http" {
    address        = "http://localhost:8080/v1/tfstate/test1"
    lock_address   = "http://localhost:8080/v1/tfstate/test1"
    unlock_address = "http://localhost:8080/v1/tfstate/test1"
    username       = "arpanrec"
  }
}

resource "null_resource" "test" {
  provisioner "local-exec" {
    command = "echo hello"
  }
}
