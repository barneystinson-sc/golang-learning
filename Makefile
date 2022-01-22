gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:pb 
clean:
	rm pb/*.go
run:
	go run cmd/main.go
streamresponse-client:
	go run cmd/client/main.go
server:
	go run cmd/server/main.go
streamrequest-client:
	go run cmd/imageuploadclient/main.go
customerservice-client:
	go run cmd/customersupportclient/main.go