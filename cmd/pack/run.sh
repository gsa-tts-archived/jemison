#!/bin/bash

export SERVICE=pack



pushd /home/vcap/app/cmd/${SERVICE}
    # Mount the serve bucket as a place to write
    # SQLite files. We go straight to S3.s3
    # mkdir -p s3
    # rm -f s3/*
    # echo "mounting serve bucket as s3fs"
    # s3fs serve s3 \
    #   -o url=http://localhost:9000/ \
    #   -o passwd_file=./.miniocreds \
    #   -o use_path_request_style \
    #   -o dbglevel=info
    echo Running the $SERVICE
    ./service.exe
popd