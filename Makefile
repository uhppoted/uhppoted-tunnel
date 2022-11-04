DIST   ?= development
DEBUG  ?= --debug
CMD     = ./bin/uhppoted-tunnel

.PHONY: sass
.PHONY: debug
.PHONY: reset
.PHONY: update
.PHONY: update-release

all: test      \
     benchmark \
     coverage

clean:
	go clean
	rm -rf bin

update:
	go get -u github.com/uhppoted/uhppoted-lib@master
	go mod tidy

update-release:
	go get -u github.com/uhppoted/uhppoted-lib
	go mod tidy

format: 
	go fmt ./...

build: format
	go build -trimpath -o bin/ ./...

test: build
	go test ./...

vet: test
	go vet ./...

lint: vet
	golint ./...

benchmark: build
	go test -count 5 -bench=.  ./system/events

coverage: build
	go test -cover ./...

build-all: vet
	mkdir -p dist/$(DIST)/windows
	mkdir -p dist/$(DIST)/darwin
	mkdir -p dist/$(DIST)/linux
	mkdir -p dist/$(DIST)/arm7
	env GOOS=linux   GOARCH=amd64       GOWORK=off go build -trimpath -o dist/$(DIST)/linux   ./...
	env GOOS=linux   GOARCH=arm GOARM=7 GOWORK=off go build -trimpath -o dist/$(DIST)/arm7    ./...
	env GOOS=darwin  GOARCH=amd64       GOWORK=off go build -trimpath -o dist/$(DIST)/darwin  ./...
	env GOOS=windows GOARCH=amd64       GOWORK=off go build -trimpath -o dist/$(DIST)/windows ./...

release: update-release build-all
	find . -name ".DS_Store" -delete
	tar --directory=dist --exclude=".DS_Store" -cvzf dist/$(DIST).tar.gz $(DIST)
	cd dist;  zip --recurse-paths $(DIST).zip $(DIST)

debug: build
	# go test -run Test ./...
	# $(CMD) --console --in http/0.0.0.0:8082 --out udp/broadcast:192.168.1.255:60004 --udp-timeout 30s --html ./examples/html
	# $(CMD) --debug --console --in tls/client:127.0.0.1:12345 --out udp/broadcast:192.168.1.255:60005 --udp-timeout 1s --max-retries 2
	# npx eslint --fix ./examples/html/javascript/*.js
	# $(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out tcp/client:216.128.182.157:8080
	$(CMD) --config "./examples/uhppoted-tunnel.toml::client"

delve: build
#   dlv exec ./bin/uhppoted-tunnel -- --debug --console
	dlv test github.com/uhppoted/uhppoted-tunnel -- run Test*

version: build
	$(CMD) version

help: build
	$(CMD) help
	$(CMD) help commands
	$(CMD) help version
	$(CMD) help help

host: build
	$(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out tcp/server:0.0.0.0:12345

client: build
	$(CMD) --debug --console --in tcp/client:127.0.0.1:12345 --out udp/broadcast:192.168.1.255:60005 --udp-timeout 1s

reverse-host: build
	$(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out tcp/client:127.0.0.1:12345

reverse-client: build
	$(CMD) --debug --console --in tcp/server:0.0.0.0:12345 --out udp/broadcast:192.168.1.255:60005 --udp-timeout 1s

tls-host: build
	$(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out tls/server:0.0.0.0:12345 --client-auth

tls-client: build
	$(CMD) --debug --console --in tls/client:127.0.0.1:12345 --out udp/broadcast:192.168.1.255:60005 --udp-timeout 1s

http: build
	npx eslint --fix ./examples/html/javascript/*.js
	$(CMD) --debug --console --in http/0.0.0.0:8082 --out udp/broadcast:192.168.1.255:60005 --udp-timeout 1s --html ./examples/html

https: build
	npx eslint --fix ./examples/html/javascript/*.js
	$(CMD) --debug --console --in https/0.0.0.0:8443 --out udp/broadcast:192.168.1.255:60005 --udp-timeout 1s --html ./examples/html

daemonize: build
	sudo $(CMD) daemonize --in  udp/listen:0.0.0.0:60000  --out tcp/server:0.0.0.0:12345          --label qwerty
	sudo $(CMD) daemonize --in tcp/client:127.0.0.1:12345 --out udp/broadcast:192.168.1.255:60005 --label uiop

undaemonize: build
	sudo $(CMD) undaemonize --label qwerty
	sudo $(CMD) undaemonize --label uiop

