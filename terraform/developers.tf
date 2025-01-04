
# Mirrors parts of 
# https://github.com/GSA-TTS/FAC/blob/main/terraform/meta/config.tf
# This is copied into the space directories (e.g. `dev`, `staging`)
# as part of the Terraform build.

locals {
  spacedevelopers = [
    "matthew.jadud@gsa.gov"
  ]

  spacemanagers = [
    "matthew.jadud@gsa.gov"
  ]
}
