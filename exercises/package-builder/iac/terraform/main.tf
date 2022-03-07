data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

resource "aws_security_group" "this" {
  name        = "sm"
  description = "Allow HTTP and SSH traffic"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "sm"
  }
}

data "aws_ami" "this" {
  most_recent = true

  filter {
    name   = "name"
    values = ["sm-package-builder-*"]
  }

  owners = ["self"]
}

resource "aws_instance" "this" {
  key_name                    = "sm"
  ami                         = data.aws_ami.this.id
  instance_type               = "c5.9xlarge"
  associate_public_ip_address = true
  subnet_id                   = data.aws_subnets.default.ids[3]
  vpc_security_group_ids      = [aws_security_group.this.id]
  ebs_optimized               = true
  disable_api_termination     = false

  root_block_device {
    volume_size = 120
    volume_type = "gp2"
  }

  tags = {
    Name = "sm"
  }

  volume_tags = {
    Name = "sm"
  }
}
