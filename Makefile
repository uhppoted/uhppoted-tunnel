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
	mkdir -p dist/$(DIST)/arm
	mkdir -p dist/$(DIST)/arm7
	env GOOS=linux   GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/linux   ./...
	env GOOS=linux   GOARCH=arm64         GOWORK=off go build -trimpath -o dist/$(DIST)/arm     ./...
	env GOOS=linux   GOARCH=arm   GOARM=7 GOWORK=off go build -trimpath -o dist/$(DIST)/arm7    ./...
	env GOOS=darwin  GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/darwin  ./...
	env GOOS=windows GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/windows ./...

release: update-release build-all
	find . -name ".DS_Store" -delete
	tar --directory=dist --exclude=".DS_Store" -cvzf dist/$(DIST).tar.gz $(DIST)

publish: release
	echo "Releasing version $(VERSION)"
	rm -f dist/development.tar.gz
	gh release create "$(VERSION)" ./dist/*.tar.gz --draft --prerelease --title "$(VERSION)-beta" --notes-file release-notes.md

debug: build
#	$(CMD) --config "./examples/uhppoted-tunnel.toml#client-en3" --console --debug
#	$(CMD) --config "./examples/uhppoted-tunnel.toml#client-lo0" --console --debug
#	$(CMD) --debug --console --in tcp/client:149.248.55.183:12345 --out udp/broadcast:192.168.1.255:60000 --udp-timeout 1s
#	$(CMD) --debug --console --in tcp/client::en3:149.248.55.183:12345 --out udp/broadcast:192.168.1.255:60000 --udp-timeout 1s
#	$(CMD) --debug --console --in tcp/client::en0:149.248.55.183:12345 --out udp/broadcast:192.168.1.255:60000 --udp-timeout 1s
#	$(CMD) --debug --console --in tls/client:::127.0.0.1:12345 --out udp/broadcast:192.168.1.255:60005 --udp-timeout 1s
#	$(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out tcp/server::lo0:127.0.0.1:12345
#	$(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out tls/server::lo0:127.0.0.1:12345 --client-auth
#	$(CMD) --debug --console --in tcp/client:155.138.144.189:12345 --out udp/broadcast::en3:192.168.1.255:60000 --udp-timeout 1s
	$(CMD) --debug --console --in udp/listen::en3:0.0.0.0:60000 --out tcp/client:155.138.144.189:12345

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
#	$(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out tcp/server:0.0.0.0:12345
	$(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out tcp/server::lo0:127.0.0.1:12345

client: build
	# $(CMD) --debug --console --in tcp/client:127.0.0.1:12345 --out udp/broadcast:192.168.1.255:60005 --udp-timeout 1s
	# $(CMD) --config "#client" --console --debug
	# $(CMD) --config "./examples/uhppoted-tunnel.toml#client"
	$(CMD) --config "./examples/uhppoted-tunnel.toml#client" --console --debug

client-ethernet: build
	$(CMD) --config "./examples/uhppoted-tunnel.toml#client-ethernet"

client-wifi: build
	$(CMD) --config "./examples/uhppoted-tunnel.toml#client-wifi"

reverse-host: build
	$(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out tcp/client:127.0.0.1:12345

reverse-client: build
	$(CMD) --debug --console --in tcp/server:0.0.0.0:12345 --out udp/broadcast:192.168.1.255:60005 --udp-timeout 1s

tls-host: build
	$(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out tls/server:0.0.0.0:12345 --client-auth

tls-client: build
	$(CMD) --debug --console --in tls/client:127.0.0.1:12345 --out udp/broadcast:192.168.1.255:60005 --udp-timeout 1s

event-client: build
	# $(CMD) --debug --console --in udp/event:0.0.0.0:60001 --out tcp/event/client:127.0.0.1:12345
	$(CMD) --debug --console --in udp/event:0.0.0.0:60001 --out tls/event/client:127.0.0.1:12345

event-host: build
	# $(CMD) --debug --console --in tcp/event/server:0.0.0.0:12345 --out udp/event:192.168.1.255:60005
	$(CMD) --debug --console --in tls/event/server:0.0.0.0:12345 --out udp/event:192.168.1.255:60005

reverse-event-client: build
	# $(CMD) --debug --console --in udp/event:0.0.0.0:60001 --out tcp/event/server:0.0.0.0:12345
	$(CMD) --debug --console --in udp/event:0.0.0.0:60001 --out tls/event/server:0.0.0.0:12345

reverse-event-host: build
	# $(CMD) --debug --console --in tcp/event/client:127.0.0.1:12345 --out udp/event:192.168.1.255:60005
	$(CMD) --debug --console --in tls/event/client:127.0.0.1:12345 --out udp/event:192.168.1.255:60005

http: build
	npx eslint --fix ./examples/html/javascript/*.js
	$(CMD) --debug --console --in http/0.0.0.0:8082 --out udp/broadcast:192.168.1.255:60005 --udp-timeout 1s --html ./examples/html

https: build
	npx eslint --fix ./examples/html/javascript/*.js
	$(CMD) --debug --console --in https/0.0.0.0:8443 --out udp/broadcast:192.168.1.255:60005 --udp-timeout 1s --html ./examples/html

daemonize: build
	sudo $(CMD) daemonize --in  udp/listen:0.0.0.0:60000  --out tcp/server:0.0.0.0:12345
	# sudo $(CMD) daemonize --in  udp/listen:0.0.0.0:60000  --out tcp/server:0.0.0.0:12345          --label qwerty
	# sudo $(CMD) daemonize --in tcp/client:127.0.0.1:12345 --out udp/broadcast:192.168.1.255:60005 --label uiop
	# sudo $(CMD) daemonize --config "./examples/uhppoted-tunnel.toml#client" --label qwerty
	# sudo $(CMD) daemonize
	# sudo $(CMD) daemonize --config "#client"

undaemonize: build
	sudo $(CMD) undaemonize --label qwerty
	# sudo $(CMD) undaemonize --label uiop

