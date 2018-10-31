PROJECT=github.com/kshmatov/headless-practice
BUILD=go build
OUTPUT=bin

WS=webserver

clean-ws:
	go clean $(PROJECT)/$(WS)

build-ws: clean-ws
	$(BUILD) -o $(OUTPUT)/$(WS) $(PROJECT)/$(WS)
