#!/bin/sh

go-bindata -pkg stormpathweb -o assets.go templates/* assets/* config/*