package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func main() {
	logfile, err := os.Create("out.log")
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(logfile)

	ctx := context.Background()
	handler := &handler{
		fileCache: make(map[string]string),
	}

	dockerClient, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	images, err := dockerClient.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		panic(err)
	}
	for _, image := range images {
		for _, tag := range image.RepoTags {
			handler.cachedImages = append(handler.cachedImages, tag)
		}
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		os.Exit(0)
	}()

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Println(err)
			continue
		}
		var contentLength uint
		fmt.Sscanf(string(line), "Content-Length: %d", &contentLength)
		_, _, _ = reader.ReadLine()
		buffer := make([]byte, contentLength)

		bufferWindow := buffer
		for len(bufferWindow) > 0 {
			n, err := reader.Read(bufferWindow)
			bufferWindow = bufferWindow[n:]
			if err != nil {
				log.Fatalln(err)
			}
		}

		var rpcCall RpcCall
		if err := json.Unmarshal(buffer, &rpcCall); err != nil {
			synErr := err.(*json.SyntaxError)
			log.Println(synErr.Offset)
			log.Println("Error unmarshalling call", err)
			continue
		}

		response, err := handler.Handle(ctx, rpcCall)
		if err != nil {
			log.Fatalln(err)
		}
		if response == nil {
			continue
		}
		responseBytes, err := json.Marshal(RpcResponse{
			Id:     rpcCall.Id,
			Reuslt: response,
		})
		log.Println("response", string(responseBytes))
		if err != nil {
			log.Fatalln(err)
		}
		os.Stdout.WriteString(fmt.Sprintf("Content-Length: %d\r\n\r\n", len(responseBytes)))
		os.Stdout.Write(responseBytes)
	}
}
