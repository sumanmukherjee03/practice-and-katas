module "acm" {
  source                    = "terraform-aws-modules/acm/aws"
  zone_id                   = aws_route53_zone.sample_app.zone_id
  domain_name               = "sampleappmrbc.com"
  subject_alternative_names = ["*.sampleappmrbc.com"]
  wait_for_validation       = false
}
