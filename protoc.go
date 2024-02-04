package pssm

//go:generate bash -c "mkdir -p gen && protoc -I protos --go_out gen --go_opt paths=source_relative --go-grpc_out gen --go-grpc_opt paths=source_relative api.proto"
