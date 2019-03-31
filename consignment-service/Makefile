build:
	# protoc -I. --go_out=plugins=grpc:. proto/consignment/consignment.proto
	protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. proto/consignment/consignment.proto
	# GOOS=linux GOARCH=amd64 go build
	# docker build -t consignment-service .
	# docker run -d -rm -p 50051:50051 consignment-service