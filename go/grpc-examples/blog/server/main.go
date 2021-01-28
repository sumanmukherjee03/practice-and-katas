package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/sumanmukherjee03/practice-and-katas/go/grpc-examples/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

////////////////////////// Implementation of the gRPC server /////////////////////////
type server struct{}

func (s *server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()
	data := blogItem{
		AuthorId: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Encountered an error inserting records into the database : %v", err))
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Could not extract an object id from the returned result when trying to insert data into the database"))
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: data.AuthorId,
			Title:    data.Title,
			Content:  data.Content,
		},
	}, nil
}

func (s *server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	blogId := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Encountered an error generating object id for blog from blogId passed in : %v", err))
	}
	data := &blogItem{}
	res := collection.FindOne(context.Background(), bson.M{"_id": oid})
	if err = res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Encountered an error fetching object from database : %v", err))
	}
	resp := &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorId,
			Title:    data.Title,
			Content:  data.Content,
		},
	}
	return resp, nil
}

func (s *server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Encountered an error generating object id for blog from blogId passed in : %v", err))
	}
	data := &blogItem{}
	filter := bson.M{"_id": oid}
	updateDocument := bson.M{
		"$set": bson.M{
			"author_id": blog.GetAuthorId(),
			"title":     blog.GetTitle(),
			"content":   blog.GetContent(),
		},
	}
	after := options.After // TODO : Not sure why this needs to be in a separate line - understand iota better
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	res := collection.FindOneAndUpdate(ctx, filter, updateDocument, opts)
	if err = res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Encountered an error updating object in database : %v", err))
	}
	resp := &blogpb.UpdateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorId,
			Title:    data.Title,
			Content:  data.Content,
		},
	}
	return resp, nil
}

func (s *server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	blogId := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Encountered an error generating object id for blog from blogId passed in : %v", err))
	}
	filter := bson.M{"_id": oid}
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Encountered an error deleting a record from the database : %v", err))
	}
	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Encountered an error deleting a record from the database because the record could not be found : %v", err))
	}
	return &blogpb.DeleteBlogResponse{BlogId: blogId}, nil
}

func (s *server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	ctx := context.Background()
	cursor, err := collection.Find(ctx, primitive.D{{}})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Encountered an error getting cursor from database for items - %v\n", err))
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		data := &blogItem{}
		if fetchErr := cursor.Decode(data); fetchErr != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Encountered an error fetching item from database - %v\n", fetchErr))
		}
		resp := &blogpb.ListBlogResponse{
			Blog: &blogpb.Blog{
				Id:       data.ID.Hex(),
				AuthorId: data.AuthorId,
				Title:    data.Title,
				Content:  data.Content,
			},
		}
		if sendErr := stream.Send(resp); sendErr != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Encountered an error sending items to client - :%v\n", sendErr))
		}
	}
	if cursorErr := cursor.Err(); cursorErr != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Encountered an error in the database cursor - :%v\n", cursorErr))
	}
	return nil
}

////////////////////////// Types for entries in the Database //////////////////////
type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorId string             `bson:"author_id,omitempty"`
	Title    string             `bson:"title,omitempty"`
	Content  string             `bson:"content,omitempty"`
}

////////////////////////// ENTRYPOINT //////////////////////////
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // This is useful for debugging to get the line number in the code where the error occurred

	////////////////////////////////// MongoDB Client setup code ////////////////////////////
	// Create a context with timeout to interact with the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Defer cancel to release resources bound to the context on program exit

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:welcome2mongodb@localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect mongodb client to mongodb server : %v", err)
	}
	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping mongodb server : %v", err)
	}

	// Defer disconnect mongo client from mongodb server
	defer func() {
		fmt.Println("Closing mongodb connection")
		if err = mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to properly disconnect mongodb client from mongodb server : %v", err)
		}
	}()

	collection = mongoClient.Database("blogdb").Collection("blog")

	////////////////////////////////// GRPC Server setup code ////////////////////////////
	// Setup the network listener
	lis, err := net.Listen("tcp", "localhost:50051") // 50051 is the default port for grpc
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	// boilerplate code to create a new grpc server listening on a host:port
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})
	reflection.Register(s) // Register the reflection for server gRPC implementation

	// Run the grpc server in a separate goroutine
	go func() {
		fmt.Println("Starting new grpc server for blog service")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to start blog server : %s", err)
		}
	}()

	// Create a channel to wait for an interrupt (Ctrl+c) to exit gracefully
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt) // This relays only the interrupt signal to the channel
	<-ch                            // Block until a signal is received from the channel
	fmt.Println("Shutting down server")
	s.Stop()    // Stop the server
	lis.Close() // Close the listener
	fmt.Println("Stopped server")
}
