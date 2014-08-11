# Vim plugins for goimps

## Installation

Add your vimrc

```vim
set rtp+=$GOPATH/src/github.com/ToQoz/goimps/misc/vim
```

## Providings

### Functions

```
goimps#Importable()       :: -> []string
goimps#Dropable(filename) :: string -> []string
goimps#Unused(filename)   :: string -> []string
```

## Tips

### Use :DropUnused

[ToQoz/unite-go-imports](http://github.com/ToQoz/unite-go-imports)

### Import/Drop by selecting unite.vim

[ToQoz/unite-go-imports](http://github.com/ToQoz/unite-go-imports)

### Use `goimps fmt` instead of `gofmt`

If you are [vim-jp/go-vim](http://github.com/vim-jp/go-vim) user.

```vim
g:gofmt_command = "goimps fmt"
```

If you are [fatih/vim-go](http://github.com/fatih/vim-go) user.

```vim
g:go_fmt_command = "goimps fmt"
```
