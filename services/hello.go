package services

import (
	"context"
	"fmt"
	"grpc/bp/hello"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type HelloService struct {
	hello.UnimplementedHelloServiceServer
}

func (HelloService) Say(ctx context.Context, request *hello.HelloRequest) (*hello.HelloResponse, error) {
	defer func() {
		md := metadata.Pairs("type", "trialer.say.reply")
		grpc.SetTrailer(ctx, md)
	}()
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Print(md)
	}
	log.Printf("say: %s", request.Name)

	md := metadata.Pairs("type", "header.say.reply")
	grpc.SendHeader(ctx, md)

	return &hello.HelloResponse{Reply: fmt.Sprintf("say.reply: %s", request.Name)}, nil
}

func (HelloService) List(request *hello.HelloRequest, stream hello.HelloService_ListServer) error {
	defer func() {
		md := metadata.Pairs("type", "trialer.list.reply")
		grpc.SetTrailer(stream.Context(), md)
	}()
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		log.Print(md["type"])
	}

	log.Printf("list: %s", request.Name)

	md := metadata.Pairs("type", "header.list.reply")
	grpc.SendHeader(stream.Context(), md)

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
	defer func() {
		md := metadata.Pairs("type", "trialer.multipart.reply")
		grpc.SetTrailer(stream.Context(), md)
	}()

	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		log.Print(md["type"])
	}

	md := metadata.Pairs("type", "header.multipart.reply")
	grpc.SendHeader(stream.Context(), md)
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

	defer func() {
		md := metadata.Pairs("type", "trialer.channel.reply")
		grpc.SetTrailer(stream.Context(), md)
	}()

	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		log.Print(md["type"])
	}

	md := metadata.Pairs("type", "header.channel.reply")
	grpc.SendHeader(stream.Context(), md)

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
