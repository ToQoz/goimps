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
