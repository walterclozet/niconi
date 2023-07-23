#!/usr/bin/env bash
rm niconi.zip
rm -r dist/assets/
rm -r dist/static/
rm -r dist/data/
rm dist/niconi
cp -rp assets dist/
cp -rp static dist/
mkdir -p dist/data/allstars.db/
rm dist/static/*en*
rm dist/static/*ko*
env CGO_ENABLED=0 GOOS=android GOARCH=arm64 go build -o dist/niconi main.go
cd dist
zip -r ../niconi.zip *
