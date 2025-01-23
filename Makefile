DIST   ?= development
DEBUG  ?= --debug
CMD     = ./bin/uhppoted-tunnel

.DEFAULT_GOAL := test
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
	go get -u github.com/uhppoted/uhppoted-lib@main
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

vet: 
	go vet ./...

lint: 
	env GOOS=darwin  GOARCH=amd64 staticcheck ./...
	env GOOS=linux   GOARCH=amd64 staticcheck ./...
	env GOOS=windows GOARCH=amd64 staticcheck ./...

benchmark: build
	go test -count 5 -bench=.  ./system/events

coverage: build
	go test -cover ./...

vuln:
	govulncheck ./...

build-all: build test vet lint
	mkdir -p dist/$(DIST)/linux
	mkdir -p dist/$(DIST)/arm
	mkdir -p dist/$(DIST)/arm7
	mkdir -p dist/$(DIST)/arm6
	mkdir -p dist/$(DIST)/darwin-x64
	mkdir -p dist/$(DIST)/darwin-arm64
	mkdir -p dist/$(DIST)/windows
	env GOOS=linux   GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/linux        ./...
	env GOOS=linux   GOARCH=arm64         GOWORK=off go build -trimpath -o dist/$(DIST)/arm          ./...
	env GOOS=linux   GOARCH=arm   GOARM=7 GOWORK=off go build -trimpath -o dist/$(DIST)/arm7         ./...
	env GOOS=linux   GOARCH=arm   GOARM=6 GOWORK=off go build -trimpath -o dist/$(DIST)/arm6         ./...
	env GOOS=darwin  GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/darwin-x64   ./...
	env GOOS=darwin  GOARCH=arm64         GOWORK=off go build -trimpath -o dist/$(DIST)/darwin-arm64 ./...
	env GOOS=windows GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/windows      ./...

release: update-release build-all
	find . -name ".DS_Store" -delete
	tar --directory=dist/$(DIST)/linux        --exclude=".DS_Store" -cvzf dist/$(DIST)-linux-x64.tar.gz    .
	tar --directory=dist/$(DIST)/arm          --exclude=".DS_Store" -cvzf dist/$(DIST)-arm-x64.tar.gz      .
	tar --directory=dist/$(DIST)/arm7         --exclude=".DS_Store" -cvzf dist/$(DIST)-arm7.tar.gz         .
	tar --directory=dist/$(DIST)/arm6         --exclude=".DS_Store" -cvzf dist/$(DIST)-arm6.tar.gz         .
	tar --directory=dist/$(DIST)/darwin-x64   --exclude=".DS_Store" -cvzf dist/$(DIST)-darwin-x64.tar.gz   .
	tar --directory=dist/$(DIST)/darwin-arm64 --exclude=".DS_Store" -cvzf dist/$(DIST)-darwin-arm64.tar.gz .
	cd dist/$(DIST)/windows && zip --recurse-paths ../../$(DIST)-windows-x64.zip . -x ".DS_Store"

publish: release
	echo "Releasing version $(VERSION)"
	gh release create "$(VERSION)" "./dist/$(DIST)-arm-x64.tar.gz"      \
	                               "./dist/$(DIST)-arm7.tar.gz"         \
	                               "./dist/$(DIST)-arm6.tar.gz"         \
	                               "./dist/$(DIST)-darwin-arm64.tar.gz" \
	                               "./dist/$(DIST)-darwin-x64.tar.gz"   \
	                               "./dist/$(DIST)-linux-x64.tar.gz"    \
	                               "./dist/$(DIST)-windows-x64.zip"     \
	                               --draft --prerelease --title "$(VERSION)-beta" --notes-file release-notes.md

debug: build
	# $(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out ip/out:192.168.1.255:60005 --udp-timeout 1s
	$(CMD) --config "./examples/uhppoted-tunnel.toml#ip"

delve: build
#   dlv exec ./bin/uhppoted-tunnel -- --debug --console
	dlv test github.com/uhppoted/uhppoted-tunnel -- run Test*

godoc:
	godoc -http=:80	-index_interval=60s

version: build
	$(CMD) version

help: build
	$(CMD) help
	$(CMD) help commands
	$(CMD) help version
	$(CMD) help help

host: build
#	$(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out tcp/server:0.0.0.0:12345
	$(CMD) --debug --console --in udp/listen:0.0.0.0:60005 --out tcp/server::lo0:127.0.0.1:12345

# host-udp: build
# 	$(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out udp/broadcast:192.168.1.255:60005 --udp-timeout 1s

# host-tcp: build
# 	$(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out tcp/client:192.168.1.100:60005

client: build
	$(CMD) --debug --console --in tcp/client:127.0.0.1:12345 --out udp/broadcast:192.168.1.255:60000 --udp-timeout 1s
	# $(CMD) --config "#client" --console --debug
	# $(CMD) --config "./examples/uhppoted-tunnel.toml#client"
	# $(CMD) --config "./examples/uhppoted-tunnel.toml#client" --console --debug

