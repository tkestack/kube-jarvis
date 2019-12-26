package httpserver

import (
	"context"
	"net/http"
	"tkestack.io/kube-jarvis/pkg/logger"
)

var needStart = false

type Handler interface {
	Hand(ctx context.Context) (response interface{})
}

func Start(logger logger.Logger, addr string) {
	if !needStart {
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
	needStart = true
	http.HandleFunc(pattern, handler)
}
