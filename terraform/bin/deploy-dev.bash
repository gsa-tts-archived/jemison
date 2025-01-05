#!/bin/bash

# This script runs in the GH Actions environment.
# Running and testing locally will not work.


# Run the Terraform for deploying to `dev`
pushd ${GITHUB_WORKSPACE}/terraform
  make dev
popd