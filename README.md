# Utils Server

## Terraform HTTP Backend

This can be used as a backend for terraform.

URL format: `<protocol>://<host>:<port>/tfstate/<workspace>`

```hcl
terraform {
  backend "http" {
    address        = "http://localhost:8080/tfstate/test"
    lock_address   = "http://localhost:8080/tfstate/test"
    unlock_address = "http://localhost:8080/tfstate/test"
    username       = "test"
    password       = "test"
  }
}
```
