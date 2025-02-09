## Slides

Slides are located at [SLIDES.md](/SLIDES.md). They can be used with
[presenterm](https://github.com/mfontanini/presenterm).

## Demo

Small terminal recording in case you want to see what it looks like:

https://asciinema.org/a/eW3vTISGvRF7gJMlTKiIMhhZV

## The LSP implementation

The implementation is meant for docker-compose files. It's not a full language
server implementation, instead it does super basic documentation via hover and
local docker registry image autocompletion.

### Usage in Neovim

To start once in the current buffer:

```
lua vim.lsp.start({cmd={"./lsp"},root_dir=vim.fn.getcwd()})
```

To permanently set it up:

```lua
local lspconfig = require("lspconfig")
require("lspconfig.configs").dclsp = {
  default_config = {
    cmd = { "./lsp" },
    filetypes = { "yaml" },
    root_dir = function(fname)
      return lspconfig.util.find_git_ancestor(fname) or vim.fn.getcwd()
    end,
    settings = {},
  },
}

lspconfig.dclsp.setup({
  on_attach = on_attach,
})
```

