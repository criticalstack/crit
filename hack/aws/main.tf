terraform {
  required_version = ">= 0.12.0"
}

provider "aws" {
  region = "us-east-1"
  version = "~> 2.0"
}

data "aws_vpc" "vpc" {
  id = "${var.vpc_id}"
}

data "aws_subnet_ids" "private" {
  vpc_id = "${var.vpc_id}"

  filter {
    name   = "tag:Name"
    values = ["*private*"]
  }

  #filter {
    #name  = "availability-zone"
    #values = ["us-east-1a", "us-east-1b", "us-east-1c", "us-east-1d"] # the other AZ 1e doesn't support t3
  #}
}

data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "tls_private_key" "ssh" {
  count = fileexists("~/.ssh/id_rsa.pub") ? 0 : 1
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "aws_key_pair" "dev" {
  key_name   = "${var.cluster_name}-ssh"
  public_key = fileexists("~/.ssh/id_rsa.pub") ? file("~/.ssh/id_rsa.pub") : tls_private_key.ssh[0].public_key_openssh
}

#####################################################################
# S3
#####################################################################

resource "aws_s3_bucket" "backup" {
  bucket        = "${var.cluster_name}-s3"
  acl           = "private"
  force_destroy = true

  versioning {
    enabled = false
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
}

resource "aws_s3_bucket_policy" "backup" {
  bucket = aws_s3_bucket.backup.id

  policy = jsonencode({
    Statement = [
      {
        Effect = "Deny"
        Principal = "*"
        Action = "s3:*"
        Resource = [
          "${aws_s3_bucket.backup.arn}",
          "${aws_s3_bucket.backup.arn}/*"
        ]
        Condition = {
          StringNotLike: {
            "aws:arn": [
              "arn:aws:iam::${local.account_id}:user/vault-github-*",
              "arn:aws:iam::${local.account_id}:role/${module.iam.this_iam_role_name}",
              "arn:aws:sts::${local.account_id}:assumed-role/${module.iam.this_iam_role_name}/*",
            ]
          }
        }
      },
    ]
    Version = "2012-10-17"
  })
}

resource "aws_s3_bucket_object" "certs" {
  for_each = fileset(".pki", "**")

  bucket = aws_s3_bucket.backup.id
  key = "pki/${each.value}"
  source = ".pki/${each.value}"
}

resource "aws_s3_bucket_object" "crit_bin_object" {
  bucket = aws_s3_bucket.backup.id
  key    = "bin/crit"
  source = "../../bin/crit"
}

#####################################################################
# IAM
#####################################################################

module "iam" {
  source = "terraform-aws-modules/iam/aws//modules/iam-assumable-role"

  role_name = "${var.cluster_name}-iam-role-control-plane"

  create_role             = true
  create_instance_profile = true
  role_requires_mfa       = false

  trusted_role_services = [
    "ec2.amazonaws.com",
  ]

  custom_role_policy_arns = [
    for _, v in aws_iam_policy.control_plane_policies: v.arn
  ]
}

resource "aws_iam_policy" "control_plane_policies" {
  for_each = fileset("./userdata/control_plane/policies", "*.json")

  name = "${var.cluster_name}-${trimsuffix(each.value, ".json")}"
  policy = file("./userdata/control_plane/policies/${each.value}")
}

resource "aws_iam_policy" "workers_policies" {
  for_each = fileset("./userdata/workers/policies", "*.json")

  name = "${var.cluster_name}-${trimsuffix(each.value, ".json")}"
  policy = file("./userdata/workers/policies/${each.value}")
}

#####################################################################
# Security group(s)
#####################################################################

module "sg" {
  source = "terraform-aws-modules/security-group/aws"

  name        = "${var.cluster_name}-sg"
  description = "${var.cluster_name} security group"
  vpc_id      = var.vpc_id

  ingress_with_self = [{ rule = "all-all" }]
  egress_with_self  = [{ rule = "all-all" }]

  ingress_with_cidr_blocks = [
    {
      rule        = "all-all"
      cidr_blocks = data.aws_vpc.vpc.cidr_block
    },
  ]
  egress_with_cidr_blocks = [
    {
      rule        = "all-all"
      cidr_blocks = "0.0.0.0/0"
    },
  ]
}

#####################################################################
# Control plane
#####################################################################

module "control_plane_cloudinit" {
  source = "./modules/cloudinit"

  vars = {
    account_id             = local.account_id
    cluster_name           = var.cluster_name
    control_plane_endpoint = local.control_plane_endpoint
    control_plane_size     = var.control_plane_size
    hosted_zone_name       = var.domain
    kubernetes_version     = var.kubernetes_version
    s3_bucket_name         = aws_s3_bucket.backup.id
    cluster_backup         = var.cluster_backup
  }

  files = {
    src_dir  = "userdata/control_plane/files"
    root_dir = "/"
  }

  scripts = {
    src_dir  = "userdata/control_plane/scripts"
    dest_dir = "/opt/criticalstack"
  }
}

resource "aws_instance" "control_plane" {
  count                  = var.control_plane_size
  ami                    = data.aws_ami.ubuntu.id
  iam_instance_profile   = module.iam.this_iam_instance_profile_name
  instance_type          = var.instance_type
  key_name               = aws_key_pair.dev.key_name
  user_data              = module.control_plane_cloudinit.rendered
  subnet_id              = element(tolist(data.aws_subnet_ids.private.ids), 0)
  vpc_security_group_ids = [module.sg.this_security_group_id]

  root_block_device {
    volume_type = "gp2"
    volume_size = "50"
  }

  tags = {
    "Name"                                      = "${var.cluster_name}-controlplane"
    "kubernetes.io/cluster/${var.cluster_name}" = "owned"
    "kubernetes.io/role/controlplane"           = "1"
  }
}



data "aws_route53_zone" "selected" {
  name         = "${var.domain}."
  private_zone = true
}

resource "aws_route53_record" "control_plane_endpoint" {
  zone_id = data.aws_route53_zone.selected.zone_id
  name    = "${var.cluster_name}.${data.aws_route53_zone.selected.name}"
  type    = "A"
  ttl     = 60
  records = aws_instance.control_plane.*.private_ip
}

#####################################################################
# Workers
#####################################################################

module "workers_cloudinit" {
  source = "./modules/cloudinit"

  vars = {
    cluster_name           = var.cluster_name
    control_plane_endpoint = local.control_plane_endpoint
    kubernetes_version     = var.kubernetes_version
    s3_bucket_name         = aws_s3_bucket.backup.id
  }

  files = {
    src_dir  = "userdata/workers/files"
    root_dir = "/"
  }

  scripts = {
    src_dir  = "userdata/workers/scripts"
    dest_dir = "/opt/criticalstack"
  }
}

resource "aws_instance" "worker" {
  count                  = var.worker_pool_size
  ami                    = data.aws_ami.ubuntu.id
  iam_instance_profile   = module.iam.this_iam_instance_profile_name
  instance_type          = var.instance_type
  key_name               = aws_key_pair.dev.key_name
  subnet_id              = element(tolist(data.aws_subnet_ids.private.ids), 0)
  user_data              = module.workers_cloudinit.rendered
  vpc_security_group_ids = [module.sg.this_security_group_id]

  root_block_device {
    volume_type = "gp2"
    volume_size = "50"
  }

  tags = {
    "Name"                                      = "${var.cluster_name}-worker"
    "kubernetes.io/cluster/${var.cluster_name}" = "owned"
    "kubernetes.io/role/worker"                 = "1"
  }
}
