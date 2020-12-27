# Elastic Search cluster monitor

Creation of a **Go microservice** that monitors the health of an **Elastic Search cluster**.

## Goal

- Use this microservice as a monitoring tool.
- Due date: 08.01.2021

## Specifications

- The microservice should provide gRPC endpoints giving enough information to determine the health and the state of the cluster and its indices.
- The code should rely as much as possible on the standard Go library (avoid fancy libraries).
- *Bonus: packaging of the full solution using a single helm chart*

## Sources

See `docs/sources.md` file.

## Instructions

### Start Elastic Search container and load sample data

- In root project folder (requires `docker-compose`):
```shell
docker-compose up --build
```
- Update docker host address in `script/load-dataset.sh` file if needed (default: localhost:9200)
- In root project folder (requires `curl`)
```shell
./script/load-dataset.sh
```

### Start gRPC Server and Client

- run cmd/server.go
    - optional flag `-server-port` : server listening port (default: 9000)
    - optional flag `-es-addr` : elastic search server address (default: http://localhost:9200)
- run cmd/client.go
    - optional flag `-server-port` : server port (default: 9000)
    
## Execution examples

### Server start-up
```shell
$ go run server.go
2020/12/27 10:35:21 Starting grpc server...
```
### Client start-up
```shell
$ go run client.go
2020/12/27 10:35:23 Cluster Health:

cluster_name: "docker-cluster"
status: YELLOW
number_of_nodes: 1
number_of_data_nodes: 1
active_primary_shards: 1
active_shards: 1
unassigned_shards: 1
active_shards_percent_as_number: 50
timestamp: 1609061723


2020/12/27 10:35:23 Indices Info:

indices: <
  health: YELLOW
  status: "open"
  index: "bank"
  uuid: "zFsb_M0gSke3mw0G_rXkkQ"
  pri: 1
  rep: 1
  docs_count: 1000
  store_size: 388312
  pri_store_size: 388312
>
timestamp: 1609061723
```

## Tests

Tests for the Monitoring Service can be executed with the following command from the project root folder:

```shell
go test ./pkg/service
```


