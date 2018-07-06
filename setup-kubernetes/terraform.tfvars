
# Name for role, policy and cloud formation stack (without DBG-DEV- prefix)
cluster_name = "kubernetes"

# AWS region where should the Minikube be deployed
aws_region = "eu-central-1"

#
aws_zones = ["eu-central-1a", "eu-central-1b", "eu-central-1c"]

# Instance type
aws_instance_type = "t2.medium"

# SSH key name for the machine
aws_ssh_key_name = "kubernetes-user"

#
vpc_name = "kubernetes"

#
vpc_cidr = "10.0.0.0/16"

#
private_subnets = "false"

# DNS zone where the domain is placed
hosted_zone = "hvt.zone"

# Tags
tags = {
  Application = "kubernetes"
}
