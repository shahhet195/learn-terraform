provider "aws" {
  region = var.region
}

module "ec2" {
  source = "./modules/ec2"

  instance_type = var.instance_type
  instance_name = var.instance_name
}