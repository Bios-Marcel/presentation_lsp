---
title: Understanding and implementing the Language Server protocol
theme:
  name: terminal-light
  override:
    footer:
      style: template
      center: https://github.com/Bios-Marcel/presentation_lsp
      right: "{current_slide} / {total_slides}"
---

A brief and incomplete history of code editing
==

<!-- pause --> 
1. Basic Editors
  * Edit, leave editor, compile, read compiler warnings, find cause in code
  * No refactoring features and no LLM 😂
  * People read documentation in books
<!-- pause -->
2. Integrated Development Environments (IDEs)
  * The editor and the compiler run alongside eachother in the IDE
  * Warnings and errors are displayed in the IDE and maybe even inline
  * The code powering the editor was specifically built for the IDE and the
  language
  * Tons of convenience features
<!-- pause -->
2. Editors with LSP
  * The editor has no understanding of any concrete language, the concrete
  features are implemented by a language server. The features may vary and
  contain things such as:
    * formatting
    * autocompletion
    * refactoring features
    * ...
  * The editor can display information such as warnings and errors reported by
  the LSP
  * Most actions can happen asynchronously

<!-- end_slide -->

Where did LSP come from?
==

<!-- pause --> 
Microsoft released the specification after the VS Code release. They didn't
invent the concept though.

The goal was to have one lightweight editor for potentially everything, without
Microsoft having to maintain explicit language support.

<!-- pause --> 
A nice bonus is, that every editor can use this, such as vim, emacs and others.

<!-- pause -->
Funnily enough, other editors such as Neovim even have better support than
VS Code opinion, as you need to write an extension for VS first.

<!-- end_slide -->

What technologies does LSP build on?
==

<!-- pause --> 
LSP is a protocol on top of JSONRPC:
  * Plaintext (not binary)
  * Line-based (akin to HTTP)
<!-- pause --> 

JSONRPC is a simple RPC (Remote Procedure Call) protocol:
  * Allows arbitrary function calls and responses
  * Allows notifications (functions without a response)

<!-- end_slide -->

The protocol - JSONRPC
==

<!-- pause --> 
Each message has 2 parts, the headers and the content.
<!-- pause --> 

A message may looks something like this:

```
Content-Length: N\r\n
\r\n
{"id": 1, "method": "x", "params": {...}}
```
<!-- pause --> 

Each RPC call, consistents of at least 3 lines, separated by Windows linebreaks.

Any response sent with the same ID as a request, is a response to that request.
These can be sent out of order.

<!-- pause --> 

Example:

```
Content-Length: 58

{"id": 1, "method": "hello", "params": {"name": "Marcel"}}
```

```
Content-Length: 48

{"id": 1, "result": {"message": "Hello Marcel"}}
```

<!-- end_slide -->

The protocol - Handshake (Initialize)
==

<!-- pause --> 
1. Editor starts the server (system process)
<!-- pause --> 
2. Editor attaches to stdin and stdout
<!-- pause --> 
3. Editor sends an `initialze` request
  * Contains the capabilities and preferences of the editor, such as:
    * Encoding
    * Capabilities (Completion, Refactoring, ...)
    * Locales
    * Project root directory
    * ...
<!-- pause --> 
4. LSP responds with `initializeResult`
  * Contains capabilities and server identity information
<!-- pause --> 
5. Both sides remember the capabilities and adjust their behaviour accordingly

<!-- end_slide -->

The protocol - Document Sync
==

LSP supports three different types of syncing. Editor and server have to agree
upon which synchronisation protocol is used.
<!-- pause --> 

### None

This one is the most primitive. It means no sync happens at all, you gotta read
files from disk, never having the latest edits.
<!-- pause --> 

### Incremental

You get the full document once and then get notified of any edits. This is the
**preferred** way in terms of performance.

However, it is harder to program.
<!-- pause --> 

### Full

This is as primitive as `None`, but without disk IO. The editor sends you full
text documents on every edit.

Also easy to program and not quite optimal in terms of performance.

<!-- end_slide -->

The protocol - hover action
==

The `hover` action is used to display information about whatever is under the
cursor. This can be a symbol or anything else that the language server supports.

The requests are quite simple:

```json
"position": {
  "line": 0,
  "character": 2
},
"textDocument": {
  "uri":"file:\/\/\/home\/marcel\/code\/presentation_lsp\/docker-compose.yml"
}
```

The server then responds with plaintext or markdown in the `contents` field:

```json
"contents": "`image` defines which Docker image is deployed."
```

The protocol - completion action
==

The `completion` action is quite similar to `hover`:

```json
"context": {
  "triggerKind": 1
},
"position": {
  "line": 2,
  "character": 12
},
"textDocument": {
  "uri":"file:\/\/\/home\/marcel\/code\/presentation_lsp\/docker-compose.yml"
}
```

Response:

```json
"isIncomplete": false,
"items": [
  {"label":"traefik:latest"},
  {"label":"biosmarcel/scribble.rs:latest"},
  {"label":"biosmarcel/scribble.rs:v0.9.3"},
  {"label":"biosmarcel/scribble.rs:v0.9.2"}
]
```

Future usecases
==

* LLMs
  * Built-in support
  * Expose semantical context via API
* ???


