package main

type ClientCapabilities struct{}

type TextDocumentSyncKind uint8

const TextDocumentSyncKindFull TextDocumentSyncKind = 1

type InitializeParams struct {
	Capabilities ClientCapabilities `json:"capabilities"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type PositionEncoding string

const PositionEncodingUTF8 PositionEncoding = "utf-8"

type CompletionProvider struct{}

type ServerCapabilities struct {
	PositionEncoding   PositionEncoding     `json:"positionEncoding"`
	HoverProvider      bool                 `json:"hoverProvider"`
	CompletionProvider CompletionProvider   `json:"completionProvider"`
	TextDocumentSync   TextDocumentSyncKind `json:"textDocumentSync"`
}
