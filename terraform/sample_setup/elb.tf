resource "aws_security_group" "sample_app_elb" {
  name        = "sample-app-elb-security-group"
  description = "Allow HTTPS traffic from outside"

  vpc_id = module.vpc.vpc_id

  # HTTPS access from anywhere
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTP access from anywhere
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # outbound internet access
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

resource "aws_elb" "sample_app" {
  name            = "sample-app-elb"
  subnets         = [module.vpc.public_subnets[0]]
  security_groups = ["${aws_security_group.sample_app_elb.id}"]

  listener {
    instance_port     = 3000
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
    # lb_port            = 443
    # lb_protocol        = "https"
    # ssl_certificate_id = module.acm.this_acm_certificate_arn
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 3
    target              = "HTTP:80/"
    interval            = 30
  }

  instances                   = ["${aws_instance.sample_app.id}"]
  cross_zone_load_balancing   = true
  idle_timeout                = 300
  connection_draining         = true
  connection_draining_timeout = 300

  tags = {
    Name  = "sample-app-security-group"
    Owner = "terraform"
  }
}

resource "aws_route53_record" "sample_app_elb" {
  zone_id = aws_route53_zone.sample_app.zone_id
  name    = "app1.sampleappmrbc.com"
  type    = "A"

  alias {
    name                   = aws_elb.sample_app.dns_name
    zone_id                = aws_elb.sample_app.zone_id
    evaluate_target_health = false
  }
}
