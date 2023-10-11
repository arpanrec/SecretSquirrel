terraform {
  backend "http" {
    address        = "http://localhost:8080/tfstate/test"
    username       = "test"
    password       = "test"
  }
}

resource "null_resource" "test" {
  provisioner "local-exec" {
    command = "echo hello"
  }
}
