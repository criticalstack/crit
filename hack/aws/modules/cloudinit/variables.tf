variable "debug" {
  type    = bool
  default = false
}

variable "files" {
  type    = object({ src_dir=string, root_dir=string })
  default = null
}

variable "scripts" {
  type    = object({ src_dir=string, dest_dir=string })
  default = null
}

variable "main" {
  type    = string
  default = null
}

variable "vars" {
  type    = map(any)
  default = null
}
