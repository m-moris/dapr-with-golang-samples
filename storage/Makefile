PORT			:= 5005
HTTP_PORT		:= 3006
GRPC_PORT		:= 50001

run:
	dapr run \
	  --app-id storage-test \
	  --app-protocol grpc \
	  --app-port $(PORT)  \
	  --dapr-grpc-port $(GRPC_PORT) \
	  --dapr-http-port $(HTTP_PORT) \
	  --log-level info \
	  --components-path ./config/ \
	  go run main.go

