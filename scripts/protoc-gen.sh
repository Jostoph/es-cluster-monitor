protoc --proto_path=api/proto --proto_path=scripts --go_opt=module=github.com/Jostoph/es-cluster-monitor/pkg/api --go_out=plugins=grpc:pkg/api es-monitor-service.proto