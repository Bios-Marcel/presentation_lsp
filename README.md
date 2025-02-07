## Slides

Slides are located at [SLIDES.md](/SLIDES.md). They can be used with
[presenterm](https://github.com/mfontanini/presenterm).

## The LSP implementation

The implementation is meant for docker-compose files. It's not a full language
server implementation, instead it does super basic documentation via hover and
local docker registry image autocompletion.

### Usage in Neovim

```
lua vim.lsp.start({cmd={"./lsp"},root_dir=vim.fn.getcwd()})
```

