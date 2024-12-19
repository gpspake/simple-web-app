#!/bin/bash

# Run tests using the docker compose test config

docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
