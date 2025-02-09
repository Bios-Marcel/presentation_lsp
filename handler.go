package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
)

func (h *handler) handleInitialize() (any, error) {
	return InitializeResult{
		Capabilities: ServerCapabilities{
			PositionEncoding:   PositionEncodingUTF8,
			TextDocumentSync:   TextDocumentSyncKindFull,
			HoverProvider:      true,
			CompletionProvider: CompletionProvider{},
		},
		ServerInfo: ServerInfo{
			Name:    "presentation_lsp",
			Version: "1",
		},
	}, nil
}

func (h *handler) handleDidOpen(params json.RawMessage) (any, error) {
	var openParams TextDocumentOpenParams
	if err := json.Unmarshal(params, &openParams); err != nil {
		return nil, err
	}

	h.fileCache[openParams.TextDocument.URI] = openParams.TextDocument.Text
	return nil, nil
}

func (h *handler) handleDidChange(params json.RawMessage) (any, error) {
	var openParams TextDocumentChangeParams
	if err := json.Unmarshal(params, &openParams); err != nil {
		return nil, err
	}

	h.fileCache[openParams.TextDocument.URI] = openParams.ContentChanges[0].Text
	return nil, nil
}

func (h *handler) handleCompletion(params json.RawMessage) (any, error) {
	var completionParams CompletionParams
	if err := json.Unmarshal(params, &completionParams); err != nil {
		return nil, err
	}

	file := h.fileCache[completionParams.TextDocument.URI]
	lines := strings.Split(file, "\n")
	line := lines[completionParams.Position.Line]

	var result CompletionList
	imageIndex := strings.Index(line, "image:")
	if imageIndex != -1 && (imageIndex+6) < completionParams.Position.Character {
		for _, image := range h.cachedImages {
			result.Items = append(result.Items, CompletionItem{
				Label: image,
			})
		}
	}

	return result, nil
}

func (h *handler) handleHover(params json.RawMessage) (any, error) {
	var hoverParams HoverParams
	if err := json.Unmarshal(params, &hoverParams); err != nil {
		return nil, err
	}

	file := h.fileCache[hoverParams.TextDocument.URI]
	lines := strings.Split(file, "\n")
	line := lines[hoverParams.Position.Line]
	pos := hoverParams.Position.Character
	char := line[pos]
	if char == ' ' {
		return HoverResult{}, nil
	}

	var word string
	startOfWord := strings.LastIndexByte(line[:pos], ' ')
	endOfWord := strings.IndexByte(line[pos:], ' ')
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
	return nil, nil
}

type handler struct {
	fileCache    map[string]string
	cachedImages []string
}

func (h *handler) Handle(ctx context.Context, req RpcCall) (any, error) {
	log.Println("Incomming request", req.Method, string(req.Params))
	switch req.Method {
	case "initialize":
		return h.handleInitialize()
	case "textDocument/didOpen":
		return h.handleDidOpen(req.Params)
	case "textDocument/didChange":
		return h.handleDidChange(req.Params)
	case "textDocument/completion":
		return h.handleCompletion(req.Params)
	case "textDocument/hover":
		return h.handleHover(req.Params)
	}
	return nil, nil
}
