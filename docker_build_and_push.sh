#!/bin/bash

docker login
docker build -t spiderswarm/spiderswarm:latest .
docker push spiderswarm/spiderswarm:latest
