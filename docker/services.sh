#!/usr/bin/env bash
cd docker && docker-compose exec services /bin/sh -c "cd cmd/$1 && TZ=UTC && go run main.go ${param} $2"
