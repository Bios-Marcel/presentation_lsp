package main

type InitializeParams struct {
	Capabilities ClientCapabilities `json:"capabilities"`
}

type TextDocumentSyncKind uint8

const TextDocumentSyncKindFull TextDocumentSyncKind = 1

type TextDocumentOpenParams struct {
	TextDocument TextDocument `json:"textDocument"`
}

type TextDocument struct {
	URI        string `json:"uri"`
	Text       string `json:"text"`
	Version    uint   `json:"version"`
	LanguageId string `json:"languageId"`
}

type ContentChange struct {
	Text string `json:"text"`
}

type TextDocumentChangeParams struct {
	ContentChanges []ContentChange `json:"contentChanges"`
	TextDocument   TextDocument    `json:"textDocument"`
}

type ClientCapabilities struct{}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type PositionEncoding string

const PositionEncodingUTF8 PositionEncoding = "utf-8"

type ServerCapabilities struct {
	PositionEncoding PositionEncoding     `json:"positionEncoding"`
	HoverProvider    bool                 `json:"hoverProvider"`
	TextDocumentSync TextDocumentSyncKind `json:"textDocumentSync"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

// Position is a simple y,x coordinate. Both values start at 0.
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type HoverParams struct {
	TextDocument TextDocument `json:"textDocument"`
	Position     Position     `json:"position"`
}

type HoverResult struct {
	Contents string `json:"contents"`
}
