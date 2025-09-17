_default:
  @just --list --list-prefix '  > '

build-all:
  just build-server
  just build-client

build-server:
  docker build . -t censys-kv-server:0.1.0 -f ./kv-service/Dockerfile

build-client:
  docker build . -t censys-kv-client:0.1.0 -f ./kv-test-client/Dockerfile

run:
  docker compose up

down:
  docker compose down
