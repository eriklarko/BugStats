#! /bin/bash

if [ ! -d "src" ]; then
  echo "You don't seem to be standing at root of the TaggyGo project. It should contain a folder named src";
  exit 1;
fi

set -e

$(./setup-gopath.sh)
echo "GOPATH is now $GOPATH"

echo "Downloading dependencies..."
go get github.com/codeskyblue/go-sh
echo "Done!"
