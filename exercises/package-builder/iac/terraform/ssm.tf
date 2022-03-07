resource "aws_ssm_parameter" "ssh_private_key" {
  name        = "sm"
  type        = "SecureString"
  value       = "<set_via_aws_console_manually>"
  description = "Private key to ssh into package-builder instance"

  lifecycle {
    ignore_changes = [value]
  }
}
