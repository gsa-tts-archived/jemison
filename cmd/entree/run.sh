#!/bin/bash

/home/vcap/app/cmd/migrate/service.exe || exit 1
/home/vcap/app/cmd/entree/service.exe || exit 1
