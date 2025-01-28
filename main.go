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
	"strings"
	"syscall"
)

func main() {
	logfile, err := os.Create("out.log")
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(logfile)

	ctx := context.Background()
	// listener := NewStdListener()
	handler := &handler{
		fileCache: make(map[string]string),
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
		// FIXME Support different headers
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

type handler struct {
	fileCache map[string]string
}

func (h *handler) Handle(ctx context.Context, req RpcCall) (any, error) {
	log.Println("Incomming request", req.Method)
	log.Println(string(req.Params))
	switch req.Method {
	case "initialize":
		return InitializeResult{
			Capabilities: ServerCapabilities{
				PositionEncoding: PositionEncodingUTF8,
				HoverProvider:    true,
				TextDocumentSync: TextDocumentSyncKindFull,
			},
			ServerInfo: ServerInfo{
				Name:    "presentation_lsp",
				Version: "1",
			},
		}, nil
	case "textDocument/didOpen":
		var openParams TextDocumentOpenParams
		if err := json.Unmarshal(req.Params, &openParams); err != nil {
			return nil, err
		}

		h.fileCache[openParams.TextDocument.URI] = openParams.TextDocument.Text
	case "textDocument/didChange":
		var openParams TextDocumentChangeParams
		if err := json.Unmarshal(req.Params, &openParams); err != nil {
			return nil, err
		}

		h.fileCache[openParams.TextDocument.URI] = openParams.ContentChanges[0].Text
	case "textDocument/hover":
		var hoverParams HoverParams
		if err := json.Unmarshal(req.Params, &hoverParams); err != nil {
			return nil, err
		}

		log.Println("Cache", h.fileCache)
		file := h.fileCache[hoverParams.TextDocument.URI]
		lines := strings.Split(file, "\n")
		log.Println(lines)
		line := lines[hoverParams.Position.Line]
		log.Println(len(line))
		pos := hoverParams.Position.Character
		char := line[pos]
		if char == ' ' {
			return HoverResult{}, nil
		}

		var word string
		startOfWord := strings.LastIndexByte(line[:pos], ' ')
		endOfWord := strings.IndexByte(line[pos:], ' ')
		log.Println(startOfWord, endOfWord)
		if startOfWord == -1 && endOfWord == -1 {
			word = line
		} else if startOfWord == -1 {
			word = line[:endOfWord+int(pos)]
		} else if endOfWord == -1 {
			word = line[startOfWord+1:]
		} else {
			word = line[startOfWord+1 : endOfWord+pos]
		}

		word = strings.TrimSpace(word)
		switch word {
		case "services:":
			log.Println("Match:", word)
			return HoverResult{
				Contents: "`services` holds a map of all services. It can contain 0 to n values.",
			}, nil
		case "image:":
			log.Println("Match:", word)
			return HoverResult{
				Contents: "`image` defines which Docker image is deployed.",
			}, nil
		}
	}
	return nil, nil
}
