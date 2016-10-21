#!/bin/bash

echo "Begin to build go develop env"

if [ ! -f "go1.7.linux-amd64.tar.gz" ];then
echo "download go pkg"
wget http://www.golangtc.com/static/go/1.7/go1.7.linux-amd64.tar.gz
fi

if [ ! -d "go" ];then
tar zxvf go1.7.linux-amd64.tar.gz
fi

if [ ! -d "workspace-go" ];then
mkdir workspace-go
fi

sed -e '/\(GOROOT\|GOPATH\|GOBIN\)/d' ~/.bash_profile  > .bash_profile.temp
mv .bash_profile.temp ~/.bash_profile

echo "export GOROOT=$PWD/go" >> ~/.bash_profile
echo "export PATH=\$PATH:\$GOROOT/bin" >> ~/.bash_profile
echo "export GOPATH=$PWD/workspace-go" >> ~/.bash_profile
echo "export GOBIN=\$GOPATH/bin" >> ~/.bash_profile

