_default:
  @just --list --list-prefix '  > '

build-all:
  just build-server
  just build-client

build-server:
  docker build --no-cache . -t censys-kv-server:0.1.0 -f ./kv-service/Dockerfile

build-client:
  docker build . -t censys-kv-client:0.1.0 -f ./kv-test-client/Dockerfile

run:
  docker network inspect censys-kv-network >/dev/null 2>&1 || docker network create --driver bridge censys-kv-network
  docker compose up

down:
  docker compose down
