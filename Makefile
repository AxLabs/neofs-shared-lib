all:
	go build -o ./libs/libneofs-darwin-amd64.so -buildmode=c-shared main.go lib.go parser.go util.go