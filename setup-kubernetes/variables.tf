
variable "aws_region" {
  description = "Region where Cloud Formation is created"
  default     = "eu-central-1"
}

variable "cluster_name" {
  description = "Name of the AWS Minikube cluster - will be used to name all created resources"
}

variable "tags" {
  description = "Tags used for the AWS resources created by this template"
  type        = "map"
}

variable "aws_instance_type" {
  description = "Type of instance"
  default     = "t2.medium"
}

variable "aws_ssh_key_name" {
  description = "AWS key pair name"
}

variable "hosted_zone" {
  description = "Hosted zone to be used for the alias"
}

variable "hosted_zone_private" {
  description = "Is the hosted zone public or private"
  default     = false
}

#
# VPC
#

variable aws_zones {
  description = "AWS AZs (Availability zones) where subnets should be created"
  type = "list"
}

variable private_subnets {
  description = "Create both private and public subnets"
  type = "string"
  default = "false"
}

variable vpc_name {
  description = "Name of the VPC"
  type = "string"
}

# Network details (Change this only if you know what you are doing or if you think you are lucky)
variable vpc_cidr {
  description = "CIDR of the VPC"
  type = "string"
}

