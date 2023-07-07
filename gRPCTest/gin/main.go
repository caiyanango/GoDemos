package main

import (
	"fmt"
	"gRPCTest/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func main() {
	conn, _ := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	web := gin.Default()
	web.GET("/:name", func(c *gin.Context) {
		name := c.Param("name")
		req := &pb.HelloRequest{Name: name}
		reply, err := client.SayHello(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(reply.Reply),
		})
	})
	web.Run(":8081")
}
