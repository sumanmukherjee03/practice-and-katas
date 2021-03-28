resource "aws_route53_zone" "sample_app" {
  name          = "sampleappmrbc.com"
  force_destroy = true
}
