package main

import (
	"context"
	"fmt"
	"log"

	"github.com/thewebdevel/grpc-blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// Create Blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Sathish",
		Title:    "My first blog",
		Content:  "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})

	blogID := createBlogRes.Blog.GetId()

	if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
	}

	fmt.Printf("Blog created: %v\n", createBlogRes)

	// Read Blog
	readBlogRes, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blogID})
	if err != nil {
		fmt.Printf("Unable to read blog: %v\n", err)
	}

	fmt.Printf("Blog read: %v\n", readBlogRes)

	// Update Blog
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Swaathi",
		Title:    "Erlang blog",
		Content:  "Content of the Erlang blog",
	}
	updateBlogRes, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if err != nil {
		fmt.Printf("Unable to update blog: %v\n", err)
	}

	fmt.Printf("Blog updated: %v\n", updateBlogRes)
}