client-ethernet: build
	$(CMD) --config "./examples/uhppoted-tunnel.toml#client-ethernet"

client-wifi: build
	$(CMD) --config "./examples/uhppoted-tunnel.toml#client-wifi"

reverse-host: build
	$(CMD) --debug --console --in udp/listen:0.0.0.0:60005 --out tcp/client:127.0.0.1:12345

reverse-client: build
	$(CMD) --debug --console --in tcp/server:0.0.0.0:12345 --out udp/broadcast:192.168.1.255:60000 --udp-timeout 1s

tls-host: build
	$(CMD) --debug --console --in udp/listen:0.0.0.0:60000 --out tls/server:0.0.0.0:12345 \
	       --ca-cert ../runtime/uhppoted-tunnel/ca.cert     \
	       --key     ../runtime/uhppoted-tunnel/server.key  \
	       --cert    ../runtime/uhppoted-tunnel/server.cert \
           --client-auth

tls-client: build
	$(CMD) --debug --console --in tls/client:127.0.0.1:12345 --out udp/broadcast:192.168.1.255:60005 \
	       --ca-cert ../runtime/uhppoted-tunnel/ca.cert     \
	       --key     ../runtime/uhppoted-tunnel/client.key  \
	       --cert    ../runtime/uhppoted-tunnel/client.cert
	       --udp-timeout 1s

event-client: build
	$(CMD) --debug --console --in udp/event:0.0.0.0:60001 --out tcp/client:127.0.0.1:12345

event-host: build
	$(CMD) --debug --console --in tcp/server:0.0.0.0:12345 --out udp/event:192.168.1.255:60005
	# $(CMD) --debug --console --in tls/server:0.0.0.0:12345 --out udp/event:192.168.1.255:60005

tls-event-client: build
	$(CMD) --debug --console --in udp/event:0.0.0.0:60001 --out tls/client:127.0.0.1:12345

tls-event-host: build
	$(CMD) --debug --console --in tls/server:0.0.0.0:12345 --out udp/event:192.168.1.255:60005

reverse-event-client: build
	$(CMD) --debug --console --in udp/event:0.0.0.0:60001 --out tcp/server:0.0.0.0:12345
	# $(CMD) --debug --console --in udp/event:0.0.0.0:60001 --out tls/server:0.0.0.0:12345

reverse-event-host: build
	$(CMD) --debug --console --in tcp/client:127.0.0.1:12345 --out udp/event:192.168.1.255:60005
	# $(CMD) --debug --console --in tls/client:127.0.0.1:12345 --out udp/event:192.168.1.255:60005

tls-reverse-event-client: build
	$(CMD) --debug --console --in udp/event:0.0.0.0:60001 --out tls/server:0.0.0.0:12345

tls-reverse-event-host: build
	$(CMD) --debug --console --in tls/client:127.0.0.1:12345 --out udp/event:192.168.1.255:60005

http: build
	npx eslint --fix ./examples/html/javascript/*.js
	$(CMD) --debug --console --in http/0.0.0.0:8082 --out udp/broadcast:192.168.1.255:60000 --udp-timeout 1s --html ./examples/html

https: build
	npx eslint --fix ./examples/html/javascript/*.js
	$(CMD) --debug --console --in https/0.0.0.0:8443 --out udp/broadcast:192.168.1.255:60000 --udp-timeout 1s --html ./examples/html

tailscale-client: build
#	$(CMD) --debug --console --workdir ../runtime/uhppoted-tunnel --in udp/listen:0.0.0.0:60000 --out tailscale/client::qwerty:uhppoted:12345,nolog
	$(CMD) --config "../runtime/uhppoted-tunnel/uhppoted-tunnel.toml#tailscale-client"

tailscale-server: build
#	$(CMD) --debug --console --workdir ../runtime/uhppoted-tunnel --in tailscale/server:uhppoted:12345,nolog --out udp/broadcast:192.168.1.255:60005 --udp-timeout 1s
	$(CMD) --config "../runtime/uhppoted-tunnel/uhppoted-tunnel.toml#tailscale-server"

tailscale-server-misconfigured: build
	$(CMD) --config "../runtime/uhppoted-tunnel/uhppoted-tunnel.toml#tailscale-server" --out udp/broadcast:192.168.1.255:60000

tailscale-event-client: build
	$(CMD) --debug --console --in udp/event:0.0.0.0:60001 --out tailscale/client::qwerty:uhppoted:12345,nolog

tailscale-event-server: build
	$(CMD) --debug --console --in tailscale/server:uhppoted:12345,nolog --out udp/event:192.168.1.255:60005

ip: build
	$(CMD) --config "./examples/uhppoted-tunnel.toml#ip"

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

