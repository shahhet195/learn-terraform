output "public_ip" {
  description = "Public IP address of the EC2 instance"
  value       = aws_instance.ubuntu.public_ip
}

output "public_dns" {
  description = "Public DNS name of the EC2 instance"
  value       = aws_instance.ubuntu.public_dns
}
