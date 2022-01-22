gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:pb 
clean:
	rm pb/*.go
run:
	go run cmd/main.go
client-call:
	go run cmd/client/main.go
server:
	go run cmd/server/main.go
stream-client:
	go run cmd/streamclient/main.go