#!/bin/bash

# This script runs in the GH Actions environment.
# Running and testing locally will not work.

apt-get update
apt-get install -y zip python

# Install CF CLI
curl -k -O -L https://github.com/cloudfoundry/cli/releases/download/v8.8.0/cf8-cli-installer_8.8.0_x86-64.deb
apt-get install --assume-yes ./cf8-cli-installer_8.8.0_x86-64.deb

# Authenticate against Cloud.gov
cf api api.fr.cloud.gov
cf auth ${{ secrets.CF_USERNAME }} ${{ secrets.CF_PASSWORD }}

# Run the Terraform for deploying to `dev`
pushd ${GITHUB_WORKSPACE}/terraform
  make dev
popd