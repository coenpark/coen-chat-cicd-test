GOOGLEAPIS=third_party/googleapis

create:
	protoc -I. -I${GOOGLEAPIS} \
 				--go_out=. --go_opt paths=source_relative \
				--go-grpc_out=. --go-grpc_opt=paths=source_relative \
				--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
				protos/channel/channel.proto
	protoc -I. -I${GOOGLEAPIS} \
				--go_out=. --go_opt=paths=source_relative \
				--go-grpc_out=. --go-grpc_opt=paths=source_relative \
				--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
				protos/message/message.proto
	protoc -I. -I${GOOGLEAPIS} \
				--go_out=. --go_opt=paths=source_relative \
				--go-grpc_out=. --go-grpc_opt=paths=source_relative \
				--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
				protos/user/user.proto
	protoc -I. -I${GOOGLEAPIS} \
				--go_out=. --go_opt=paths=source_relative \
				--go-grpc_out=. --go-grpc_opt=paths=source_relative \
				--grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
				protos/search/search.proto
