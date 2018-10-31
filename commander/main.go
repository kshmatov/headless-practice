package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/raff/godet"
)

type handler func(godet.Params)
type handlers map[string]handler

type ChromeConnection struct {
	remote    *godet.RemoteDebugger
	handlers  handlers
	enableAll bool
}

func getConnection(addr string) (*ChromeConnection, error) {
	remote, err := godet.Connect(addr, false)
	if err != nil {
		return nil, err
	}
	c := ChromeConnection{remote: remote, handlers: make(handlers)}
	return &c, nil
}

func (cc *ChromeConnection) SetHandler(event string, handler handler) {
	if handler == nil {
		delete(cc.handlers, event)
		return
	} else {
		cc.handlers[event] = handler
	}
}

func (cc *ChromeConnection) EnableAll(all bool) {
	cc.enableAll = all
}

func (cc *ChromeConnection) setEventListeners() {
	if cc.remote == nil {
		return
	}

	m := make(map[string]bool)

	for k, v := range cc.handlers {
		event := strings.Split(k, ".")[0]
		m[event] = true
		cc.remote.CallbackEvent(k, godet.EventCallback(v))
	}
	var err error

	for k, _ := range m {
		switch k {
		case "Runtime":
			err = cc.remote.RuntimeEvents(true)
		case "Network":
			err = cc.remote.NetworkEvents(true)
		case "Page":
			err = cc.remote.PageEvents(true)
		case "DOM":
			err = cc.remote.DOMEvents(true)
		case "Log":
			err = cc.remote.LogEvents(true)
		default:
			continue
		}
		fmt.Printf("Set --- > %v\n", err)
	}
}

func (cc *ChromeConnection) Navigate(url string) {
	if cc.remote == nil {
		return
	}

	cc.setEventListeners()
	if cc.enableAll {
		cc.remote.AllEvents(true)
	}

	s, e := cc.remote.Navigate(url)
	fmt.Printf("---> %v, %v\n", s, e)
}

func (cc *ChromeConnection) Close() {
	if cc.remote == nil {
		return
	}
	cc.remote.Close()
	cc.remote = nil
	cc.handlers = nil
}

func main() {
	host := flag.String("host", "localhost", "Host where remote debugger is running")
	port := flag.Int("port", 9222, "Port where remote debugger is listening")

	addr := fmt.Sprintf("%v:%d", *host, *port)

	r1, e := getConnection(addr)
	if e != nil {
		panic(e)
	}

	r2, e := getConnection(addr)
	if e != nil {
		panic(e)
	}

	f := func(params godet.Params) {
		p, ok := params["response"]
		if !ok {
			return
		}

		ph := p.(map[string]interface{})

		for k, v := range ph["headers"].(map[string]interface{}) {
			fmt.Println(k, " - > ", v)
		}
	}

	r1.SetHandler("Network.responseReceived", f)
	r2.SetHandler("Network.responseReceived", f)

	r1.Navigate("http://localhost:8080/hello")
	time.Sleep(time.Second * 3)
	r2.Navigate("localhost:8080/hello-R2")
	time.Sleep(time.Second)
	// r1.Navigate("http://localhost:8080/hello-2")
	time.Sleep(time.Second * 5)
}
