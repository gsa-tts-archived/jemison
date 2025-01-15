#!/bin/bash

export SERVICE=resultsapi

pushd /home/vcap/app/cmd/${SERVICE}
    echo Running the $SERVICE
    make run
popd