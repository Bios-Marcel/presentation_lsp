package main

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

// Position is a simple y,x coordinate. Both values start at 0.
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type HoverParams TextDocumentPositionParams

type CompletionParams TextDocumentPositionParams

type TextDocumentPositionParams struct {
	TextDocument TextDocument `json:"textDocument"`
	Position     Position     `json:"position"`
}

type CompletionItem struct {
	Label string `json:"label"`
}

type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

type HoverResult struct {
	Contents string `json:"contents"`
}
