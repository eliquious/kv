package main

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"sync"

	"sync/atomic"
	"time"

	"github.com/eliquious/sandbox/kv/kv-proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	runtime.GOMAXPROCS(8)
	// conn, err := grpc.Dial("8.34.219.130:9034", grpc.WithInsecure())
	// if err != nil {
	// 	grpclog.Fatalf("dial error: %s\n", err)
	// }
	// defer conn.Close()
	//
	// client := kv.NewKeyValueServiceClient(conn)
	// val, err := client.Get(context.Background(), &kv.Key{[]byte("key")})
	// if err != nil {
	// 	grpclog.Printf("Get error: %s\n", err)
	// } else {
	// 	grpclog.Printf("Get value: %s\n", val.Data)
	// }
	//
	// val, err = client.Set(context.Background(), &kv.KVPair{[]byte("key"), []byte("value")})
	// if err != nil {
	// 	grpclog.Printf("Set error: %s\n", err)
	// } else {
	// 	grpclog.Printf("Set value: %s\n", val.Data)
	// }
	//
	// val, err = client.Get(context.Background(), &kv.Key{[]byte("key")})
	// if err != nil {
	// 	grpclog.Printf("Get error: %s\n", err)
	// } else {
	// 	grpclog.Printf("Get value: %s\n", val.Data)
	// }

	var totalMessages uint64
	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for worker := 0; worker < 32; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := grpc.Dial("8.34.219.130:9034", grpc.WithInsecure())
			if err != nil {
				grpclog.Fatalf("dial error: %s\n", err)
				return
			}
			defer conn.Close()

			client := kv.NewKeyValueServiceClient(conn)

			keybuf := make([]byte, 8)
			for {
				select {
				case <-ctx.Done():
					grpclog.Printf("Context done: %s\n", ctx.Err())
					return
				default:
					current := atomic.LoadUint64(&totalMessages)
					if current >= 1e5 {
						grpclog.Printf("Goal reached: %s\n")
						cancel()
						return
					}
					binary.BigEndian.PutUint64(keybuf, current)
					client.Set(ctx, &kv.KVPair{keybuf, []byte("value")})
					atomic.AddUint64(&totalMessages, 1)
				}
			}
		}()
	}

	var lastMessage uint64
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		case <-time.After(time.Second):
			current := atomic.LoadUint64(&totalMessages)
			if current > lastMessage {
				grpclog.Printf("QPS: %d Total: %d\n", current-lastMessage, totalMessages)
				lastMessage = current
			}
		}
	}
	wg.Wait()
	duration := time.Now().Sub(start)
	fmt.Println("Elapsed: ", duration)
	fmt.Println("Messages Sent: ", totalMessages)
	fmt.Println("ns Per Message: ", uint64(duration.Nanoseconds())/totalMessages)
	fmt.Println("Messages Per Second: ", 1e9/(uint64(duration.Nanoseconds())/totalMessages))

	// start = time.Now()
	// stream, err := client.SetStream(context.Background())
	// if err != nil {
	// 	grpclog.Fatal(err)
	// }
	// waitc := make(chan struct{})
	// go func() {
	// 	for {
	// 		_, err := stream.Recv()
	// 		if err == io.EOF {
	// 			// read done.
	// 			close(waitc)
	// 			return
	// 		}
	// 		if err != nil {
	// 			close(waitc)
	// 			log.Fatalf("Failed to receive a value: %v", err)
	// 			return
	// 		}
	// 		// log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
	// 	}
	// }()
	//
	// keybuf := make([]byte, 8)
	// value := []byte("value")
	// for index := 0; index < int(totalMessages); index++ {
	// 	binary.BigEndian.PutUint64(keybuf, uint64(index))
	// 	if err := stream.Send(&kv.KVPair{keybuf, value}); err != nil {
	// 		log.Fatalf("Failed to send a kvpair: %v", err)
	// 	}
	// }
	// stream.CloseSend()
	// <-waitc
	// duration = time.Now().Sub(start)
	// fmt.Println("\nElapsed: ", duration)
	// fmt.Println("Messages Sent: ", totalMessages)
	// fmt.Println("ns Per Message: ", uint64(duration.Nanoseconds())/totalMessages)
	// fmt.Println("Messages Per Second: ", 1e9/(uint64(duration.Nanoseconds())/totalMessages))
}
