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
  * Edit, Leave editor, compile, read compiler warnings, find cause in code
<!-- pause -->
2. Integrated Development Environments (IDEs)
  * The editor and the compiler run alongside eachother in the IDE
  * Warnings and errors are displayed in the IDE and maybe even inline
  * The code powering the editor was specifically built for the IDE and the
  language
<!-- pause -->
2. Editors with LSP
  * The editor has no understanding of any concrete language, the concrete
  features are implemented by a language server. The features may vary and
  contain things such as:
    * formatting
    * autocompletion
    * refactoring features
    * ...
  * The editor can show warnings and errors
<!-- pause -->

Where did LSP come from?
==

<!-- pause --> 
Microsoft released the specification when it worked on VS Code. The goal was
to have one lightweight editor for potentially everything, without Microsoft
having to maintain everything.

<!-- pause --> 
A nice bonus is, that every editor can use this, such as vim, emacs and others.

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

The protocol - JSONRPC
==

<!-- pause --> 
Each message has 2 parts, the headers and the content.
<!-- pause --> 

A message may looks something like this:

```
Content-Length: N\r\n
\r\n
{"id": 1, "metohd": "x", "params": {...}}
```
<!-- pause --> 

Each RPC call, consistents of at least 3 lines, separated by Windows linebreaks.

The protocol - Sync
==

TODO How does syncing work? Point out different options.

The protocol - hover action
==

The protocol - code action
==

