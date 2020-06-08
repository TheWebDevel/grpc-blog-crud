package main

import (
	"context"
	"fmt"
	"io"
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

	// Delete Blog
	deleteBlogRes, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if err != nil {
		fmt.Printf("Unable to delete blog: %v\n", err)
	}

	fmt.Printf("Blog deleted: %v\n", deleteBlogRes)

	// List Blog
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}
