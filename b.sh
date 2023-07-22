#!/usr/bin/env bash
rm niconi.zip
rm -r dist/assets/
rm -r dist/static/
rm dist/niconi
cp -r assets dist/
cp -r static dist/
env CGO_ENABLED=0 GOOS=android GOARCH=arm64 go build -o dist/niconi main.go
zip -r -j niconi.zip dist/*

