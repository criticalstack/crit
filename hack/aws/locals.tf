data "aws_caller_identity" "current" {}

locals {
  account_id             = data.aws_caller_identity.current.account_id
  current_user           = data.aws_caller_identity.current.arn
  control_plane_endpoint = "${var.cluster_name}.${var.domain}"
}
