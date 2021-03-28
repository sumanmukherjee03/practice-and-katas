resource "aws_security_group" "sample_app_db" {
  name        = "sample-app-db-security-group"
  description = "Allow MYSQL traffic"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "TCP"
    from_port   = 3306
    to_port     = 3306
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  tags = {
    Name  = "sample-app-db-security-group"
    Owner = "terraform"
  }
}

module "sample_app_db" {
  source = "terraform-aws-modules/rds/aws"

  identifier = "sample-app"

  create_db_option_group    = false
  create_db_parameter_group = false

  # All available versions: http://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_MySQL.html#MySQL.Concepts.VersionMgmt
  engine               = "mysql"
  engine_version       = "5.7.30"
  family               = "mysql5.7" # DB parameter group
  major_engine_version = "5.7"      # DB option group
  instance_class       = "db.t3.small"
  port                 = 3306

  allocated_storage = 5

  name     = "sample_app"
  username = var.db_username
  password = var.db_password

  skip_final_snapshot = true
  deletion_protection = false

  subnet_ids             = module.vpc.database_subnets
  vpc_security_group_ids = [aws_security_group.sample_app_db.id]

  maintenance_window = "Mon:00:00-Mon:03:00"
  backup_window      = "03:00-06:00"

  backup_retention_period = 0

  tags = {
    Owner = "terraform"
  }
}
