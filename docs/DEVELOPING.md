# developing

Developing `jemison` locally is possible, and largely requires Docker.

## configuring your system

You will need the following packages installed on your system:

* git
* make
* [docker](https://docs.docker.com/engine/install/)


## getting the code

If you are a member of the team, your onboarding process should have walked you through obtaining access.

If you are a member of the public, you can fork the code, or check out a copy directly.

```
git clone https://github.com/GSA-TTS/jemison
```

## building the system

From the top of the tree, run

`make docker`

This will build two images:

* `jemison/build`: for running builds in a `cfslinux4` containerized space 
* `jemison/dev`: for running the application in a `cfslinux4` environment

If you are on a Mac, or choose to do your development entirely through the containerized environment, then you should now be able to run

`mac macup`

This will build and run all of the services, including pull images for several additional services (e.g. `minio`, `pgweb`) to support interacting with the stack via the browser.

### local development

To develop locally (YMMV, perhaps best for Linux-based systems), you will also need

* jsonnet
* [sqlc](https://docs.sqlc.dev/en/stable/overview/install.html)
* go (1.23+)
* [go tools](https://cs.opensource.google/go/x/tools)

This will let you compile/sometimes test services locally (by setting `ENV=LOCALHOST`), but it is likely only good for dev/unit-test cycles. On an Intel-based system, it is possible to develop and run the stack iteratively (using `make up`). 

On a Mac, it is necessary to use `make macup` to build the tooling in the containerized environment. (This is because the architecture of the containers/deployed environment and the architecture of the Mac are different.)

We have no information on developing on Windows. It is likely that under WSL2, a Linux-based process will work. It should be the case that `make macup` will work regardless.

## testing

... talk about the test server/container ...

... develop E2E tests around the test server/container ...
