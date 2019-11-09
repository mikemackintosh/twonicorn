#!/usr/bin/env bash
FILES=$(ls -d cmd/twonicorn/*.go | grep -v "_test")
for arch in 386 amd64; do
  for os in linux darwin windows; do
    echo "Building for $os/$arch"
    OUTFILE=bin/twonicornd-$os-$arch
    GOARCH=$arch GOOS=$os go build -v -o $OUTFILE $FILES;
    du -hs $OUTFILE
  done
done
