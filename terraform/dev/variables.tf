# This is a bodge until we add a different access control mechanism
variable api_key {
  type = string
  sensitive = true
}

variable cf_api_url {
  type = string
  default = "https://api.fr.cloud.gov"
}

variable cf_org {
  type = string
  description = "cloud.gov organization name (e.g. `gsa-tts-usagov`)"
  default = "gsa-tts-usagov"
}

variable cf_org_guid {
  type = string
  sensitive = true
}

variable cf_space_guid {
  type = string
  sensitive = true
}

# Provided in terraform.tfvars
# This is a GH secret
variable cf_password {
  type = string
  sensitive = true
}

# Provided in terraform.tfvars
# This is a GH secret
variable cf_username {
  type = string
  sensitive = true
}

variable gin_debug_level {
  type = string
  default = "release"
}

variable "gitref" {
  type        = string
  description = "gitref for the specific version of app that you want to use"
  default     = "refs/heads/main"
  # You can also specify a specific commit, eg "7487f882903b9e834a5133a883a88b16fb8b16c9"
}
