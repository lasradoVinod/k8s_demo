/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"time"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	for true {
		func() {
			conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			c := pb.NewGreeterClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 300 * time.Second)
			defer cancel()
			for i := 1; i <= 100; i++ {
				r, err := c.SayHello(ctx, &pb.HelloRequest{Name: strconv.Itoa(i)})
				if err != nil {
					log.Fatalf("could not greet: %v", err)
				}
				time.Sleep(100 * time.Millisecond)

				log.Printf("Greeting: %s", r.GetMessage())
			}
			conn.Close()
		} ()
	}
}
