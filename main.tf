# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "aws" {
  region = var.region
}

resource "aws_instance" "ubuntu" {
  ami           = "ami-04b70fa74e45c3917"
  instance_type = var.instance_type

  tags = {
    Name = var.instance_name
  }
}
