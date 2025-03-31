#!/bin/sh

# リリース用。
version=$(git tag -l | tail -n1)

export GOOS=windows

go build \
	-a -tags netgo -trimpath \
	-ldflags='-s -w -extldflags="-static" -X main.version='$version' -buildid=' \
	-o ./go-zenhan.exe \
	&& upx --lzma ./go-zenhan.exe
