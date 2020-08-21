variable "aws_region" {
  default = "us-east-1"
}

variable "instance_type" {
  default = "t3.large"
}

variable "cluster_name" {
  type = string
}

variable "control_plane_size" {
  type = string
  default = "1"
  description = "# of control plane nodes {1,3,5}"
}

variable "worker_pool_size" {
  type = string
  default = "1"
  description = "# of worker nodes"
}

variable "vpc_id" {
  type = string
}

variable "domain" {
  type = string
  description = "hosted zone name used for domain name (no trailing dot)"
}

variable "kubernetes_version" {
  type = string
  default = "1.16.8"
}

variable "cluster_backup" {
  type = bool
  default = true
}
