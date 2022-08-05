all:
	go build -o ./libs/libneofs-darwin-amd64.so -buildmode=c-shared main.go accounting.go container.go netmap.go object.go reputation.go session.go
