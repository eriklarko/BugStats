if [ ! -d "src" ]; then
  echo "You don't seem to be standing at root of the TaggyGo project. It should contain a folder named src";
  exit 1;
fi

PWD=`pwd`
DEPS_PATH="${PWD}"
echo "export GOPATH=$DEPS_PATH"
