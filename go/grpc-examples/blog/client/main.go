package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sumanmukherjee03/practice-and-katas/go/grpc-examples/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting grpc client for blog service")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not setup connection to grpc server from client : %v", err)
	}
	defer cc.Close()
	c := blogpb.NewBlogServiceClient(cc)
	resp, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "John Doe",
			Title:    "Starting web development",
			Content:  "Start web development by learning the basics of how the web works",
		},
	})
	if err != nil {
		log.Fatalf("Encountered an error in creating a blog : %v", err)
	}
	fmt.Printf("Created blog in the database : [%s]", resp.GetBlog())
}
