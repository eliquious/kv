package main

import (
	"testing"

	"github.com/eliquious/sandbox/kv/kv-proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func BenchmarkGet(b *testing.B) {
	conn, err := grpc.Dial(":9034", grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalf("dial error: %s\n", err)
	}
	defer conn.Close()
	ctx := context.Background()
	client := kv.NewKeyValueServiceClient(conn)
	client.Set(ctx, &kv.KVPair{[]byte("key"), []byte("value")})
	b.StartTimer()
	for index := 0; index < b.N; index++ {
		_, err = client.Get(ctx, &kv.Key{[]byte("key")})
		if err != nil {
			b.Fatalf("Get error: %s\n", err)
		}
	}
	b.StopTimer()
}

func BenchmarkSet(b *testing.B) {
	conn, err := grpc.Dial(":9034", grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalf("dial error: %s\n", err)
	}
	defer conn.Close()

	ctx := context.Background()
	client := kv.NewKeyValueServiceClient(conn)
	b.StartTimer()
	for index := 0; index < b.N; index++ {
		_, err = client.Set(ctx, &kv.KVPair{[]byte("key"), []byte("value")})
		if err != nil {
			b.Fatalf("Set error: %s\n", err)
		}
	}
	b.StopTimer()
}
