terraform {
  required_providers {
    brickbybrick = {
      source = "registry.terraform.io/brickbybrickfitness/brickbybrick"
    }
  }
}

variable "brickbybrick_api_key" {
  type = string
}


provider "brickbybrick" {
  api_key = var.brickbybrick_api_key
}

resource "brickbybrick_exercise" "my_exercise_created_with_tf" {
  name = "Dumbbell floor press"
}



