package services

import (
	"context"
	"fmt"
	"grpc/bp/hello"
	"io"
	"log"
)

type HelloService struct {
	hello.UnimplementedHelloServiceServer
}

func (HelloService) Say(ctx context.Context, request *hello.HelloRequest) (*hello.HelloResponse, error) {
	log.Printf("say: %s", request.Name)
	return &hello.HelloResponse{Reply: fmt.Sprintf("say.reply: %s", request.Name)}, nil
}

func (HelloService) List(request *hello.HelloRequest, stream hello.HelloService_ListServer) error {
	log.Printf("list: %s", request.Name)
	for i := 0; i < 10; i++ {
		err := stream.Send(&hello.HelloResponse{Reply: fmt.Sprintf("list.reply: %s.%d", request.Name, i)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (HelloService) Multipart(stream hello.HelloService_MultipartServer) error {
	log.Print("multipart")
	cnt := 0
	for {
		req, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		log.Printf("multipart: %s", req.Name)
		cnt++
	}
	return stream.SendAndClose(&hello.HelloResponse{Reply: fmt.Sprintf("multipart.reply: %d", cnt)})
}

func (HelloService) Channel(stream hello.HelloService_ChannelServer) error {
	log.Print("channel")
	cnt := 0
	for {
		req, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		log.Printf("channel: %s", req.Name)
		cnt++
		err = stream.Send(&hello.HelloResponse{Reply: fmt.Sprintf("channel.reply: %d", cnt)})
		if err != nil {
			return err
		}
	}
	return nil
}
