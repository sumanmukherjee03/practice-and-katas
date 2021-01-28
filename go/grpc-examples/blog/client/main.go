package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/sumanmukherjee03/practice-and-katas/go/grpc-examples/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Starting grpc client for blog service")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not setup connection to grpc server from client : %v", err)
	}
	defer cc.Close()
	c := blogpb.NewBlogServiceClient(cc)

	serialNum := rand.Intn(10000)

	// Create blog
	resp, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: fmt.Sprintf("John Doe %d", serialNum),
			Title:    fmt.Sprintf("Youtube tutorial - %d", serialNum),
			Content:  fmt.Sprintf("Start by learning the basics of how the web works - %d", serialNum),
		},
	})
	if err != nil {
		log.Fatalf("Encountered an error in creating a blog : %v", err)
	}
	createdBlog := resp.GetBlog()
	fmt.Printf("Created blog in the database : [%s]\n", createdBlog)

	// Read blog
	readBlogResp, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: createdBlog.GetId()})
	if err != nil {
		log.Fatalf("Encountered an error in reading a blog : %v", err)
	}
	readBlog := readBlogResp.GetBlog()
	fmt.Printf("Read blog from the database : [%s]\n", readBlog)

	// Update blog
	updateBlogResp, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       readBlog.GetId(),
			AuthorId: fmt.Sprintf("Jane Doe %d", serialNum),
			Title:    fmt.Sprintf("Wordpress blog - %d", serialNum),
			Content:  fmt.Sprintf("Start by learning the basics of how wordpress works - %d", serialNum),
		},
	})
	if err != nil {
		log.Fatalf("Encountered an error in updating a blog : %v", err)
	}
	updatedBlog := updateBlogResp.GetBlog()
	fmt.Printf("Updated blog in the database : [%s]\n", updatedBlog)

	delBlogResponse, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: readBlog.GetId()})
	if err != nil {
		log.Fatalf("Encountered an error in deleting a blog : %v", err)
	}
	fmt.Printf("Deleted blog in the database with id : [%s]\n", delBlogResponse.GetBlogId())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	listBlogReq := &blogpb.ListBlogRequest{}
	stream, listBlogErr := c.ListBlog(ctx, listBlogReq)
	if listBlogErr != nil {
		log.Fatalf("Received an error from the blog server while setting up server streaming : %v", listBlogErr)
	}
	for {
		listBlogResp, listBlogStreamErr := stream.Recv()
		if listBlogStreamErr == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Received an error during streaming of the response : %v", listBlogStreamErr)
		}
		fmt.Printf("Received blog : [%v]\n", listBlogResp.GetBlog())
	}
}
