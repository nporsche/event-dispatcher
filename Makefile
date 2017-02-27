pb:
	protoc --gofast_out=plugins=grpc:. *.proto

.PHONY:
	pb
