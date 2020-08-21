# tfconfig

`tfconfig` is a tool for creating Terraform variable definitions files. It reads Terraform configuration files (`.tf`, `.hcl`, `.tfvars`) from a directory, parses the files, and generates an interactive prompt for each variable.

```
$ tfconfig create path/to/dir
environment = "dev"
region      = "east"
tag         = "1.0.0"

? use existing config No
? environment dev
? region west
? Tag of the container image 1.0.1

environment = "dev"
region      = "west"
tag         = "1.0.1"

? write this file to path/to/dir/terraform.tfvars Yes
```
