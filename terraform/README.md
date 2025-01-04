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