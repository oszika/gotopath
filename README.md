# gotopath

## Synopsis
Go to path is a tool to facilitate navigation in the shell. It's autocomplete suggestions and it learns from himself where are your favorite paths.

## Motivation
It's too boring to handle shell aliases manually.

## How it's work?
A very light daemon manages shortcuts and communicates with client using unix sockets. A client submits path request using shortcut, relative or absolute path. If shortcut is requested and if exists, server return the main used absolute path. If a complete path is used, server just add a count for all subpaths.

## Limitations
For now, it's work only with Zsh. Script has not been tested with other shells.

## Installation
go get github.com/oszika/gotopath && $GOPATH/src/github.com/oszika/gotopath/install.sh

## Zsh integration
Set in your .zshrc:
```
source ~/.config/gotopath/gotopath.zsh
```

### Autocompletion
Set in your .zshrc:
```
fpath=(~/.config/gotopath/ $fpath)
```
All shortcuts appear but also all paths associated. For example:

```
$ g tat<TAB>
Shortcuts:
tata                  tata:=/tmp/tata       tata:=/tmp/titi/tata
```

## Systemd integration
You can start gotopath service using systemctl:
```
$ systemctl --user start gotopath.service
```

## Usage
Go to /etc/zsh using absolute path:
```
$ g /etc/zsh
Go to /etc/zsh
```

Go to /etc/zsh using shortcut:
```
$ g zsh
Go to /etc/zsh

```

Autocomplete z*:
```
$ g z<TAB>
Shortcuts:
zsh            zsh:=/etc/zsh
zprofile
```
