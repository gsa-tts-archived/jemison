#!/bin/bash

if [ -z "${SERVICE_KEY_NAME}" ]; then
  echo "You need to set the env var SERVICE_KEY_NAME"
  exit
fi

SERVICE_INSTANCE_NAME=tf_state

cf create-service-key "${SERVICE_INSTANCE_NAME}" "${SERVICE_KEY_NAME}"
S3_CREDENTIALS=$(cf service-key "${SERVICE_INSTANCE_NAME}" "${SERVICE_KEY_NAME}" | tail -n +2)

export AWS_ACCESS_KEY_ID=$(echo "${S3_CREDENTIALS}" | jq -r '.credentials.access_key_id')
export AWS_SECRET_ACCESS_KEY=$(echo "${S3_CREDENTIALS}" | jq -r '.credentials.secret_access_key')
export BUCKET_NAME=$(echo "${S3_CREDENTIALS}" | jq -r '.credentials.bucket')
export AWS_DEFAULT_REGION=$(echo "${S3_CREDENTIALS}" | jq -r '.credentials.region')

# cf delete-service-key -f -w "${SERVICE_INSTANCE_NAME}" "${SERVICE_KEY_NAME}"
