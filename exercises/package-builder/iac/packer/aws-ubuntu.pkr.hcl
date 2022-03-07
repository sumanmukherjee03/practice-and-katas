packer {
  required_version = ">= 1.8.0, < 2.0.0"
  required_plugins {
    amazon = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "package_release_path" {
  type    = string
  default = "/brave-browser/src/out/Release"
}

source "amazon-ebs" "ubuntu" {
  ami_name                    = "sm-package-builder-{{timestamp}}"
  instance_type               = "t2.micro"
  region                      = "us-west-2"
  associate_public_ip_address = true
  ssh_username                = "ubuntu"
  shutdown_behavior           = "terminate"
  ami_regions                 = ["us-west-2"]

  run_tags = {
    Name        = "sm"
    Description = "ami builder for package builder"
  }

  tags = {
    Name        = "sm"
    Description = "ami for package builder"
  }

  source_ami_filter {
    filters = {
      name                             = "ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"
      root-device-type                 = "ebs"
      virtualization-type              = "hvm"
      architecture                     = "x86_64"
      "block-device-mapping.volume-type" = "gp2"
    }

    most_recent = true
    owners      = ["099720109477"]
  }
}

build {
  name = "sm"

  sources = [
    "source.amazon-ebs.ubuntu",
  ]

  provisioner "shell" {
    environment_vars = [
      "PACKAGE_RELEASE_PATH=${var.package_release_path}",
    ]

    script          = "provision.sh"
    execute_command = "echo 'packer' | sudo -S /bin/bash -c '{{ .Vars }} {{ .Path }}'"
  }
}
