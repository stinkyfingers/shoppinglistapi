terraform {
  backend "s3" {
    bucket = "remotebackend"
    key    = "shoppinglistapi/terraform.tfstate"
    region = "us-west-1"
    profile = "jds"
  }
}
