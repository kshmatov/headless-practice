PROJECT=github.com/kshmatov/headless-practice
BUILD=go build
OUTPUT=bin

WS=webserver
COM=commander

clean-ws:
	go clean $(PROJECT)/$(WS)

build-ws: clean-ws
	$(BUILD) -o $(OUTPUT)/$(WS) $(PROJECT)/$(WS)

build-c:
	$(BUILD) -o $(OUTPUT)/$(COM) $(PROJECT)/$(COM)

run:
	go run $(COM)/main.go
