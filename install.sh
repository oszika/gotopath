#!/bin/sh

mkdir -p ~/.config/gotopath
mkdir -p ~/.config/systemd/user

echo "Install gotopath"
go install

echo "Install systemd gotopath service"
perl -pe 's/#GOTOPATH#/$ENV{GOPATH}\/bin\/gotopath/' config/systemd/user/gotopath.service > ~/.config/systemd/user/gotopath.service
