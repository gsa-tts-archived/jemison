# deploying jemison

The Terraform layout is inspired by the conversation here:

https://stackoverflow.com/a/74655690

Because we're a fundamentally small application (or, a small number of services), we're going to take a simpler "dev/staging/prod" approach to organization.

# establishing the org/space

The TF included in this repository does *not* attempt to stand up everything from scratch. That is, if we suffer a catastrophic failure, and lose all of Cloud.gov, and have to rebuild from scratch... there will be manual steps.

## create the org and space

We are currently in the org `gsa-tts-usagov`, and that may change. 

We intend to have three spaces: `dev`, `staging`, and `production`. 

* Every merge to `main` will push to `dev` and run tests.
* Every morning we will push to `staging` and run E2E tests.
* Pushes to `production` will be via release cuts, and be manual.

## create the state bucket

In each space, we need a bucket called `tf_state`. This currently is *not* be a TF managed resource. 

```
cf create-service s3 basic tf_state
```

Then, get the credentials. (Note the `source` in there.)

https://cloud.gov/docs/services/s3/#:~:text=Using%20the%20S3%20credentials

Set `SERVICE_KEY_NAME` to something (e.g. `<username>-s3-service-key`) in your env, and then run the script in `bin`.

```
export SERVICE_KEY_NAME="service-key-name-of-some-sort"
source ./bin/get-s3-creds-from-cgov.bash
```

One you have run the script, those variables need to be loaded into GH for use in actions. They will be `environment` secrets---meaning, we will need a `tf_state` bucket in each env/space, and we will need a set of credentials for each space in GH.

* TF_VAR_BUCKET_NAME
* TF_VAR_AWS_DEFAULT_REGION
* TF_VAR_SECRET_ACCESS_KEY
* TF_VAR_ACCESS_KEY_ID

By prefixing these with `TF_VAR_`, they are directly accessible in our Terraform scripts.

### for reference...

```
git add ../.github/workflows/deploy-to-dev.yaml ; git commit -m "Iterating GH Action $(date)" ; git push

gh workflow run --ref jadudm/tf-0103 --field environment=dev deploy-to-dev.yaml
```


### obtaining local credentials (remove)

This is only during local dev... this goes away once we're in Github Actions.

https://cloud.gov/docs/services/cloud-gov-service-account/

These will go in `terraform.tfvars`

```
cf_username   = ""
cf_password   = ""
api_key      = ""
```


## launching the stack

```
make dev
```

at the top of the tree will deploy the `dev` stack. More work needs to be done in order to store the TF state in S3, so that we can run this from Github Actions. For now, this is not complete; if different devs deploy, they will have to completely destroy (tear down) the state of the other devs. This will become... annoying... once we start storing data in buckets. (Buckets must be empty in order to be torn down.) 

So, the deploy to Cloud.gov is still a work-in-progress. But, it is possible, while testing/developing, to do a deploy from a local machine. Once we have GH Actions in place, we will *never* do a deploy from a local machine. We will always do our deploys from an action.

## layout

At the top of the `terraform` directory are two files that matter:

* Makefile
* developers.tf

`developers.tf` will become part of our onboarding. This file is where devs add themselves as an initial commit so that they gain access to the Cloud.gov environment. We will control access to Cgov through this file. (This wiring is not in place yet, but the file is there. The access controls have to be implemented as scripts executed in a Github Action that call the CF API on Cloud.gov.)

Cgov deployments are organized into `organizations` and `spaces`. An organization might be `gsa-tts-search`, and a space might be `dev`, `staging`, or `production`. 

There are two directories (currently) that contain the Terraform deploy scripts:

* dev
* shared

`dev` contains the variables and drivers for deploying to our (eventual) `dev` space. Every service that we deploy will get a section in this file:

```
module "fetch" {
  source = "../shared/services/fetch"
  # disk_quota = 256
  # memory = 128
  # instances = 1
  space_name = data.cloudfoundry_space.app_space.name
  app_space_id = data.cloudfoundry_space.app_space.id
  domain_id = data.cloudfoundry_domain.public.id
  databases = module.databases.ids
  buckets = module.buckets.ids
}
```

I have not yet determined if this can be made reusable between spaces (meaning, avoiding the boilerplate-ness of this). Each service has to be wired up to the correct databases and S3 buckets _in its space_ in order to execute. Further, we might want to allocate different amounts of RAM, disk, and instances to services in the different spaces. That is, we might one 1 instance of `fetch` in the `dev` environment, but 3 instances of `fetch` in `production`. Because we only have one pool of RAM for all of the spaces combined, we will probably run light in lower environments, and run a fuller stack in `production`. 

The service itself is defined in `shared/services/<service-name>`. We apparently have to include the provider (?), define the variables for the module, the outputs, and the module itself. Put another way:

* `providers.tf` is boilerplate. It will need to change when we switch to the official `cloudfoundry/cloudfoundry` provider.
* `variables.tf` defines the variables that the service needs to have defined in order to execute. For example, when instantiating the module, we need to provide the amount of RAM, disk, and the number of instances the service will be created with.
* `service.tf` defines the service itself.

We can see the `fetch` service:

```
resource "cloudfoundry_app" "fetch" {
  name                 = "fetch"
  space                = var.app_space_id # data.cloudfoundry_space.app_space.id
  buildpacks            = ["https://github.com/cloudfoundry/apt-buildpack", "https://github.com/cloudfoundry/binary-buildpack.git"]
  path                 = "${path.module}/../app.tar.gz"
  source_code_hash     = filesha256("${path.module}/../app.tar.gz")
  disk_quota           = var.disk_quota
  memory               = var.memory
  instances            = var.instances
  strategy             = "rolling"
  timeout              = 200
  health_check_type    = "port"
  health_check_timeout = 180
  health_check_http_endpoint = "/api/heartbeat"

  service_binding {
    service_instance = var.databases.queues
  }

  service_binding {
    service_instance = var.databases.work
  }

  service_binding {
    service_instance = var.buckets.fetch
  }
}
```

All of the services get the entire codebase; this is because we then launch, on a per-instance basis, different code from `cmd`.

Variables include the ID of the space we are deploying to (e.g. we do not deploy to `dev`, but to a UUID4 value representing `dev`), the disk, memory, and instances, and more importantly, bindings to the databases and S3 buckets. 

### buckets and databases

In `shared/cloudgov` are module definitions for our databases and S3 buckets.

In `dev/main.tf`, we instantiate these as follows:

```
module "databases" {
  source              = "../shared/cloudgov/databases"
  cf_org              = local.cf_org
  cf_space            = local.cf_space
  queue_db_plan_name  = "micro-psql"
  search_db_plan_name = "micro-psql"
  work_db_plan_name   = "micro-psql"
}
```

For `dev`, we might only use `micro` instances. For production, however, we might instantiate `xl` instances. This lets us configure the databases on a per-space basis. (S3 buckets are all the same, so there is no configuration.)

This module has outputs. Once instantiated, we can refer to `module.databases` as a `map(string)` and reference the `id` of each of the databases (or buckets). In this way, we can pass the entire map of IDs to the services, and they can then bind to the correct databases/S3 buckets. Most (all?) services will want to bind to the queues databases; only some need to bind to `work`, and some need to bind to `serve`. 

