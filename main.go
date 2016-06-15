package main

import (
	"io"
	"net"
	"runtime"

	"golang.org/x/net/context"

	"github.com/coocood/freecache"
	"github.com/eliquious/sandbox/kv/kv-proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	runtime.GOMAXPROCS(8)
	lis, err := net.Listen("tcp", ":9034")
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	cache := freecache.NewCache(1024 * 1024 * 256)

	grpcServer := grpc.NewServer()
	kv.RegisterKeyValueServiceServer(grpcServer, &server{cache})
	grpcServer.Serve(lis)
}

type server struct {
	cache *freecache.Cache
}

func (s *server) Get(ctx context.Context, k *kv.Key) (*kv.Value, error) {
	select {
	case <-ctx.Done():
		return &kv.Value{}, ctx.Err()
	default:
		val, err := s.cache.Get(k.Data)
		if err != nil {
			return &kv.Value{}, err
		}
		return &kv.Value{Data: val}, nil
	}
}

func (s *server) Set(ctx context.Context, k *kv.KVPair) (*kv.Value, error) {
	select {
	case <-ctx.Done():
		return &kv.Value{}, ctx.Err()
	default:
		err := s.cache.Set(k.Key, k.Value, 0)
		if err != nil {
			return &kv.Value{}, err
		}
		return &kv.Value{Data: k.Value}, nil
	}
}

func (s *server) GetStream(stream kv.KeyValueService_GetStreamServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		val, _ := s.cache.Get(in.Data)
		if err := stream.Send(&kv.Value{val}); err != nil {
			return err
		}
	}
}

func (s *server) SetStream(stream kv.KeyValueService_SetStreamServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		s.cache.Set(in.Key, in.Value, 30)
		if err := stream.Send(&kv.Value{in.Value}); err != nil {
			return err
		}
	}
}
