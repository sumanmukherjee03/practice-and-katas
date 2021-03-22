provider "aws" {
  region = "us-west-2"
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "sample-setup"

  cidr = "10.0.0.0/16"

  azs              = ["us-west-2a", "us-west-2b", "us-west-2c"]
  private_subnets  = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets   = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  database_subnets = ["10.0.121.0/24", "10.0.122.0/24", "10.0.123.0/24"]

  create_database_subnet_group = true

  enable_ipv6 = true

  enable_nat_gateway = false
  single_nat_gateway = true

  enable_s3_endpoint       = true
  enable_dynamodb_endpoint = true

  tags = {
    Owner       = "terraform"
    Environment = "dev"
  }

  vpc_tags = {
    Name = "sample-setup"
  }
}
