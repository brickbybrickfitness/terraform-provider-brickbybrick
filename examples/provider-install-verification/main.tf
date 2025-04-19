# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    brickbybrick = {
      source = "registry.terraform.io/brickbybrickfitness/brickbybrick"
    }
  }
}

variable "api_key" {
  type = string
}

provider "brickbybrick" {
  api_key = var.api_key
}

data "brickbybrick_exercises" "my_exercises" {}

output "print_exercises" {
  value = data.brickbybrick_exercises.my_exercises
}
