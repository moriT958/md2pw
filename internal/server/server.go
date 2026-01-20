package server

import (
	"log/slog"
	"net/http"
)

type Converter struct {
	http.Server
}

func NewServer(addr string) *Converter {
	srv := new(http.Server)
	srv.Addr = addr

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("ok")); err != nil {
			slog.Error("unhealthy",
				"error", err.Error(),
			)
		}
	})

	return &Converter{
		Server: *srv,
	}
}

func Start() {
	// gracefull start & shutdonw
}
