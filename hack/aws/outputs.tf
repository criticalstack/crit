output "control_plane_endpoint" {
  value = local.control_plane_endpoint
}

output "control_plane_ips" {
  value = aws_instance.control_plane.*.private_ip
}

output "worker_ips" {
  value = aws_instance.worker.*.private_ip
}

