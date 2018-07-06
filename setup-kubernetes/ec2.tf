#
# variables from other files
#

locals {
  aws_subnet_id = "${aws_subnet.public_subnet.0.id}"
  kubeadm_token = "${data.template_file.kubeadm_token.rendered}"
}

#
# Security Group
#

data "aws_subnet" "minikube_subnet" {
  id = "${local.aws_subnet_id}"
}

resource "aws_security_group" "minikube" {
  vpc_id = "${data.aws_subnet.minikube_subnet.vpc_id}"
  name   = "${var.cluster_name}"

  tags = "${merge(map("Name", var.cluster_name, format("kubernetes.io/cluster/%v", var.cluster_name), "owned"), var.tags)}"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 6443
    to_port     = 6443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

#
# IAM role
#

data "template_file" "iam_policy_json" {
  template = "${file("${path.module}/templates/iam_policy.json.tpl")}"

  vars {}
}

data "template_file" "iam_role_json" {
  template = "${file("${path.module}/templates/iam_role.json.tpl")}"

  vars {}
}

resource "aws_iam_policy" "minikube_policy" {
  name        = "${var.cluster_name}"
  path        = "/"
  description = "Policy for role ${var.cluster_name}"
  policy      = "${data.template_file.iam_policy_json.rendered}"
}

resource "aws_iam_role" "minikube_role" {
  name = "${var.cluster_name}"
  assume_role_policy = "${data.template_file.iam_role_json.rendered}"
}

resource "aws_iam_policy_attachment" "minikube-attach" {
  name       = "minikube-attachment"
  roles      = ["${aws_iam_role.minikube_role.name}"]
  policy_arn = "${aws_iam_policy.minikube_policy.arn}"
}

resource "aws_iam_instance_profile" "minikube_profile" {
  name = "${var.cluster_name}"
  role = "${aws_iam_role.minikube_role.name}"
}

#
# Bootstraping scripts
#

data "template_file" "init_minikube" {
  template = "${file("${path.module}/templates/init-aws-minikube.sh")}"

  vars {
    kubeadm_token = "${local.kubeadm_token}"
    dns_name      = "${var.cluster_name}.${var.hosted_zone}"
    ip_address    = "${aws_eip.minikube.public_ip}"
    cluster_name  = "${var.cluster_name}"
  }
}

data "template_file" "cloud-init-config" {
  template = "${file("${path.module}/templates/cloud-init-config.yaml")}"

  vars {
    calico_yaml = "${base64gzip("${file("${path.module}/templates/calico.yaml")}")}"
  }
}

data "template_cloudinit_config" "minikube_cloud_init" {
  gzip          = true
  base64_encode = true

  part {
    filename     = "cloud-init-config.yaml"
    content_type = "text/cloud-config"
    content      = "${data.template_file.cloud-init-config.rendered}"
  }

  part {
    filename     = "init-aws-minikube.sh"
    content_type = "text/x-shellscript"
    content      = "${data.template_file.init_minikube.rendered}"
  }
}

#
# EC2 instance
#

# the latest CentOS 7 image will be used, see REAMDE how to activate it on AWS
data "aws_ami" "centos7" {
  most_recent = true
  owners      = ["aws-marketplace"]

  filter {
    name   = "product-code"
    values = ["aw0evgkw8e5c1q413zgy5pjce"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_eip" "minikube" {
  vpc = true
}

resource "aws_instance" "minikube" {
  instance_type = "${var.aws_instance_type}"

  ami = "${data.aws_ami.centos7.id}"

  key_name = "${var.aws_ssh_key_name}"

  subnet_id = "${data.aws_subnet.minikube_subnet.id}"

  associate_public_ip_address = false

  vpc_security_group_ids = [
    "${aws_security_group.minikube.id}",
  ]

  iam_instance_profile = "${aws_iam_instance_profile.minikube_profile.name}"

  user_data = "${data.template_cloudinit_config.minikube_cloud_init.rendered}"

  tags = "${merge(map("Name", var.cluster_name, format("kubernetes.io/cluster/%v", var.cluster_name), "owned"), var.tags)}"

  root_block_device {
    volume_type           = "gp2"
    volume_size           = "50"
    delete_on_termination = true
  }

  lifecycle {
    ignore_changes = [
      "ami",
      "user_data",
      "associate_public_ip_address",
    ]
  }
}


resource "aws_eip_association" "minikube_assoc" {
  instance_id   = "${aws_instance.minikube.id}"
  allocation_id = "${aws_eip.minikube.id}"
}

#
# DNS records
#

data "aws_route53_zone" "dns_zone" {
  name         = "${var.hosted_zone}."
  private_zone = "${var.hosted_zone_private}"
}

resource "aws_route53_record" "minikube" {
  zone_id = "${data.aws_route53_zone.dns_zone.zone_id}"
  name    = "${var.cluster_name}.${var.hosted_zone}"
  type    = "A"
  records = ["${aws_eip.minikube.public_ip}"]
  ttl     = 300
}

resource "aws_route53_record" "basic_hvt_zone" {
  zone_id = "${data.aws_route53_zone.dns_zone.zone_id}"
  name    = "${data.aws_route53_zone.dns_zone.name}"
  type    = "A"
  ttl     = "300"
  records = ["${aws_eip.minikube.public_ip}"]
}

resource "aws_route53_record" "prefix_hvt_zone" {
  zone_id = "${data.aws_route53_zone.dns_zone.zone_id}"
  name    = "*.${data.aws_route53_zone.dns_zone.name}"
  type    = "A"
  ttl     = "300"
  records = ["${aws_eip.minikube.public_ip}"]
}

