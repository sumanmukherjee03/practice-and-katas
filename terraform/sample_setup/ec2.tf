resource "aws_security_group" "sample_app" {
  name        = "sample-app-security-group"
  description = "Allow HTTP and SSH traffic"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["96.48.69.126/32"]
  }

  ingress {
    description = "HTTP"
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = ["96.48.69.126/32"]
  }

  ingress {
    description = "HTTP"
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name  = "sample-app-security-group"
    Owner = "terraform"
  }
}

data "aws_ami" "sample_app" {
  most_recent = true

  filter {
    name   = "name"
    values = ["sample-app-image"]
  }

  owners = ["self"]
}

resource "aws_instance" "sample_app" {
  key_name                    = "sample-app-key"
  ami                         = data.aws_ami.sample_app.id
  instance_type               = "t2.micro"
  associate_public_ip_address = true
  subnet_id                   = module.vpc.public_subnets[0]
  vpc_security_group_ids      = [aws_security_group.sample_app.id]
  ebs_optimized               = false
  disable_api_termination     = false

  tags = {
    Name  = "sample-app"
    Owner = "terraform"
  }
}
