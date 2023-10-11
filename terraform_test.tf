terraform {
  backend "http" {
    address        = "http://localhost:8080/tfstate/test1"
    lock_address   = "http://localhost:8080/tfstate/test1"
    unlock_address = "http://localhost:8080/tfstate/test1"
    username       = "test"
    password       = "test"
  }
}

resource "null_resource" "test" {
  provisioner "local-exec" {
    command = "echo hello"
  }
}
