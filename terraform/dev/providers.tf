terraform {
  required_version = "~> 1.0"
  required_providers {
    cloudfoundry = {
      source  = "cloudfoundry-community/cloudfoundry"
      version = "~>0.51.3"
    }
  }

  backend "s3" {
    key     = "state/terraform.tfstate"
    encrypt = true
  }
}

provider "cloudfoundry" {
  api_url      = var.cf_api_url
  user         = var.cf_username
  password     = var.cf_password
  app_logs_max = 30
}
