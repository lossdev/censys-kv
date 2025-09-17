# censys-kv

## Purpose

This project contains 2 go modules and Docker containers for each service. The KV service runs in the background and supports adding, retrieving, and deleting keys, while a test client also runs and waits to accept calls to either test overwriting keys (`/test_overwrite`) or testing deletion (`/test_deletion`).

### Installation and Setup

You should have docker installed on your system. A [justfile](https://github.com/casey/just) is provided for you to run the `docker build` commands easily via `just build-all`, but just in case you don't have it or don't want to install it, I'll include copy/paste build steps below to build the two containers. Please issue these commands at the root of the repo.

First, `kv-service`:
```sh
docker build --no-cache . -t censys-kv-server:0.1.0 -f ./kv-service/Dockerfile
```

And secondly, `kv-client`:
```sh
docker build . -t censys-kv-client:0.1.0 -f ./kv-test-client/Dockerfile
```

### Running

There is an included `docker-compose` file at the root of the repo. You can either run `just run`, or `docker compose up` in order to start the stack. `kv-service` will be available at :8080, and `kv-client` will be available at :8081.

To test the completeness of the `kv-service`, simply run these commands:

```sh
curl -X GET -o - http://localhost:8081/test_overwrite
```

And:
```sh
curl -X GET -o - http://localhost:8081/test_deletion
```

### Documentation

#### `kv-service`

The first is the `kv-service`, which runs a gin HTTP server that contains a simple KV store (string:string). It supports:

##### `PUT /key/:key/:value`

example:

`curl -X PUT -o - http://localhost:8080/key/foo/bar`

Will create a key "foo" with the value "bar". Since this is a `PUT` operation, keys can be overwritten at will.

##### `DELETE /key/:key`

example:

`curl -X DELETE -o - http://localhost:8080/key/foo`

Will delete a key "foo". If the key does not exist, there will give you an error stating that it couldn't find a key with that name.

##### `GET /key/:key`

example:

`curl -X GET -o - http://localhost:8080/key/foo`

Will retrieve a key "foo". Similar to DELETE, if the key does not exist, it will give you an error stating so.

#### `kv-client`

The second service is the `kv-client`, which ensures a testing completeness of the server implementation. The two endpoints sufiiciently test all 3 functionalities of the KV server, as test_deletion tests deletion while test_overwrite tests both adding and getting elements.

##### `GET /test_deletion`

example:

`curl -X GET -o - http://localhost:8081/test_deletion`

Will issue a PUT request then a DELETE request on the same key.

##### `GET /test_overwrite`

example:

`curl -X GET -o - http://localhost:8081/test_overwrite`

Will issue a PUT and GET to ensure the first value, then another PUT and GET to ensure that the first value is overwritten by the second.