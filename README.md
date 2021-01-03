# Elastic Search cluster monitor

Creation of a **Go microservice** that monitors the health of an **Elastic Search cluster**.

## Goal

- Use this microservice as a monitoring tool.

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
{"level":"INFO","ts":"2021-01-03T14:30:31.889+0100","msg":"Starting grpc server."}
{"level":"INFO","ts":"2021-01-03T14:30:40.117+0100","msg":"finished unary call with code OK",
"grpc.start_time":"2021-01-03T14:30:40+01:00","system":"grpc","span.kind":"server",
"grpc.service":"proto.MonitorService","grpc.method":"ReadClusterHealth","grpc.code":"OK",
"grpc.time_ms":9.640000343322754}
{"level":"INFO","ts":"2021-01-03T14:30:40.127+0100","msg":"finished unary call with code OK",
"grpc.start_time":"2021-01-03T14:30:40+01:00","system":"grpc","span.kind":"server",
"grpc.service":"proto.MonitorService","grpc.method":"ReadIndicesInfo","grpc.code":"OK",
"grpc.time_ms":9.831999778747559}
^C{"level":"WARN","ts":"2021-01-03T14:31:21.327+0100","msg":"Stopping grpc server."}
```
### Client start-up
```shell
$ go run client.go
2021/01/03 14:30:40 Cluster Health:

cluster_name: "docker-cluster"
status: YELLOW
number_of_nodes: 1
number_of_data_nodes: 1
active_primary_shards: 1
active_shards: 1
unassigned_shards: 1
active_shards_percent_as_number: 50
timestamp: 1609680640


2021/01/03 14:30:40 Indices Info:

indices: <
  health: YELLOW
  status: "open"
  index: "bank"
  uuid: "5_7liZyeQRWv2X2zB0Qn8w"
  pri: 1
  rep: 1
  docs_count: 1000
  store_size: 388312
  pri_store_size: 388312
>
timestamp: 1609680640
```

## Tests

Tests for the Monitoring Service can be executed with the following command from the project root folder:

```shell
go test ./pkg/service
```


