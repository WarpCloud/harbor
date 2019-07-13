#!/bin/bash

 make install GOBUILDIMAGE=golang:1.11.2 COMPILETAG=compile_core


  docker build -t 172.16.1.99/tmp/harbor-core:v1.8.1-debug -f /home/qls/code/golang/src/github.com/goharbor/harbor/make/photon/core/Dockerfile  .