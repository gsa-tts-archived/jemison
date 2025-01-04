# deploying jemison

The Terraform layout is inspired by the conversation here:

https://stackoverflow.com/a/74655690

Because we're a fundamentally small application (or, a small number of services), we're going to take a simpler "dev/staging/prod" approach to organization.

## layout

The `terraform` folder has folders for each environment. The `sandbox` environment (if we create it) is code that is intended to run locally on a per-developer basis for deploying the application to a cloud.gov sandbox (1GB RAM).


# deploying

This is a very early README. At this point, it assume only one person is running deploys. We have work to do still.

## credentials

First, you need to set up deployment creds.

https://cloud.gov/docs/services/cloud-gov-service-account/

These will go in `terraform.tfvars`

```
cf_username   = ""
cf_password   = ""
api_key      = ""
```

(Where does the API key come from? ...)

## make 

Next, 

```
make apply_all
```

which will run a deploy from start to finish. 

It always deletes everything before proceeding. Once we have an S3 backend, this will change.