# deployment

The point of this documentation is to help lay out some of the pieces (Cloud.gov, Cloud Foundry, Terraform, CircleCI) and the way they relate. It suggests a series of steps we can take to iteratively get to a working set of deployments. The steps might not be right: a good DevOps engineer will see other paths, perhaps. But, this provides an example of the kind of thinking that might help us get things stood up quickly.

## foundations

We'll be deploying to [cloud.gov](https://cloud.gov), which is essentially [Cloud Foundry](https://www.cloudfoundry.org/). 

Under the hood, Cloud.gov is a bunch of AWS EC2 instances that bounce containerized apps around for load balancing. The bouncing and balancing is invisible to us as application developers. Caching, WAF, and other network-related support is all invisible and part of the stack.

Further, CGov provides a set of *brokered services*. These are AWS services (e.g. Postgres, S3, Elastic) that we can use, but do not have to provision and manage directly. For example, we can declare we want a Postgres database, and the platform will provide credentials to our application as part of the execution environment. *Magic*, as they say. 

Our applications run as part of an **org**anization, in which there are multiple **spaces**. These are largely about RBAC, and will have minimal impact on us. The important piece is that we will want to target three spaces (to be named/decided): `dev`, `staging`, and `prod`. 

All of this is managed via the `cf` API. (See the Cloud.gov/Cloud Foundry docs.) We can, for example, see what apps are running:

```
cf apps
```

or what services are brokered:

```
cf services
```

or what routes exist between the services and apps:

```
cf routes
```

we can stop apps

```
cf app fetch stop
```

or restart

```
cf app fetch restart
```

or, perhaps most importantly, tail the logs:

```
cf logs fetch
```

In short, via the `cf` command line tool (which makes API calls to Cloud.gov), we can manage the entire application. 

**Generally, we won't**.

## terraform all the things

We can use Terraform to manage the Cloud.gov environment.

This branch contains a simple `main.tf` that will deploy to a given org/space. In this case, we're (currently) deploying to `gsa-tts-usagov` as an org, and `search-dev` as a space. This might later become `searchgov` as an org and `dev` as a space, or something else. Consider these variables. (You could change them to a sandbox to test your own deployments; they might want to be moved out to a variables file, and be left out of version control.)

The current example is using the community-provided TF provider. Hashicorp has released an officially supported provider; we will switch to that whenever we feel like it. (Could be now, could be later.)

To run a deploy locally, you need to first have the app building/running locally. That is described in `DEVELOPING.md` in the `docs` folder. Then, you will need to obtain a set of credentials for Cloud.gov. Specifically, we need a service account and releated creds:

https://cloud.gov/docs/services/cloud-gov-service-account/

Follow the docs. If you were Darth Vader (preferred username `vaderd`), it would boil down to:

```
cf create-service cloud-gov-service-account space-auditor vaderd-service-account

cf create-service-key vaderd-service-account vaderd-service-key
```

The above would create Vader's service account and key. Then, Darth would need to get the key:

```
cf service-key vaderd-service-account vaderd-service-key
```

At which point, he would get:

```
{
 "password": "oYasdfliaweinasfdliecV",
 "username": "deadbeef-aabb-1234-feha0987654321000"
}
```

Those do not need to be saved. They would be put in a `terraform.tfvars` file as `cf_username` and `cf_password`. Or, they might be stored as secrets in a CI/CD pipeline. However, it is assumed they are disposable. You can always delete the key (invalidating the creds) and create a new one. It will be good practice for us to rotate this account/cred.

Once this is in hand and in the `tfvars`, we can build and deploy.

### deploying

To deploy from local:

```
make terraform
```

This will:

1. Build each of the apps into a binary (named `service.exe`)
2. Copy and zip those files in the `apps` directory inside of `terraform`
3. `plan` and `apply` the TF
4. Delete any existing apps
5. Upload the binary packages to CGov
6. Profit!

This has a number of drawbacks:

1. The application is build as a set of binaries. These are build inside a container. Why? Because while Golang is portable as a binary, it has expectations about the libc it is running on top of. So, we compile the binaries in a `cflinux4` container, which is the same linux/libc as running in Cloud.gov. It's a hack. We don't want to keep this.
2. The TF state is stored locally. This means that if Alice does a deploy, then Bob does a deploy, Bob's deploy will die a messy terradeath. This is because Bob's TF state will absolutely not match Alice's, even though it is "the same."

Hm. That might be it. But the second one is a big problem.

## getting to production

Our goal is a deployment pattern that looks like:

1. Every PR to `main` triggers a deploy to `dev`
2. Every morning, we auto-deploy everything in `main` to `staging`, and (ultimately) run E2E tests.
3. When desired, we cut a release. On release, we deploy to `prod`.

Optionally, we might want to be able to deploy to `dev` from *any branch at any time*. This way, devs can test code in the cloud. TBD.

Ideally:

1. Our first step might be to integrate the `make terraform` into CI/CD. That is, just use the Makefile as-is to get binaries pushed to Cloud.gov. This will have the same problems as above, because we won't be preserving our CI/CD state. This is not a situation we want to live with for long.
2. We want our CI/CD state in S3. We can create a bucket (with the `cf` CLI) for the TF state. From there, we should be able to configure an S3 provider. At this point, we have external state, and should be able to start moving with velocity in the CI/CD env.

Once we hit #2, the thought is that we should be able to do `dev`, `staging`, and `prod` deploys. That is, it should be possible to build in the containerized environment in CircleCI, and ship the binaries. Demonstrating this in at least one space (e.g. `dev`) would give us confidence that the process is working. 

We will want to then explore deploying as a [buildpack](https://docs.cloudfoundry.org/buildpacks/go/index.html) instead of binaries. This lets us pull the code down and push the app as code as opposed to a binary. It will then run "natively" in Cloud.gov. This is preferred for all languages/environments where a buildpack exists. We will probably have to always be a dual buildpack environment: one is the `apt` buildpack (that lets us install some libraries/packages needed by our services) and the `go` buildpack.

We can test buildpack deploys locally. Once we have them working, we can then adapt the Makefile, and our remote builds should "just work." 

## conclusion

The point of this documentation is to help lay out some of the pieces (Cloud.gov, Cloud Foundry, Terraform, CircleCI) and the way they relate. It suggests a series of steps we can take to iteratively get to a working set of deployments. The steps might not be right: a good DevOps engineer will see other paths, perhaps. But, this provides an example of the kind of thinking that might help us get things stood up quickly.