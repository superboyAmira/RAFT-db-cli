VER_CLUSTER=v1

gen_cluster:
	protoc --go_opt=paths=import --go_out=.  --go-grpc_opt=paths=import --go-grpc_out=. ./api/proto/cluster-contract.proto

protobuf_install : 
	sudo apt-get install -y protobuf-compiler
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	export GO_PATH=~/go
	export PATH=$PATH:$(go env GOPATH)/bin
