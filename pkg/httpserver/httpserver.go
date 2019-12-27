package httpserver

import (
	"context"
	"net/http"
	"sync"
	"tkestack.io/kube-jarvis/pkg/logger"
)

var handlers = map[string]func(http.ResponseWriter, *http.Request){}
var handlersLock sync.Mutex

type Handler interface {
	Hand(ctx context.Context) (response interface{})
}

func Start(logger logger.Logger, addr string) {
	if len(handlers) == 0 {
		return
	}

	if addr == "" {
		addr = ":9005"
	}
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			panic(err.Error())
		}
	}()
	logger.Infof("http server start at %s", addr)
}

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	handlersLock.Lock()
	defer handlersLock.Unlock()
	if _, exist := handlers[pattern]; exist {
		return
	}
	handlers[pattern] = handler
	http.HandleFunc(pattern, handler)
}
