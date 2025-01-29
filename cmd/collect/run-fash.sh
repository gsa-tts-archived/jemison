#!/bin/bash

export SERVICE=collect

pushd /home/vcap/app/cmd/${SERVICE}
    echo Running the $SERVICE
    make run
popd