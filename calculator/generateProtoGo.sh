#!/bin/bash

protoc calculatepb/calculate.proto --go_out=plugins=grpc:.