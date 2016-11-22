#!/bin/sh

mkdir -p ~/.config/gotopath
mkdir -p ~/.config/systemd/user

echo "Install gotopath"
go install

echo "Install systemd gotopath service"
perl -pe 's/#GOTOPATH#/$ENV{GOPATH}\/bin\/gotopath/' config/systemd/user/gotopath.service > ~/.config/systemd/user/gotopath.service
echo -e "\tTo enable: systemctl --user start gotopath.service"

echo "Install zsh functions"
cp config/zsh/_gotopath ~/.config/gotopath/_gotopath
cp config/zsh/gotopath.zsh ~/.config/gotopath/gotopath.zsh
echo -e "\tTo enable: source ~/.config/gotopath/gotopath.zsh"
echo -e "\tAutocompletion: fpath=(~/.config/gotopath/ \$fpath)"
