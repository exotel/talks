// INSTANCE START OMIT
resource "aws_instance" "api" {
  ami = "${lookup(var.aws_amis, var.aws_region)}"
  instance_type = "t2.micro"
  key_name = "kube_aws_rsa"
  security_groups = ["${aws_security_group.default_sg.id}"]
  subnet_id = "${aws_subnet.private.id}"
  iam_instance_profile = "webserver-role"
  count = 1
  tags {
    Name = "API"
    env = "${lookup(var.aws_env_tag, var.aws_region)}"
    service = "api"
  }
}
// INSTANCE STOP OMIT

// ELB START OMIT
resource "aws_elb" "api_private" {
  name = "private-api"
  subnets         = ["${aws_subnet.private.id}"]
  security_groups = ["${aws_security_group.default_sg.id}"]
  instances       = ["${aws_instance.api.*.id}"]
  cross_zone_load_balancing = true
  internal = "true"
  listener {
    instance_port     = 1212
    instance_protocol = "tcp"
    lb_port           = 80
    lb_protocol       = "tcp"
  }
  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 5
    target = "HTTP:4240/metrics"
    interval = 30
  }
  tags {
    Name = "prod-api-internal-elb"
    env = "${lookup(var.aws_env_tag, var.aws_region)}"
    service = "api-internal"
  }
}
// ELB STOP OMIT

// DNS START OMIT
#Create the public DNS record for API loadbalancer
resource "aws_route53_record" "api" {
  zone_id = "${var.aws_public_hosted_zone}"
  name = "api"
  type = "A"

  alias {
    name = "${aws_elb.api.dns_name}"
    zone_id = "${aws_elb.api.zone_id}"
    evaluate_target_health = true
  }
}
// DNS STOP OMIT

// SRV START OMIT
#Create a SRV DNS record so that other services can discover api services
resource "aws_route53_record" "api-srv-private" {
  zone_id = "${aws_route53_zone.private.id}"
  name = "_api._http.prod.api.internal.suffix.io."
  type = "SRV"
  ttl = "300"
  records = [
  "${formatlist("1 10 4240 %s", aws_instance.api.*.private_ip)}",
  "${formatlist("1 10 9100 %s", aws_instance.api.*.private_ip)}"
  ]
}
// SRV STOP OMIT
