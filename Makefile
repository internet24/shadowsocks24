.PHONY: install run build

install:
	./third_party/outline-ss-server.sh

run: install
	go run main.go start

build: install
	go build main.go -o shadowsocks24
