terraform {
  required_version = "~> 1.0"
  required_providers {
    cloudfoundry = {
      source  = "cloudfoundry-community/cloudfoundry"
      version = "~>0.51.3"
    }
  }
  backend "s3" {
    bucket      = ""
    key         = ""
    access_key  = ""
    secret_key  = ""
    region      = ""
    }
}

provider "cloudfoundry" {
  api_url      = "https://api.fr.cloud.gov"
  user         = var.cf_username
  password     = var.cf_password
  app_logs_max = 30
}