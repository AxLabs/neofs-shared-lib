#!/usr/bin/env bash

echo "#############"
echo "darwin/arm64:"
echo "#############"
docker run --rm \
  -v `pwd`:/go/src/github.com/user/go-project \
  -w /go/src/github.com/user/go-project \
  -e CGO_ENABLED=1 \
  docker.elastic.co/beats-dev/golang-crossbuild:1.19.3-darwin-arm64-debian10 \
  --build-cmd "go mod init; go mod tidy; go build -o ./libs/libneofs-darwin-arm64.so -buildmode=c-shared main.go lib.go parser.go util.go" \
  -p "darwin/arm64" || {
    # Fallback: if the docker build fails for darwin/arm64, then try it to build natively...
    MACHINE_TYPE=`uname -m`
    if [[ "$OSTYPE" == "darwin"* ]]; then
          if [ "${MACHINE_TYPE}" == "arm64" ]; then
            go build -o ./libs/libneofs-darwin-arm64.so -buildmode=c-shared main.go lib.go parser.go util.go && echo 'Go build: success.' || { echo 'ERROR: building shared lib failed. Exiting...'; exit 1; }
          else
            echo 'ERROR: failed to build shared lib locally *and* with docker. Exiting...'; exit 1;
          fi
    else
      echo 'ERROR: failed to build shared lib locally *and* with docker. Exiting...'; exit 1;
    fi
  }

echo "#############"
echo "darwin/amd64:"
echo "#############"
docker run --rm \
  -v `pwd`:/go/src/github.com/user/go-project \
  -w /go/src/github.com/user/go-project \
  -e CGO_ENABLED=1 \
  docker.elastic.co/beats-dev/golang-crossbuild:1.19.3-darwin-debian10 \
  --build-cmd "go build -o ./libs/libneofs-darwin-amd64.so -buildmode=c-shared main.go lib.go parser.go util.go" \
  -p "darwin/amd64" || { echo 'ERROR: building shared lib failed. Exiting...'; exit 1; }

echo "#############"
echo "linux/amd64:"
echo "#############"
docker run --rm \
  -v `pwd`:/go/src/github.com/user/go-project \
  -w /go/src/github.com/user/go-project \
  -e CGO_ENABLED=1 \
  docker.elastic.co/beats-dev/golang-crossbuild:1.19.3-main-debian10 \
  --build-cmd "go build -o ./libs/libneofs-linux-amd64.so -buildmode=c-shared main.go lib.go parser.go util.go" \
  -p "linux/amd64" || { echo 'ERROR: building shared lib failed. Exiting...'; exit 1; }

echo "#############"
echo "linux/386:"
echo "#############"
docker run --rm \
  -v `pwd`:/go/src/github.com/user/go-project \
  -w /go/src/github.com/user/go-project \
  -e CGO_ENABLED=1 \
  docker.elastic.co/beats-dev/golang-crossbuild:1.19.3-main-debian10 \
  --build-cmd "go build -o ./libs/libneofs-linux-i386.so -buildmode=c-shared main.go lib.go parser.go util.go" \
  -p "linux/386" || { echo 'ERROR: building shared lib failed. Exiting...'; exit 1; }

echo "#############"
echo "linux/arm64:"
echo "#############"
docker run --rm \
  -v `pwd`:/go/src/github.com/user/go-project \
  -w /go/src/github.com/user/go-project \
  -e CGO_ENABLED=1 \
  docker.elastic.co/beats-dev/golang-crossbuild:1.19.3-arm-debian10 \
  --build-cmd "go build -o ./libs/libneofs-linux-arm64.so -buildmode=c-shared main.go lib.go parser.go util.go" \
  -p "linux/arm64" || { echo 'ERROR: building shared lib failed. Exiting...'; exit 1; }

echo "#############"
echo "linux/armv7:"
echo "#############"
docker run --rm \
  -v `pwd`:/go/src/github.com/user/go-project \
  -w /go/src/github.com/user/go-project \
  -e CGO_ENABLED=1 \
  docker.elastic.co/beats-dev/golang-crossbuild:1.19.3-armhf-debian10 \
  --build-cmd "go build -o ./libs/libneofs-linux-armhf.so -buildmode=c-shared main.go lib.go parser.go util.go" \
  -p "linux/armv7" || { echo 'ERROR: building shared lib failed. Exiting...'; exit 1; }

echo "#############"
echo "windows/amd64:"
echo "#############"
docker run --rm \
  -v `pwd`:/go/src/github.com/user/go-project \
  -w /go/src/github.com/user/go-project \
  -e CGO_ENABLED=1 \
  docker.elastic.co/beats-dev/golang-crossbuild:1.19.3-main-debian10 \
  --build-cmd "go build -o ./libs/libneofs-windows-amd64.so -buildmode=c-shared main.go lib.go parser.go util.go" \
  -p "windows/amd64" || { echo 'ERROR: building shared lib failed. Exiting...'; exit 1; }


echo "#############"
echo "windows/386:"
echo "#############"
docker run --rm \
  -v `pwd`:/go/src/github.com/user/go-project \
  -w /go/src/github.com/user/go-project \
  -e CGO_ENABLED=1 \
  docker.elastic.co/beats-dev/golang-crossbuild:1.19.3-main-debian10 \
  --build-cmd "go build -o ./libs/libneofs-windows-386.so -buildmode=c-shared main.go lib.go parser.go util.go" \
  -p "windows/386" || { echo 'ERROR: building shared lib failed. Exiting...'; exit 1; }

# TODO: windows/arm64 is not yet supported by https://github.com/elastic/golang-crossbuild
# Check out: https://github.com/elastic/golang-crossbuild/issues/258
#docker run --rm \
#  -v `pwd`:/go/src/github.com/user/go-project \
#  -w /go/src/github.com/user/go-project \
#  -e CGO_ENABLED=1 \
#  docker.elastic.co/beats-dev/golang-crossbuild:1.19.3-main-debian10 \
#  --build-cmd "go build -o ./libs/libneofs-windows-arm64.so -buildmode=c-shared main.go lib.go parser.go util.go" \
#  -p "windows/arm64" || { echo 'ERROR: building shared lib failed. Exiting...'; exit 1; }
