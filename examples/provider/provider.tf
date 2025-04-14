variable "api_key" {
  type = string
}

provider "brickbybrick" {
  api_key = var.api_key
}
