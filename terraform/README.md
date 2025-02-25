# Local config files
- terraform.tfvars  - contains local values for all items marked "sensitive" in variables.tf
- state.config - for use by the makefile's calling of terraform init, contains the values for the following variables related to state file management:
    - bucket
    - key
    - access_key
    - secret_key
    - region
    
# credentials

https://cloud.gov/docs/services/cloud-gov-service-account/

